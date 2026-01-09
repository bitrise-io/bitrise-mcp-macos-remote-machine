package tool

import (
	"context"
	"net/http"

	"github.com/bitrise-io/bitrise-mcp-macos-remote-machine/internal/bitrise"
	"github.com/mark3labs/mcp-go/mcp"
)

var CreateRemoteMachine = bitrise.Tool{
	Definition: mcp.NewTool("bitrise_remote_machine_create",
		mcp.WithDescription(
			`Create a new remote macOS virtual machine for executing commands, running builds, or testing.

IMPORTANT CONSTRAINTS:
- Users can have up to 5 machines total (in any state).
- Only ONE machine can be running at a time per user.
- Before creating a new VM, ALWAYS call bitrise_remote_machine_list first to check existing machines.
- If a machine is already running, reuse it instead of creating a new one to save time.
- Creating or starting a VM takes time (typically 30-60 seconds) as it provisions/boots the macOS environment.

LIFECYCLE INFORMATION:
- The returned machine_id is required for ALL subsequent operations (execute, upload, download, etc.).
- The VM will appear in bitrise_remote_machine_list immediately after creation, but it needs time to boot up.
- The first bitrise_remote_machine_execute call may take longer as it waits for the VM to become ready.
- VMs will automatically terminate after 1 hour if not manually stopped or deleted.
- Use bitrise_remote_machine_stop to preserve disk state for later, or bitrise_remote_machine_delete to permanently remove.

WHEN TO CREATE A NEW VM:
- When you need a fresh macOS environment.
- When no running machine is available and you need a clean state.

BEST PRACTICES:
- FIRST call bitrise_remote_machine_list to check for existing machines.
- Reuse an already running machine whenever possible to avoid startup wait time.
- ALWAYS provide a description when you know the purpose of the machine (e.g., "iOS build for ProjectX",
  "Flutter development", "CI test environment"). This helps identify machines later.
- Stop machines when done to preserve state; delete only when you want a fresh start.

PARAMETERS:
- description (RECOMMENDED): A user-friendly description for the machine. Always provide this when you know
  why the machine is being created. Examples: "iOS build environment", "React Native dev", "Unit test runner".

RETURNS: A JSON object containing 'machine_id' (string) - save this for all subsequent operations.`,
		),
		mcp.WithString("description",
			mcp.Description("Optional description for the machine (e.g., 'iOS build environment')"),
		),
	),
	Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		description := request.GetString("description", "")

		body := map[string]any{}
		if description != "" {
			body["description"] = description
		}

		res, err := bitrise.CallAPI(ctx, bitrise.CallAPIParams{
			Method:  http.MethodPost,
			BaseURL: bitrise.APIBaseURL(),
			Path:    "/platform/me/machines",
			Body:    body,
		})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to create remote machine", err), nil
		}
		return mcp.NewToolResultText(res), nil
	},
}
