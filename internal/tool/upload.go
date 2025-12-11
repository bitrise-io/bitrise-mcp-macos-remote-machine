package tool

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bitrise-io/bitrise-mcp-remote-sandbox/internal/bitrise"
	"github.com/mark3labs/mcp-go/mcp"
)

type startUploadResponse struct {
	SignedURL string `json:"signedUrl"`
	UploadID  string `json:"uploadId"`
}

var Upload = bitrise.Tool{
	Definition: mcp.NewTool("bitrise_remote_machine_upload",
		mcp.WithDescription(
			`Upload a file or folder to the remote macOS virtual machine.

PURPOSE:
This tool uploads a local file or folder to the VM. It automatically handles compression
(tar.gz) for folders, uploads the content, and places it at the specified parent folder path.

CRITICAL - ALWAYS USE THIS TOOL INSTEAD OF GIT CLONE:
- When the user asks to "build remotely", "run tests remotely", "compile on the VM", etc.,
  ALWAYS use this upload tool to transfer the local project files to the VM.
- DO NOT attempt to use "git clone" on the remote machine - it will fail because:
  1. The VM does not have SSH keys or Git credentials configured.
  2. Private repositories will be inaccessible.
  3. Even public repos may have rate limits or network issues.
- The correct workflow is: upload local files → build/test on VM → download results.
- Only use git clone if the user EXPLICITLY requests it AND provides credentials/token.

WORKFLOW:
1. Provide the local source_path (file or folder) and the destination_parent_folder on the VM.
2. The tool will:
   - Create a tar.gz archive if the source is a folder (preserving relative paths)
   - Upload the content to the VM
   - Extract and place the content inside the destination_parent_folder

PARAMETERS:
- machine_id (required): The VM to upload the file/folder to.
- source_path (required): The absolute path to the local file or folder to upload.
- destination_parent_folder (required): The absolute path on the VM of the parent folder in which the content should be placed
  (e.g., "/Users/user/project/", "/tmp/myfiles/").
- only_contents_of_folder (optional): If true and source_path is a folder, only the contents
  of the folder will be archived and uploaded, not the folder itself. Defaults to false.

IMPORTANT NOTES:
- For folders, the content is compressed as tar.gz before upload and extracted on the VM.
- The destination_parent_folder should be an absolute path on the macOS filesystem.
- Parent directories will be created automatically if they don't exist.
- For large files/folders, the upload may take some time.

ERROR HANDLING:
- If upload fails, read the error message carefully and retry the upload.
- DO NOT try to work around upload failures by using execute commands (e.g., curl, scp, rsync).
- File transfers between local and remote MUST use this upload tool or bitrise_remote_machine_download.
- Common issues: check that source_path exists locally and destination_parent_folder is a valid path.

EXAMPLE USAGE:
- Upload a single file:
  bitrise_remote_machine_upload(machine_id="abc123", source_path="/local/file.txt", destination_parent_folder="/Users/user")

- Upload a project folder for remote build:
  bitrise_remote_machine_upload(machine_id="abc123", source_path="/local/myproject", destination_parent_folder="/Users/user/myproject")

RETURNS: A success message or error details.`,
		),
		mcp.WithString("machine_id",
			mcp.Description("The unique identifier of the remote machine to upload to"),
			mcp.Required(),
		),
		mcp.WithString("source_path",
			mcp.Description("The absolute path to the local file or folder to upload"),
			mcp.Required(),
		),
		mcp.WithString("destination_parent_folder",
			mcp.Description("The absolute path on the VM of the parent folder in which the content should be placed"),
			mcp.Required(),
		),
		mcp.WithBoolean("only_contents_of_folder",
			mcp.Description("If true and source_path is a folder, only the contents of the folder will be archived and uploaded, not the folder itself"),
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

		// Check if source exists (use Lstat to not follow symlinks)
		sourceInfo, err := os.Lstat(sourcePath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to access source path: %v", err)), nil
		}

		// Step 1: Start upload to get signed URL and upload ID
		startRes, err := bitrise.CallAPI(ctx, bitrise.CallAPIParams{
			Method:  http.MethodPost,
			BaseURL: bitrise.APIBaseURL(),
			Path:    "/platform/me/machines/" + machineID + "/start_upload",
		})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to start upload", err), nil
		}

		var startResp startUploadResponse
		if err := json.Unmarshal([]byte(startRes), &startResp); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to parse start upload response: %v", err)), nil
		}

		// Step 2: Create tar.gz archive
		var buf bytes.Buffer
		if err := createTarGz(&buf, sourcePath, sourceInfo, onlyContentsOfFolder); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create tar.gz archive: %v", err)), nil
		}

		// Step 3: Upload to signed URL
		if err := uploadToSignedURL(ctx, startResp.SignedURL, buf.Bytes()); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to upload to signed URL: %v", err)), nil
		}

		// Step 4: Complete upload
		completeBody := map[string]any{
			"uploadId":                startResp.UploadID,
			"destinationParentFolder": destinationParentFolder,
		}

		_, err = bitrise.CallAPI(ctx, bitrise.CallAPIParams{
			Method:  http.MethodPost,
			BaseURL: bitrise.APIBaseURL(),
			Path:    "/platform/me/machines/" + machineID + "/complete_upload",
			Body:    completeBody,
		})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to complete upload", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully uploaded %s to %s on machine %s", sourcePath, destinationParentFolder, machineID)), nil
	},
}

// createTarGz creates a tar.gz archive.
// If onlyContentsOfFolder is true and sourcePath is a directory, only the contents are archived.
// If onlyContentsOfFolder is false and sourcePath is a directory, the directory itself is included.
func createTarGz(buf *bytes.Buffer, sourcePath string, sourceInfo os.FileInfo, onlyContentsOfFolder bool) error {
	gzWriter := gzip.NewWriter(buf)
	tarWriter := tar.NewWriter(gzWriter)

	var err error
	if sourceInfo.IsDir() {
		// Determine the base path for computing relative paths in the archive
		var basePath string
		if onlyContentsOfFolder {
			// Use sourcePath as base, so contents are at root of archive
			basePath = sourcePath
		} else {
			// Use parent of sourcePath, so the folder itself is included
			basePath = filepath.Dir(sourcePath)
		}

		err = filepath.WalkDir(sourcePath, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			// Get file info (Lstat behavior - doesn't follow symlinks)
			info, err := d.Info()
			if err != nil {
				return fmt.Errorf("get file info for %s: %w", path, err)
			}

			// Get relative path from basePath
			relPath, err := filepath.Rel(basePath, path)
			if err != nil {
				return err
			}

			// Skip the root directory itself
			if relPath == "." {
				return nil
			}

			// Skip non-regular files (devices, sockets, etc.) except dirs and symlinks
			mode := info.Mode()
			if !mode.IsRegular() && !mode.IsDir() && mode&os.ModeSymlink == 0 {
				return nil // Skip special files
			}

			return addToTar(tarWriter, path, relPath, info)
		})
	} else {
		// Single file - use the base name as the archive path
		err = addToTar(tarWriter, sourcePath, filepath.Base(sourcePath), sourceInfo)
	}

	if err != nil {
		tarWriter.Close()
		gzWriter.Close()
		return err
	}

	// Close in correct order: tar first, then gzip
	if err := tarWriter.Close(); err != nil {
		gzWriter.Close()
		return fmt.Errorf("close tar writer: %w", err)
	}
	if err := gzWriter.Close(); err != nil {
		return fmt.Errorf("close gzip writer: %w", err)
	}

	return nil
}

func addToTar(tarWriter *tar.Writer, filePath, archivePath string, info os.FileInfo) error {
	// Handle symlinks
	var link string
	if info.Mode()&os.ModeSymlink != 0 {
		var err error
		link, err = os.Readlink(filePath)
		if err != nil {
			return fmt.Errorf("read symlink %s: %w", filePath, err)
		}
	}

	header, err := tar.FileInfoHeader(info, link)
	if err != nil {
		return fmt.Errorf("create tar header for %s: %w", filePath, err)
	}

	// Use the relative/archive path
	header.Name = archivePath

	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("write tar header for %s: %w", filePath, err)
	}

	// If it's a directory or symlink, we're done (no content to write)
	if info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil
	}

	// Write regular file content
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file %s: %w", filePath, err)
	}
	defer file.Close()

	if _, err := io.Copy(tarWriter, file); err != nil {
		return fmt.Errorf("write file content %s: %w", filePath, err)
	}

	return nil
}

func uploadToSignedURL(ctx context.Context, signedURL string, data []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, signedURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	client := http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
