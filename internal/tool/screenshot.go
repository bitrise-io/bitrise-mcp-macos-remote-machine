package tool

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bitrise-io/bitrise-mcp-macos-remote-machine/internal/bitrise"
	"github.com/mark3labs/mcp-go/mcp"
)

type screenshotResponse struct {
	SignedURL string `json:"signedUrl"`
}

var Screenshot = bitrise.Tool{
	Definition: mcp.NewTool("bitrise_remote_machine_screenshot",
		mcp.WithDescription(
			`Take a screenshot of the current display on a remote macOS virtual machine.

PURPOSE:
This tool captures the current state of the VM's graphical display as an image. This is
essential for visual verification, debugging GUI applications, identifying coordinates
for click/drag operations, and documenting the state of the remote machine.

PREREQUISITES:
- You MUST have a running VM before calling this.
- Call bitrise_remote_machine_list first to get an existing machine_id, or bitrise_remote_machine_create if none exists.

SCREEN RESOLUTION - IMPORTANT:
The remote machine screen resolution is ALWAYS 1024x768 pixels. When identifying coordinates
from screenshots for click or drag operations, you MUST use absolute coordinates based on
this 1024x768 resolution. NEVER use relative coordinates or percentages.

CRITICAL: The screenshot image may be scaled or displayed at a different size, but the actual
screen resolution is ALWAYS 1024x768. You MUST calculate coordinates as if the screen is
1024x768, scaling from the screenshot dimensions if necessary.

Valid coordinate ranges for subsequent click/drag operations:
- x: 0 to 1023 (horizontal, absolute pixels)
- y: 0 to 767 (vertical, absolute pixels)

PARAMETERS:
- machine_id (required): The unique identifier of the remote machine to take a screenshot of.

RETURNS: The screenshot image data that can be displayed directly, along with the file path
where the screenshot was saved locally.

USAGE:
Use this tool to:
- Verify the current state of the VM's display
- Identify coordinates for subsequent click or drag operations
- Debug GUI-related issues
- Document the visual state of applications running on the VM

WORKFLOW EXAMPLE:
1. Take a screenshot to see current state
2. Identify target coordinates from the screenshot
3. Use bitrise_remote_machine_click or bitrise_remote_machine_type to interact
4. Take another screenshot to verify the result`,
		),
		mcp.WithString("machine_id",
			mcp.Description("The unique identifier of the remote machine to take a screenshot of"),
			mcp.Required(),
		),
	),
	Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		machineID, err := request.RequireString("machine_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		body := map[string]any{}

		res, err := bitrise.CallAPI(ctx, bitrise.CallAPIParams{
			Method:  http.MethodPost,
			BaseURL: bitrise.APIBaseURL(),
			Path:    fmt.Sprintf("/platform/me/machines/%s/screenshot", machineID),
			Body:    body,
		})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to take screenshot", err), nil
		}

		var screenshotResp screenshotResponse
		if err := json.Unmarshal([]byte(res), &screenshotResp); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to parse screenshot response: %v", err)), nil
		}

		// Download the image from signed URL
		imageData, err := downloadScreenshot(ctx, screenshotResp.SignedURL)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to download screenshot: %v", err)), nil
		}

		// Save to temporary file
		tempDir := os.TempDir()
		filename := fmt.Sprintf("screenshot_%s_%d.jpg", machineID, time.Now().UnixNano())
		filePath := filepath.Join(tempDir, filename)

		if err := os.WriteFile(filePath, imageData, 0644); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to save screenshot to file: %v", err)), nil
		}

		// Return both the embedded image and the file path
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewImageContent(base64.StdEncoding.EncodeToString(imageData), "image/jpeg"),
				mcp.NewTextContent(fmt.Sprintf("Screenshot saved to: %s", filePath)),
			},
		}, nil
	},
}

func downloadScreenshot(ctx context.Context, signedURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, signedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	client := http.Client{Timeout: 2 * time.Minute}
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
