package tool

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bitrise-io/bitrise-mcp-macos-remote-machine/internal/bitrise"
	"github.com/mark3labs/mcp-go/mcp"
)

var Scroll = bitrise.Tool{
	Definition: mcp.NewTool("bitrise_remote_machine_scroll",
		mcp.WithDescription(
			`Perform a scroll action on a remote macOS virtual machine.

PURPOSE:
This tool allows you to simulate scroll wheel actions on the VM's graphical interface.
This is useful for scrolling through documents, web pages, lists, or any scrollable content
in GUI applications.

PREREQUISITES:
- You MUST have a running VM before calling this.
- Call bitrise_remote_machine_list first to get an existing machine_id, or bitrise_remote_machine_create if none exists.
- For visual feedback, consider using bitrise_remote_machine_screenshot before and after scrolling.

SCREEN RESOLUTION - IMPORTANT:
The remote machine screen resolution is ALWAYS 1024x768 pixels.

CRITICAL: Even if a screenshot appears to have a different resolution, the actual screen
resolution is 1024x768. Keep this in mind when determining scroll amounts and positions.

PARAMETERS:
- machine_id (required): The unique identifier of the remote machine to perform the scroll on.
- direction (required): Direction to scroll - "up" or "down".
- amount (required): The amount to scroll (the unit typically corresponds to lines).

RETURNS: An empty response on success.

USAGE:
Specify the scroll direction and amount. Use "up" to scroll up (content moves down),
and "down" to scroll down (content moves up). The screen resolution is 1024x768.`,
		),
		mcp.WithString("machine_id",
			mcp.Description("The unique identifier of the remote machine to perform the scroll on"),
			mcp.Required(),
		),
		mcp.WithString("direction",
			mcp.Description("Direction to scroll: 'up' or 'down'"),
			mcp.Required(),
		),
		mcp.WithNumber("amount",
			mcp.Description("The amount to scroll"),
			mcp.Required(),
		),
	),
	Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		machineID, err := request.RequireString("machine_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		direction, err := request.RequireString("direction")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		amount, err := request.RequireFloat("amount")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		body := map[string]any{
			"direction": direction,
			"amount":    int(amount),
		}

		res, err := bitrise.CallAPI(ctx, bitrise.CallAPIParams{
			Method:  http.MethodPost,
			BaseURL: bitrise.APIBaseURL(),
			Path:    fmt.Sprintf("/platform/me/machines/%s/scroll", machineID),
			Body:    body,
		})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to perform scroll", err), nil
		}
		return mcp.NewToolResultText(res), nil
	},
}
