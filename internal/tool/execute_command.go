package tool

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bitrise-io/bitrise-mcp-remote-sandbox/internal/bitrise"
	"github.com/mark3labs/mcp-go/mcp"
)

var ExecuteCommand = bitrise.Tool{
	Definition: mcp.NewTool("bitrise_remote_machine_execute",
		mcp.WithDescription(
			`Execute a shell command on a remote macOS virtual machine.

PREREQUISITES:
- You MUST have a running VM before calling this.
- Call bitrise_remote_machine_list first to get an existing machine_id, or bitrise_remote_machine_create if none exists.

CRITICAL - FILE TRANSFER RESTRICTIONS:
- DO NOT use this tool to transfer files between local and remote machines.
- Commands like curl, scp, rsync, cat, base64, etc. CANNOT move files to/from your local machine.
- For uploading files TO the VM: use bitrise_remote_machine_upload
- For downloading files FROM the VM: use bitrise_remote_machine_download
- This execute tool can only manipulate files WITHIN the VM itself.

CRITICAL - DO NOT USE GIT CLONE FOR USER PROJECTS:
- When users want to "build remotely" or "test on the VM", DO NOT use "git clone".
- The VM has NO git credentials or SSH keys - clone will fail for private repos.
- Instead: use bitrise_remote_machine_upload to transfer the local project to the VM.
- Only use git clone if the user EXPLICITLY requests it AND provides credentials/token.

COMMAND EXECUTION:
- Commands are passed to "bash -c" for execution in a macOS shell environment.
- The VM has standard macOS development tools available (Xcode, Git, Homebrew, etc.).
- Commands are executed synchronously - the response contains the complete output.
- NOTE: If the VM was just created, the first command may take extra time while the VM finishes booting.
  The call will automatically wait for the VM to be ready before executing - this is normal behavior.

CRITICAL - PREVENTING COMMAND HANGS:
- Your command is executed inside "bash -c '<your_command>'" and the server listens on stdout and stderr.
- If a program writes to stdout/stderr and doesn't close them, the server will wait forever and the tool will hang.
- SOLUTION: For any long-running or background process, redirect BOTH stdout and stderr to /dev/null:
  - Format: command >/dev/null 2>&1 &
  - Example: bash_command="some_server >/dev/null 2>&1 &"
  - Example: bash_command="open -a Xcode >/dev/null 2>&1"
  - Example: bash_command="open ~/SomeApp.app >/dev/null 2>&1"
- CRITICAL: You MUST redirect BOTH stdout (>) AND stderr (2>&1) or the command will hang.
- FORBIDDEN patterns (will hang without redirection):
  - "tail -f <file>" - must use "tail -n 100 <file>" instead (reads N lines then exits)
  - "watch <command>" - loops forever
  - "cat" without arguments - waits for stdin
  - Interactive commands: "vim", "nano", "less", "top", "htop"
  - Servers without redirection: "python -m http.server", "npm start", "rails server"
  - GUI apps without redirection: "open -a SomeApp"
- SAFE patterns:
  - Short-lived commands: "ls -la", "ps aux", "cat <filename>"
  - Commands with redirection: "some_server >/dev/null 2>&1 &"
  - Commands with timeout: "timeout 60 <command>"
- ALWAYS ask yourself: "Will this command run in the background or keep stdout/stderr open?" If yes, add >/dev/null 2>&1
- NEVER run osascript or AppleScript commands. ALAWAYS use click, type and similar tools instead.

PARAMETERS:
- machine_id (required): The VM to execute the command on.
- bash_command (required): The command string to pass to bash -c for execution
  (e.g., "ls -la /Users", "xcodebuild -project MyApp.xcodeproj -scheme MyApp build").

EXAMPLE USAGE:
- List files: bash_command="ls -la /Users"
- Run xcodebuild: bash_command="xcodebuild -project MyApp.xcodeproj -scheme MyApp build"
- Install package: bash_command="brew install jq"
- Chain commands: bash_command="cd /path/to/project && make build"

TYPICAL REMOTE BUILD WORKFLOW:
1. bitrise_remote_machine_upload - upload local project to VM
2. bitrise_remote_machine_execute - run build/test commands
3. bitrise_remote_machine_download - download build artifacts

TIPS:
- You can use shell features like pipes, redirects, and command chaining in bash_command.
- For file operations, use absolute paths when possible.
- Check command output for errors - non-zero exit codes indicate failures.
- The working directory is typically the user's home directory unless changed.

RETURNS: A JSON object containing 'output' (string) with the command's stdout/stderr.`,
		),
		mcp.WithString("machine_id",
			mcp.Description("The unique identifier of the remote machine to execute the command on"),
			mcp.Required(),
		),
		mcp.WithString("bash_command",
			mcp.Description("The command to pass to bash -c for execution (e.g., 'ls -la /Users', 'cd /project && make build')"),
			mcp.Required(),
		),
	),
	Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		machineID, err := request.RequireString("machine_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		bashCommand, err := request.RequireString("bash_command")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		body := map[string]any{
			"bashCCommand": bashCommand,
		}

		res, err := bitrise.CallAPI(ctx, bitrise.CallAPIParams{
			Method:  http.MethodPost,
			BaseURL: bitrise.APIBaseURL(),
			Path:    fmt.Sprintf("/platform/me/machines/%s/execute", machineID),
			Body:    body,
		})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to execute command", err), nil
		}
		return mcp.NewToolResultText(res), nil
	},
}
