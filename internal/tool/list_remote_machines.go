package tool

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bitrise-io/bitrise-mcp-macos-remote-machine/internal/bitrise"
	"github.com/mark3labs/mcp-go/mcp"
)

var ListRemoteMachines = bitrise.Tool{
	Definition: mcp.NewTool("bitrise_remote_machine_list",
		mcp.WithDescription(
			`List all remote macOS virtual machines owned by the authenticated user.

THIS SHOULD BE YOUR FIRST CALL when you need to work with VMs.

MULTI-MACHINE SUPPORT:
- Users can have up to 5 machines total (in any state).
- Only ONE machine can be running at a time.
- Machines can be in different states: running, pending, terminated, failed.
- Terminated machines preserve their disk state and can be restarted.

DECISION FLOW:
1. Call bitrise_remote_machine_list to see all your machines
2. Check the 'machines' array for machine details including status
3. If you need a running machine:
   - If a machine has status 'running': use it directly
   - If a machine has status 'terminated'/'failed': use bitrise_remote_machine_start
   - If no machines exist or you need a fresh one: use bitrise_remote_machine_create
4. If another machine is running and you need a different one, stop it first with bitrise_remote_machine_stop

MACHINE STATES:
- 'pending': Machine is starting up, wait for it to become 'running'
- 'running': Machine is active and ready for commands
- 'terminated': Machine was stopped or timed out, disk state preserved
- 'failed': Machine encountered an error

RETURNS: An array of machine objects, each containing:
- 'machine_id': Unique identifier
- 'description': User-provided description
- 'status': Current state (pending/running/terminated/failed)`,
		),
	),
	Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		res, err := bitrise.CallAPI(ctx, bitrise.CallAPIParams{
			Method:  http.MethodGet,
			BaseURL: bitrise.APIBaseURL(),
			Path:    "/platform/me/machines",
		})

		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to list remote machines", err), nil
		}

		// Extract only the 'machines' array from the response
		var response struct {
			Machines json.RawMessage `json:"machines"`
		}
		if err := json.Unmarshal([]byte(res), &response); err != nil {
			return mcp.NewToolResultText(res), nil // fallback to full response if parsing fails
		}

		return mcp.NewToolResultText(string(response.Machines)), nil
	},
}
