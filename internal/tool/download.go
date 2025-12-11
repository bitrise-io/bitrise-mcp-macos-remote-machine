package tool

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitrise-io/bitrise-mcp-remote-sandbox/internal/bitrise"
	"github.com/mark3labs/mcp-go/mcp"
)

type downloadResponse struct {
	SignedURL string `json:"signedUrl"`
}

var Download = bitrise.Tool{
	Definition: mcp.NewTool("bitrise_remote_machine_download",
		mcp.WithDescription(
			`Download a file or folder from the remote macOS virtual machine.

PURPOSE:
This tool downloads content from the VM's filesystem to your local machine.
The content is transferred as a tar.gz archive and automatically extracted inside the destination parent folder.

WORKFLOW:
1. Provide the source_path on the VM and the local destination_parent_folder.
2. The tool will:
   - Request a download URL from the VM for the specified source
   - Download the tar.gz archive
   - Extract the content inside the destination_parent_folder

PARAMETERS:
- machine_id (required): The VM to download from.
- source_path (required): The absolute path on the VM of the file/folder to download
  (e.g., "/Users/user/project/build/output.ipa", "/Users/user/project/results/").
- destination_parent_folder (required): The absolute path on your local machine of
  the parent folder in which the content should be extracted (e.g., "/local/downloads", "/tmp/artifacts").
- only_contents_of_folder (optional): If true and source_path is a folder, only the contents
  of the folder will be archived and downloaded, not the folder itself. Defaults to false.
- open_after_download (optional): If true, automatically opens the downloaded file with the system's
  default application. Defaults to false. Useful for viewing images, opening documents or applications.

IMPORTANT NOTES:
- The content is downloaded as tar.gz and automatically extracted inside the destination_parent_folder.
- Parent directories will be created automatically if they don't exist.
- For large files like build artifacts, the download may take some time.

ERROR HANDLING:
- If download fails, read the error message carefully and RETRY the download.
- DO NOT try to work around download failures by using execute commands (e.g., cat, base64, scp).
- File transfers between remote and local MUST use this download tool or bitrise_remote_machine_upload.
- Common issues to check before retrying:
  1. Verify the source_path exists on the VM (use bitrise_remote_machine_execute with "ls" to check).
  2. Ensure destination_parent_folder is a valid writable path locally.
  3. For "file not found" errors, double-check the exact path on the VM.
- If a path issue is identified, fix it and retry - do not attempt alternative transfer methods.

COMMON USE CASES:
- Downloading iOS app builds (.ipa files) after xcodebuild
- Retrieving test results and logs
- Getting generated configuration or output files
- Extracting any files created during command execution

EXAMPLE USAGE:
1. bitrise_remote_machine_execute(machine_id="abc123", command="xcodebuild", args=["archive", ...])
2. bitrise_remote_machine_download(machine_id="abc123", source_path="/Users/user/build/MyApp.ipa", destination_parent_folder="/local/builds")

RETURNS: A success message or error details.`,
		),
		mcp.WithString("machine_id",
			mcp.Description("The unique identifier of the remote machine to download from"),
			mcp.Required(),
		),
		mcp.WithString("source_path",
			mcp.Description("The absolute path on the VM of the file/folder to download"),
			mcp.Required(),
		),
		mcp.WithString("destination_parent_folder",
			mcp.Description("The absolute path on your local machine of the parent folder in which the content should be extracted"),
			mcp.Required(),
		),
		mcp.WithBoolean("only_contents_of_folder",
			mcp.Description("If true and source_path is a folder, only the contents of the folder will be archived and downloaded, not the folder itself"),
		),
		mcp.WithBoolean("open_after_download",
			mcp.Description("If true, automatically opens the downloaded file with the system's default application"),
		),
	),
	Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		machineID, err := request.RequireString("machine_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		sourcePath, err := request.RequireString("source_path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		destinationParentFolder, err := request.RequireString("destination_parent_folder")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		onlyContentsOfFolder := request.GetBool("only_contents_of_folder", false)
		openAfterDownload := request.GetBool("open_after_download", false)

		// Step 1: Get download URL
		body := map[string]any{
			"sourcePath":           sourcePath,
			"onlyContentsOfFolder": onlyContentsOfFolder,
		}

		res, err := bitrise.CallAPI(ctx, bitrise.CallAPIParams{
			Method:  http.MethodPost,
			BaseURL: bitrise.APIBaseURL(),
			Path:    "/platform/me/machines/" + machineID + "/download",
			Body:    body,
		})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to get download URL", err), nil
		}

		var dlResp downloadResponse
		if err := json.Unmarshal([]byte(res), &dlResp); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to parse download response: %v", err)), nil
		}

		// Step 2: Download from signed URL
		data, err := downloadFromSignedURL(ctx, dlResp.SignedURL)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to download from signed URL: %v", err)), nil
		}

		// Step 3: Extract tar.gz to destination (like: tar -xzf archive -C destination)
		extractedPaths, err := extractTarGzWithPaths(data, destinationParentFolder)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to extract archive: %v", err)), nil
		}

		successMsg := fmt.Sprintf("Successfully downloaded %s from machine %s to %s", sourcePath, machineID, destinationParentFolder)

		// Step 4: Open the downloaded file if requested
		if openAfterDownload && len(extractedPaths) > 0 {
			// Find top-level items (files or folders directly in the destination folder)
			topLevelItems := getTopLevelItems(extractedPaths, destinationParentFolder)

			// If only one top-level item exists, open it (works for files, folders, and .app bundles)
			if len(topLevelItems) == 1 {
				if err := openURL(topLevelItems[0]); err != nil {
					successMsg += fmt.Sprintf("\n(Note: Failed to automatically open: %v)", err)
				} else {
					successMsg += "\nItem opened automatically."
				}
			} else {
				// If multiple top-level items were extracted, open the destination folder
				if err := openURL(destinationParentFolder); err != nil {
					successMsg += fmt.Sprintf("\n(Note: Failed to automatically open folder: %v)", err)
				} else {
					successMsg += "\nDestination folder opened automatically."
				}
			}
		}

		return mcp.NewToolResultText(successMsg), nil
	},
}

func downloadFromSignedURL(ctx context.Context, signedURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, signedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	client := http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("download failed with status %d: %s", resp.StatusCode, string(body))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return data, nil
}

func extractTarGzWithPaths(data []byte, destPath string) ([]string, error) {
	var extractedPaths []string
	destPath = filepath.Clean(destPath)
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return nil, fmt.Errorf("create destination directory: %w", err)
	}

	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read tar header: %w", err)
		}

		// Skip macOS AppleDouble files (extended attributes stored as ._filename)
		baseName := filepath.Base(header.Name)
		if strings.HasPrefix(baseName, "._") {
			continue
		}

		targetPath := filepath.Join(destPath, filepath.Clean(header.Name))

		// Security: ensure target is within destination
		if !isWithinDir(targetPath, destPath) {
			return nil, fmt.Errorf("path traversal detected: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return nil, fmt.Errorf("create directory %s: %w", targetPath, err)
			}
			extractedPaths = append(extractedPaths, targetPath)
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return nil, fmt.Errorf("create parent directory: %w", err)
			}
			if err := extractFile(tarReader, targetPath, header.Mode); err != nil {
				return nil, err
			}
			extractedPaths = append(extractedPaths, targetPath)
		case tar.TypeSymlink:
			// Validate symlink target
			linkTarget := filepath.Join(filepath.Dir(targetPath), header.Linkname)
			if !isWithinDir(linkTarget, destPath) {
				return nil, fmt.Errorf("symlink escapes destination: %s -> %s", header.Name, header.Linkname)
			}
			os.Remove(targetPath)
			if err := os.Symlink(header.Linkname, targetPath); err != nil {
				return nil, fmt.Errorf("create symlink: %w", err)
			}
			extractedPaths = append(extractedPaths, targetPath)
		}
	}
	return extractedPaths, nil
}

func isWithinDir(path, dir string) bool {
	path = filepath.Clean(path)
	dir = filepath.Clean(dir)
	return path == dir || strings.HasPrefix(path, dir+string(filepath.Separator))
}

func extractFile(tarReader *tar.Reader, targetPath string, mode int64) error {
	file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(mode))
	if err != nil {
		return fmt.Errorf("create file %s: %w", targetPath, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, tarReader); err != nil {
		return fmt.Errorf("write file %s: %w", targetPath, err)
	}

	return nil
}

// getTopLevelItems returns only the top-level items (files or folders) directly in the destination folder.
// This helps distinguish between a single .app bundle (which should be opened) vs multiple items (where we should open the parent folder).
func getTopLevelItems(extractedPaths []string, destinationParentFolder string) []string {
	topLevelMap := make(map[string]bool)
	destPath := filepath.Clean(destinationParentFolder)

	for _, path := range extractedPaths {
		cleanPath := filepath.Clean(path)

		// Get the relative path from destination
		relPath, err := filepath.Rel(destPath, cleanPath)
		if err != nil {
			continue
		}

		// Split the relative path and get the first component (top-level item)
		parts := strings.Split(relPath, string(filepath.Separator))
		if len(parts) > 0 && parts[0] != "." && parts[0] != ".." {
			topLevelPath := filepath.Join(destPath, parts[0])
			topLevelMap[topLevelPath] = true
		}
	}

	// Convert map to slice
	topLevelItems := make([]string, 0, len(topLevelMap))
	for item := range topLevelMap {
		topLevelItems = append(topLevelItems, item)
	}

	return topLevelItems
}
