# Install Bitrise Remote Machine MCP Server in Cursor

## Prerequisites

1. [Cursor](https://cursor.com/download) IDE installed (latest version)
2. [Create a Bitrise API Token](https://devcenter.bitrise.io/api/authentication):
   - Go to your [Bitrise Account Settings/Security](https://app.bitrise.io/me/account/security).
   - Navigate to the "Personal access tokens" section.
   - Copy the generated token.
3. [Go](https://go.dev/) (>=1.25) installed

## Setup

The Bitrise Remote Machine MCP server runs locally via Go.

### Install Steps

1. Go directly to your global MCP configuration file at `~/.cursor/mcp.json` and enter the code block below
2. In Tools & Integrations > MCP tools, click the pencil icon next to "bitrise-remote-machine"
3. Replace `YOUR_BITRISE_PAT` with your actual [Bitrise Personal Access Token](https://devcenter.bitrise.io/api/authentication)
4. Save the file
5. Restart Cursor

### Configuration

```json
{
  "mcpServers": {
    "bitrise-remote-machine": {
      "command": "go",
      "args": [
        "run",
        "github.com/bitrise-io/bitrise-mcp-macos-remote-machine@latest"
      ],
      "env": {
        "BITRISE_TOKEN": "YOUR_BITRISE_PAT"
      }
    }
  }
}
```

## Configuration Files

- **Global (all projects)**: `~/.cursor/mcp.json`
- **Project-specific**: `.cursor/mcp.json` in project root

## Verification

1. Restart Cursor completely
2. Check for green dot in Settings → Tools & Integrations → MCP Tools
3. In chat/composer, check "Available Tools"
4. Test with: "Create a remote machine and list the files in /Users"

## Troubleshooting

### General Issues

- **MCP not loading**: Restart Cursor completely after configuration
- **Invalid JSON**: Validate that JSON format is correct
- **Tools not appearing**: Check server shows green dot in MCP settings
- **Go not found**: Ensure Go is installed and in your PATH
- **Check logs**: Look for MCP-related errors in Cursor logs

## Important Notes

- **Cursor specifics**: Supports both project and global configurations, uses `mcpServers` key
