# Install Bitrise Remote Machine MCP Server in VS Code

## Prerequisites

1. [VS Code](https://code.visualstudio.com/Download) installed
2. [Create a Bitrise API Token](https://devcenter.bitrise.io/api/authentication):
   - Go to your [Bitrise Account Settings/Security](https://app.bitrise.io/me/account/security).
   - Navigate to the "Personal access tokens" section.
   - Copy the generated token.
3. [Go](https://go.dev/) (>=1.25) installed

## Setup

Follow [VS Code | Add an MCP server](https://code.visualstudio.com/docs/copilot/customization/mcp-servers#_add-an-mcp-server) and add the following configuration to your settings:

```json
{
  "servers": {
    "bitrise-remote-machine": {
      "type": "stdio",
      "command": "go",
      "args": [
        "run",
        "github.com/bitrise-io/bitrise-mcp-macos-remote-machine@latest"
      ],
      "env": {
        "BITRISE_TOKEN": "${input:bitrise-token}"
      }
    }
  },
  "inputs": [
    {
      "id": "bitrise-token",
      "type": "promptString",
      "description": "Bitrise token",
      "password": true
    }
  ]
}
```

Save the configuration. VS Code will automatically recognize the change and load the tools into Copilot Chat.

## Verification

1. Restart VS Code completely
2. Check for green dot in Settings → Tools & Integrations → MCP Tools
3. In chat/composer, check "Available Tools"
4. Test with: "Create a remote machine and list the files in /Users"

## Troubleshooting

- **MCP not loading**: Restart VS Code completely after configuration
- **Invalid JSON**: Validate that JSON format is correct
- **Tools not appearing**: Check server shows green dot in MCP settings
- **Go not found**: Ensure Go is installed and in your PATH
