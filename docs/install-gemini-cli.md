# Install Bitrise Remote Machine MCP Server in Google Gemini CLI

## Prerequisites

1. The latest version of Google Gemini CLI installed (see [official Gemini CLI documentation](https://github.com/google-gemini/gemini-cli))
2. [Create a Bitrise API Token](https://devcenter.bitrise.io/api/authentication):
   - Go to your [Bitrise Account Settings/Security](https://app.bitrise.io/me/account/security).
   - Navigate to the "Personal access tokens" section.
   - Copy the generated token.
3. [Go](https://go.dev/) (>=1.25) installed

<details>
<summary><b>Storing Your PAT Securely</b></summary>
<br>

For security, avoid hardcoding your token. Create or update `~/.gemini/.env` (where `~` is your home or project directory) with your PAT:

```bash
# ~/.gemini/.env
BITRISE_PAT=your_token_here
```

</details>

## Setup

MCP servers for Gemini CLI are configured in its settings JSON under an `mcpServers` key.

- **Global configuration**: `~/.gemini/settings.json` where `~` is your home directory
- **Project-specific**: `.gemini/settings.json` in your project directory

After securely storing your PAT, you can add the Bitrise Remote Machine MCP server configuration to your settings file. You may need to restart the Gemini CLI for changes to take effect.

### Configuration

Add this to `~/.gemini/settings.json`:

```json
{
    "mcpServers": {
        "bitrise-remote-machine": {
            "command": "go",
            "args": [
                "run",
                "github.com/bitrise-io/bitrise-mcp-remote-sandbox@latest"
            ],
            "env": {
                "BITRISE_TOKEN": "$BITRISE_PAT"
            }
        }
    }
}
```

## Verification

To verify that the Bitrise Remote Machine MCP server has been configured, start Gemini CLI in your terminal with `gemini`, then:

1. **Check MCP server status**:

    ```
    /mcp list
    ```

    ```
    ℹ Configured MCP servers:

    🟢 bitrise-remote-machine - Ready (6 tools)
        - bitrise_remote_machine_list
        - bitrise_remote_machine_create
        - bitrise_remote_machine_delete
        - bitrise_remote_machine_execute
        - bitrise_remote_machine_upload
        - bitrise_remote_machine_download
    ```

2. **Test with a prompt**
    ```
    Create a remote machine and list the files in /Users
    ```

## Available Tools

The Bitrise Remote Machine MCP server provides 12 tools:

- **VM Lifecycle**: bitrise_remote_machine_list, bitrise_remote_machine_create, bitrise_remote_machine_delete
- **Command & Files**: bitrise_remote_machine_execute, bitrise_remote_machine_upload, bitrise_remote_machine_download
- **GUI Interaction**: bitrise_remote_machine_screenshot, bitrise_remote_machine_click, bitrise_remote_machine_mouse_drag, bitrise_remote_machine_type, bitrise_remote_machine_scroll
- **Remote Access**: bitrise_remote_machine_open_vnc

## Troubleshooting

### Authentication Issues

- **Token expired**: Generate a new Bitrise token

### Configuration Issues

- **Invalid JSON**: Validate your configuration:
    ```bash
    cat ~/.gemini/settings.json | jq .
    ```
- **MCP connection issues**: Check logs for connection errors:
    ```bash
    gemini --debug "test command"
    ```
- **Go not found**: Ensure Go is installed and in your PATH
