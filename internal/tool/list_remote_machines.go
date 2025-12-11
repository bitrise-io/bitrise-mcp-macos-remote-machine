package tool

import (
	"context"
	"net/http"

	"github.com/bitrise-io/bitrise-mcp-remote-sandbox/internal/bitrise"
	"github.com/mark3labs/mcp-go/mcp"
)

var ListRemoteMachines = bitrise.Tool{
	Definition: mcp.NewTool("bitrise_remote_machine_list",
		mcp.WithDescription(
			`List all remote macOS virtual machines currently running for the authenticated user.

THIS SHOULD BE YOUR FIRST CALL when you need to execute commands on a VM.

WHY CALL THIS FIRST:
- Users can only have ONE VM running at a time.
- If a VM already exists, you MUST reuse it instead of creating a new one.
- Creating unnecessary VMs wastes time (30-60 seconds provisioning) and will fail.
- This call is fast and helps you determine the correct next action.

DECISION FLOW:
1. Call bitrise_remote_machine_list
2. If machine_ids array is NOT empty: use the existing machine_id for subsequent operations
3. If machine_ids array IS empty: call bitrise_remote_machine_create to provision a new VM

RETURNS: A JSON object containing 'machine_ids' (array of strings).
- Empty array [] means no VMs are running - you need to create one.
- Array with one ID means a VM exists - use that machine_id for operations.`,
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
		return mcp.NewToolResultText(res), nil
	},
}
