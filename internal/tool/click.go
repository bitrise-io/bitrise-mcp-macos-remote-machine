package tool

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bitrise-io/bitrise-mcp-macos-remote-machine/internal/bitrise"
	"github.com/mark3labs/mcp-go/mcp"
)

var Click = bitrise.Tool{
	Definition: mcp.NewTool("bitrise_remote_machine_click",
		mcp.WithDescription(
			`Perform a mouse click action on a remote macOS virtual machine.

PURPOSE:
This tool allows you to simulate mouse clicks on the VM's graphical interface at specified
coordinates. This is useful for automating GUI interactions, testing applications, or
performing tasks that require mouse input.

PREREQUISITES:
- You MUST have a running VM before calling this.
- Call bitrise_remote_machine_list first to get an existing machine_id, or bitrise_remote_machine_create if none exists.
- For visual feedback, consider using bitrise_remote_machine_screenshot before and after clicks.

SCREEN RESOLUTION - IMPORTANT:
The remote machine screen resolution is ALWAYS 1024x768 pixels. You MUST use absolute coordinates
based on this 1024x768 resolution. NEVER use relative coordinates or percentages.

CRITICAL: Even if a screenshot appears to have a different resolution, you MUST calculate and
provide coordinates as if the screen is 1024x768. Scale your coordinates accordingly.

Valid coordinate ranges:
- x: 0 to 1023 (horizontal, absolute pixels)
- y: 0 to 767 (vertical, absolute pixels)

PARAMETERS:
- machine_id (required): The unique identifier of the remote machine to perform the click on.
- x (required): The x-coordinate (horizontal position) for the click (0-1023).
- y (required): The y-coordinate (vertical position) for the click (0-767).
- button (required): The mouse button to click - "left", "right", or "middle".
- double_click (optional): Whether to perform a double click. Defaults to false.

RETURNS: An empty response on success.

USAGE:
Use this tool in combination with bitrise_remote_machine_screenshot to identify coordinates
and verify click results. Coordinates are relative to the screen's top-left corner (0,0).
The screen resolution is 1024x768.`,
		),
		mcp.WithString("machine_id",
			mcp.Description("The unique identifier of the remote machine to perform the click on"),
			mcp.Required(),
		),
		mcp.WithNumber("x",
			mcp.Description("The x-coordinate (horizontal position) for the click"),
			mcp.Required(),
		),
		mcp.WithNumber("y",
			mcp.Description("The y-coordinate (vertical position) for the click"),
			mcp.Required(),
		),
		mcp.WithString("button",
			mcp.Description("The mouse button to click: 'left', 'right', or 'middle'"),
			mcp.Required(),
		),
		mcp.WithBoolean("double_click",
			mcp.Description("Whether to perform a double click"),
		),
	),
	Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		machineID, err := request.RequireString("machine_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		x, err := request.RequireFloat("x")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		y, err := request.RequireFloat("y")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		button, err := request.RequireString("button")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		body := map[string]any{
			"x":      int(x),
			"y":      int(y),
			"button": button,
		}

		doubleClick := request.GetBool("double_click", false)
		if doubleClick {
			body["doubleClick"] = doubleClick
		}

		res, err := bitrise.CallAPI(ctx, bitrise.CallAPIParams{
			Method:  http.MethodPost,
			BaseURL: bitrise.APIBaseURL(),
			Path:    fmt.Sprintf("/platform/me/machines/%s/click", machineID),
			Body:    body,
		})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to perform click", err), nil
		}
		return mcp.NewToolResultText(res), nil
	},
}
