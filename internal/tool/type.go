package tool

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bitrise-io/bitrise-mcp-remote-sandbox/internal/bitrise"
	"github.com/mark3labs/mcp-go/mcp"
)

var Type = bitrise.Tool{
	Definition: mcp.NewTool("bitrise_remote_machine_type",
		mcp.WithDescription(
			`Type text input on a remote macOS virtual machine.

PURPOSE:
This tool allows you to simulate keyboard input on the VM's graphical interface, typing
the specified text as if entered from a physical keyboard. This is useful for filling
forms, entering data into applications, writing in text editors, or any task requiring
keyboard input in GUI applications.

PREREQUISITES:
- You MUST have a running VM before calling this.
- Call bitrise_remote_machine_list first to get an existing machine_id, or bitrise_remote_machine_create if none exists.
- Ensure the appropriate text field or application is focused before typing.
- Use bitrise_remote_machine_click to focus on a text input field if needed.

PARAMETERS:
- machine_id (required): The unique identifier of the remote machine to type on.
- text (required): The text to type on the remote machine. Supports control characters and special keys.

RETURNS: An empty response on success.

USAGE:
1. Use bitrise_remote_machine_screenshot to see the current state
2. Use bitrise_remote_machine_click to focus on a text field if needed
3. Use this tool to type the desired text
4. Use bitrise_remote_machine_screenshot to verify the input

CONTROL CHARACTERS AND SPECIAL KEYS:
This tool supports control characters for special key input. Always use this tool for key input
instead of AppleScript or other scripting methods. When specifying control characters, use the
literal escape sequences (do not double-escape them). Supported control characters include:
- \n or \r - Enter/Return key (use literal \n, not \\n)
- \t - Tab key (use literal \t, not \\t)
- \b - Backspace key (use literal \b, not \\b)
- \e or \x1b - Escape key (use literal \e or \x1b, not \\e or \\x1b)

For keyboard shortcuts, combine control characters as needed.

IMPORTANT: Send raw control character sequences, not escaped versions:
- CORRECT: text="hello\nworld" (sends "hello" + Enter + "world")
- INCORRECT: text="hello\\nworld" (would send literal backslash-n)

NOTES:
- The text is typed as-is, including special characters and control sequences
- Always prefer this tool for ALL keyboard input, including special keys like Enter, Tab, Escape
- Do NOT use AppleScript or shell commands for keyboard input - use this tool instead
- Long text strings are typed sequentially`,
		),
		mcp.WithString("machine_id",
			mcp.Description("The unique identifier of the remote machine to type on"),
			mcp.Required(),
		),
		mcp.WithString("text",
			mcp.Description("The text to type on the remote machine"),
			mcp.Required(),
		),
	),
	Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		machineID, err := request.RequireString("machine_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		text, err := request.RequireString("text")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		body := map[string]any{
			"text": text,
		}

		res, err := bitrise.CallAPI(ctx, bitrise.CallAPIParams{
			Method:  http.MethodPost,
			BaseURL: bitrise.APIBaseURL(),
			Path:    fmt.Sprintf("/platform/me/machines/%s/type", machineID),
			Body:    body,
		})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to type text", err), nil
		}
		return mcp.NewToolResultText(res), nil
	},
}
