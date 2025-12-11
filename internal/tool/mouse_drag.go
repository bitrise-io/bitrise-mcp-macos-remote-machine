package tool

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bitrise-io/bitrise-mcp-remote-sandbox/internal/bitrise"
	"github.com/mark3labs/mcp-go/mcp"
)

var MouseDrag = bitrise.Tool{
	Definition: mcp.NewTool("bitrise_remote_machine_mouse_drag",
		mcp.WithDescription(
			`Perform a mouse drag action on a remote macOS virtual machine.

PURPOSE:
This tool allows you to simulate mouse drag operations on the VM's graphical interface,
moving from one coordinate to another while holding the mouse button. This is useful for
drag-and-drop operations, selecting text, resizing windows, or drawing operations.

PREREQUISITES:
- You MUST have a running VM before calling this.
- Call bitrise_remote_machine_list first to get an existing machine_id, or bitrise_remote_machine_create if none exists.
- For visual feedback, consider using bitrise_remote_machine_screenshot before and after drags.

SCREEN RESOLUTION - IMPORTANT:
The remote machine screen resolution is ALWAYS 1024x768 pixels. You MUST use absolute coordinates
based on this 1024x768 resolution. NEVER use relative coordinates or percentages.

CRITICAL: Even if a screenshot appears to have a different resolution, you MUST calculate and
provide coordinates as if the screen is 1024x768. Scale your coordinates accordingly.

Valid coordinate ranges:
- x: 0 to 1023 (horizontal, absolute pixels)
- y: 0 to 767 (vertical, absolute pixels)

PARAMETERS:
- machine_id (required): The unique identifier of the remote machine to perform the drag on.
- start_x (required): The starting x-coordinate (horizontal position) for the drag (0-1023).
- start_y (required): The starting y-coordinate (vertical position) for the drag (0-767).
- end_x (required): The ending x-coordinate (horizontal position) for the drag (0-1023).
- end_y (required): The ending y-coordinate (vertical position) for the drag (0-767).

RETURNS: An empty response on success.

USAGE:
Use this tool in combination with bitrise_remote_machine_screenshot to identify coordinates
and verify drag results. The screen resolution is 1024x768.`,
		),
		mcp.WithString("machine_id",
			mcp.Description("The unique identifier of the remote machine to perform the drag on"),
			mcp.Required(),
		),
		mcp.WithNumber("start_x",
			mcp.Description("The starting x-coordinate (horizontal position) for the drag"),
			mcp.Required(),
		),
		mcp.WithNumber("start_y",
			mcp.Description("The starting y-coordinate (vertical position) for the drag"),
			mcp.Required(),
		),
		mcp.WithNumber("end_x",
			mcp.Description("The ending x-coordinate (horizontal position) for the drag"),
			mcp.Required(),
		),
		mcp.WithNumber("end_y",
			mcp.Description("The ending y-coordinate (vertical position) for the drag"),
			mcp.Required(),
		),
	),
	Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		machineID, err := request.RequireString("machine_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		startX, err := request.RequireFloat("start_x")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		startY, err := request.RequireFloat("start_y")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		endX, err := request.RequireFloat("end_x")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		endY, err := request.RequireFloat("end_y")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		body := map[string]any{
			"startX": int(startX),
			"startY": int(startY),
			"endX":   int(endX),
			"endY":   int(endY),
		}

		res, err := bitrise.CallAPI(ctx, bitrise.CallAPIParams{
			Method:  http.MethodPost,
			BaseURL: bitrise.APIBaseURL(),
			Path:    fmt.Sprintf("/platform/me/machines/%s/mouse_drag", machineID),
			Body:    body,
		})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to perform mouse drag", err), nil
		}
		return mcp.NewToolResultText(res), nil
	},
}
