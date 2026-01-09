package tool

import (
	"context"
	"net/http"

	"github.com/bitrise-io/bitrise-mcp-macos-remote-machine/internal/bitrise"
	"github.com/mark3labs/mcp-go/mcp"
)

var StopRemoteMachine = bitrise.Tool{
	Definition: mcp.NewTool("bitrise_remote_machine_stop",
		mcp.WithDescription(
			`Stop a running remote macOS virtual machine, preserving its disk state.

PURPOSE:
This tool gracefully stops a running machine. The disk state is preserved, allowing you
to restart the machine later with all files and configurations intact. This is useful
when you want to pause work without losing progress.

WHEN TO USE:
- When you want to pause work on a machine but keep its state for later.
- When you need to free up your running machine slot to start a different machine.
- When you want to preserve the current state before a potential failure.

BENEFITS OF STOPPING VS DELETING:
- Disk state is preserved (files, installed software, configurations).
- Machine can be restarted with bitrise_remote_machine_start (note: starting takes 30-60 seconds).
- No need to re-upload files or reconfigure the environment.

CONSTRAINTS:
- The machine must be in 'running' or 'pending' state.
- Users can have up to 5 machines total (terminated machines count toward this limit).

REQUIRED PARAMETERS:
- machine_id: The ID of the machine to stop (from bitrise_remote_machine_list).

RETURNS: Empty response on success. The machine will transition to 'terminated' state.`,
		),
		mcp.WithString("machine_id",
			mcp.Description("The unique identifier of the remote machine to stop"),
			mcp.Required(),
		),
	),
	Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		machineID, err := request.RequireString("machine_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		res, err := bitrise.CallAPI(ctx, bitrise.CallAPIParams{
			Method:  http.MethodPost,
			BaseURL: bitrise.APIBaseURL(),
			Path:    "/platform/me/machines/stop",
			Body: map[string]any{
				"machineId": machineID,
			},
		})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to stop remote machine", err), nil
		}
		return mcp.NewToolResultText(res), nil
	},
}
