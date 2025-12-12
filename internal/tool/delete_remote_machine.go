package tool

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bitrise-io/bitrise-mcp-macos-remote-machine/internal/bitrise"
	"github.com/mark3labs/mcp-go/mcp"
)

var DeleteRemoteMachine = bitrise.Tool{
	Definition: mcp.NewTool("bitrise_remote_machine_delete",
		mcp.WithDescription(
			`Delete (terminate) a remote macOS virtual machine.

WHEN TO DELETE:
- ONLY when you are completely finished with ALL tasks the user has requested.
- When the user explicitly asks to clean up or terminate the VM.
- When you encounter an unrecoverable error and need a fresh environment.

WHEN NOT TO DELETE:
- If the user might have follow-up tasks or questions that require command execution.
- If you're in the middle of a multi-step workflow.
- If unsure whether the user is done - ASK FIRST before deleting.

IMPORTANT CONSIDERATIONS:
- Deletion is IMMEDIATE and IRREVERSIBLE - all data on the VM is lost.
- After deletion, any subsequent commands will require creating a new VM (30-60 second wait).
- VMs auto-expire after 1 hour anyway, so deletion is mainly for immediate cleanup.
- Being conservative about deletion provides better user experience than having to recreate VMs.

COST AWARENESS:
- Keeping a VM running unnecessarily consumes resources.
- Deleting too early and recreating wastes more time and resources than keeping it a bit longer.
- Strike a balance: delete when genuinely done, but don't delete prematurely.

REQUIRED PARAMETERS:
- machine_id: The ID of the VM to delete (obtained from bitrise_remote_machine_create or bitrise_remote_machine_list).`,
		),
		mcp.WithString("machine_id",
			mcp.Description("The unique identifier of the remote machine to delete"),
			mcp.Required(),
		),
	),
	Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		machineID, err := request.RequireString("machine_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		res, err := bitrise.CallAPI(ctx, bitrise.CallAPIParams{
			Method:  http.MethodDelete,
			BaseURL: bitrise.APIBaseURL(),
			Path:    fmt.Sprintf("/platform/me/machines/%s", machineID),
		})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to delete remote machine", err), nil
		}
		return mcp.NewToolResultText(res), nil
	},
}
