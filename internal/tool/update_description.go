package tool

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bitrise-io/bitrise-mcp-macos-remote-machine/internal/bitrise"
	"github.com/mark3labs/mcp-go/mcp"
)

var UpdateDescription = bitrise.Tool{
	Definition: mcp.NewTool("bitrise_remote_machine_update_description",
		mcp.WithDescription(
			`Update the description of a remote macOS virtual machine.

PURPOSE:
This tool allows you to set or update a user-friendly description for a machine.
Descriptions help identify machines when you have multiple machines in different states.

IMPORTANT - ALWAYS UPDATE WHEN PURPOSE CHANGES:
You MUST use this tool to update a machine's description whenever its current purpose
differs significantly from its existing description. For example:
- If a machine labeled "iOS build" is now being used for "Android testing", update it.
- If a machine has no description but you know what it's being used for, set one.
- If the scope of work on a machine has changed, reflect that in the description.
Keeping descriptions accurate is essential for managing multiple machines effectively.

WHEN TO USE:
- ALWAYS when a machine's actual use differs from its current description.
- When you want to label a machine with its purpose (e.g., "iOS build environment").
- When you want to add notes about what's installed or configured on the machine.
- When managing multiple machines and need to distinguish between them.

CONSTRAINTS:
- Works on machines in any state (running, pending, terminated, failed).
- The machine must exist and be owned by the authenticated user.

REQUIRED PARAMETERS:
- machine_id: The ID of the machine to update.
- description: The new description text for the machine.

RETURNS: Empty response on success.`,
		),
		mcp.WithString("machine_id",
			mcp.Description("The unique identifier of the remote machine to update"),
			mcp.Required(),
		),
		mcp.WithString("description",
			mcp.Description("The new description for the machine"),
			mcp.Required(),
		),
	),
	Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		machineID, err := request.RequireString("machine_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		description, err := request.RequireString("description")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		res, err := bitrise.CallAPI(ctx, bitrise.CallAPIParams{
			Method:  http.MethodPatch,
			BaseURL: bitrise.APIBaseURL(),
			Path:    fmt.Sprintf("/platform/me/machines/%s/description", machineID),
			Body: map[string]any{
				"description": description,
			},
		})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to update machine description", err), nil
		}
		return mcp.NewToolResultText(res), nil
	},
}
