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
- You can only have ONE VM running at a time per user.
- Before creating a new VM, ALWAYS call bitrise_remote_machine_list first to check if one already exists.
- If a VM already exists, reuse it instead of trying to create a new one.
- Creating a VM takes time (typically 30-60 seconds) as it provisions a fresh macOS environment.

LIFECYCLE INFORMATION:
- The returned machine_id is required for ALL subsequent operations (execute, upload, download, delete).
- The VM will appear in bitrise_remote_machine_list immediately after creation, but it needs time to boot up.
- The first bitrise_remote_machine_execute call may take longer as it waits for the VM to become ready.
- VMs will automatically expire and terminate after 1 hour if not manually deleted.
- ALWAYS store the machine_id and reuse the same VM for related tasks to avoid unnecessary provisioning time.

WHEN TO CREATE A NEW VM:
- When bitrise_remote_machine_list returns an empty list and you need to execute commands.
- When you need a clean macOS environment for builds, tests, or shell operations.

BEST PRACTICES:
- FIRST call bitrise_remote_machine_list to check for existing VMs.
- Reuse existing VMs whenever possible - creating new ones wastes time.
- Only delete the VM when you are completely finished with ALL tasks the user requested.
- If the user might have follow-up tasks, ask before deleting the VM.
- Remember the 1-hour expiration: for long-running tasks, be aware of elapsed time.

RETURNS: A JSON object containing 'machine_id' (string) - save this for all subsequent operations.`,
		),
	),
	Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		res, err := bitrise.CallAPI(ctx, bitrise.CallAPIParams{
			Method:  http.MethodPost,
			BaseURL: bitrise.APIBaseURL(),
			Path:    "/platform/me/machines",
		})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to create remote machine", err), nil
		}
		return mcp.NewToolResultText(res), nil
	},
}
