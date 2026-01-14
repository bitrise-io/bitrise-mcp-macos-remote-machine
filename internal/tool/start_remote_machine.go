package tool

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bitrise-io/bitrise-mcp-macos-remote-machine/internal/bitrise"
	"github.com/mark3labs/mcp-go/mcp"
)

var StartRemoteMachine = bitrise.Tool{
	Definition: mcp.NewTool("bitrise_remote_machine_start",
		mcp.WithDescription(
			`Start an existing terminated or failed remote macOS virtual machine.

PURPOSE:
This tool starts a machine that was previously stopped or has terminated/failed. The machine
resumes from its preserved disk state, so any files and configurations from the previous
session are retained.

WHEN TO USE:
- When you have a terminated machine that you want to resume working with.
- When a machine has failed and you want to restart it.
- When you want to continue work on a machine without creating a new one.

CONSTRAINTS:
- Only ONE machine can be running at a time per user.
- The machine must be in 'terminated' or 'failed' state.
- If another machine is already running, you must stop it first.
- Starting a VM takes time (typically 30-60 seconds) as it boots the macOS environment.

MACHINE STATES THAT CAN BE STARTED:
- 'terminated': Machine was stopped or timed out, disk state preserved.
- 'failed': Machine encountered an error.

REQUIRED PARAMETERS:
- machine_id: The ID of the machine to start (from bitrise_remote_machine_list).
- description: A brief description of the machine's purpose. Always provide this parameter when
  the machine's purpose is changed compared to its current description.

RETURNS: Empty response on success. The machine will transition to 'pending' then 'running'.`,
		),
		mcp.WithString("machine_id",
			mcp.Description("The unique identifier of the remote machine to start"),
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

		body := map[string]any{}
		description := request.GetString("description", "")
		if description != "" {
			body["description"] = description
		}

		res, err := bitrise.CallAPI(ctx, bitrise.CallAPIParams{
			Method:  http.MethodPost,
			BaseURL: bitrise.APIBaseURL(),
			Path:    fmt.Sprintf("/platform/me/machines/%s/start", machineID),
			Body:    body,
		})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to start remote machine", err), nil
		}
		return mcp.NewToolResultText(res), nil
	},
}
