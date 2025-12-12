package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitrise-io/bitrise-mcp-macos-remote-machine/internal/bitrise"
	"github.com/mark3labs/mcp-go/mcp"
)

var OpenVNC = bitrise.Tool{
	Definition: mcp.NewTool("bitrise_remote_machine_open_vnc",
		mcp.WithDescription(
			`Open a VNC connection to a remote macOS virtual machine for graphical access.

PURPOSE:
This tool enables VNC (Virtual Network Computing) access to the VM, allowing you to
connect with a VNC client for graphical remote desktop access. This is useful when
the user would like to interact with the VM's GUI, run applications, or perform tasks
that require a graphical interface.

PREREQUISITES:
- You MUST have a running VM before calling this.
- Call bitrise_remote_machine_list first to get an existing machine_id, or bitrise_remote_machine_create if none exists.

PARAMETERS:
- machine_id (required): The unique identifier of the remote machine to open VNC connection to.

RETURNS: A JSON object containing:
- vncAddress: The address of the VNC server to connect to.
- vncUsername: The username for VNC authentication.
- vncPassword: The password for VNC authentication.

USAGE:
This tool automatically opens a VNC connection using your system's default VNC client.
It will return the VNC connection details and attempt to open the connection directly.
If automatic opening fails, you can manually connect using the returned credentials.`,
		),
		mcp.WithString("machine_id",
			mcp.Description("The unique identifier of the remote machine to open VNC connection to"),
			mcp.Required(),
		),
	),
	Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		machineID, err := request.RequireString("machine_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		body := map[string]any{}

		res, err := bitrise.CallAPI(ctx, bitrise.CallAPIParams{
			Method:  http.MethodPost,
			BaseURL: bitrise.APIBaseURL(),
			Path:    fmt.Sprintf("/platform/me/machines/%s/open_vnc", machineID),
			Body:    body,
		})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to open VNC connection", err), nil
		}

		// Parse the response to extract VNC details
		var vncDetails map[string]string
		if err := json.Unmarshal([]byte(res), &vncDetails); err == nil {
			vncAddress, addrOk := vncDetails["vncAddress"]
			vncUsername, userOk := vncDetails["vncUsername"]
			vncPassword, passOk := vncDetails["vncPassword"]

			if addrOk && vncAddress != "" {
				// Construct VNC URL with credentials if available
				// Format: vnc://username:password@host:port
				var vncURL string
				if userOk && passOk && vncUsername != "" && vncPassword != "" {
					vncURL = fmt.Sprintf("vnc://%s:%s@%s", vncUsername, vncPassword, vncAddress)
				} else {
					vncURL = fmt.Sprintf("vnc://%s", vncAddress)
				}

				// Attempt to open with system's default VNC client
				if err := openURL(vncURL); err != nil {
					// If automatic opening fails, just return the details
					return mcp.NewToolResultText(fmt.Sprintf("Failed to automatically open VNC client (%v), but VNC connection details are ready:\n\n%s", err, res)), nil
				}

				return mcp.NewToolResultText(fmt.Sprintf("VNC connection opened automatically:\n\n%s", res)), nil
			}
		}

		return mcp.NewToolResultText(res), nil
	},
}
