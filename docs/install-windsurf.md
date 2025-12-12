# Install Bitrise Remote Machine MCP Server in Windsurf

## Prerequisites

1. [Windsurf IDE](https://windsurf.com/) installed (latest version)
2. [Create a Bitrise API Token](https://devcenter.bitrise.io/api/authentication):
   - Go to your [Bitrise Account Settings/Security](https://app.bitrise.io/me/account/security).
   - Navigate to the "Personal access tokens" section.
   - Copy the generated token.
3. [Go](https://go.dev/) (>=1.25) installed

## Setup

The Bitrise Remote Machine MCP server runs locally via Go.

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

## Installation Steps

### Manual Configuration

1. Click the hammer icon (🔨) in Cascade
2. Click **Configure** to open `~/.codeium/windsurf/mcp_config.json`
3. Add the configuration from above
4. Replace `YOUR_BITRISE_PAT` with your actual token
5. Save the file
6. Click **Refresh** (🔄) in the MCP toolbar

## Configuration Details

- **File path**: `~/.codeium/windsurf/mcp_config.json`
- **Scope**: Global configuration only (no per-project support)
- **Format**: Must be valid JSON (use a linter to verify)

## Verification

After installation:

1. Look for "1 available MCP server" in the MCP toolbar
2. Click the hammer icon to see available Bitrise Remote Machine tools
3. Test with: "Create a remote machine and list the files in /Users"
4. Check for green dot next to the server name

## Troubleshooting

### General Issues

- **Authentication failures**: Verify PAT hasn't expired
- **Invalid JSON**: Validate with [jsonlint.com](https://jsonlint.com)
- **Tools not appearing**: Restart Windsurf completely
- **Go not found**: Ensure Go is installed and in your PATH
- **Check logs**: `~/.codeium/windsurf/logs/`

## Important Notes

- **Windsurf limitations**: No environment variable interpolation, global config only