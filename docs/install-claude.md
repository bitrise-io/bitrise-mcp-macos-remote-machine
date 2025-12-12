# Install Bitrise Remote Machine MCP Server in Claude Applications

## Claude Code CLI

### Prerequisites

- Claude Code CLI installed
- [Create a Bitrise API Token](https://devcenter.bitrise.io/api/authentication):
   - Go to your [Bitrise Account Settings/Security](https://app.bitrise.io/me/account/security).
   - Navigate to the "Personal access tokens" section.
   - Copy the generated token.
- [Go](https://go.dev/) (>=1.25) installed
- Open Claude Code inside the directory for your project (recommended for best experience and clear scope of configuration)

<details>
<summary><b>Storing Your PAT Securely</b></summary>
<br>

For security, avoid hardcoding your token. One common approach:

1. Store your token in `.env` file
```
BITRISE_PAT=your_token_here
```

2. Add to .gitignore
```bash
echo -e ".env\n.mcp.json" >> .gitignore
```

</details>

### Setup

1. Run the following command in the Claude Code CLI:
```bash
claude mcp add bitrise-remote-machine -e BITRISE_TOKEN=YOUR_BITRISE_PAT -- go run github.com/bitrise-io/bitrise-mcp-macos-remote-machine@latest
```

With an environment variable:
```bash
claude mcp add bitrise-remote-machine -e BITRISE_TOKEN=$(grep BITRISE_PAT .env | cut -d '=' -f2) -- go run github.com/bitrise-io/bitrise-mcp-macos-remote-machine@latest
```

2. Restart Claude Code
3. Run `claude mcp list` to see if the Bitrise Remote Machine server is configured

### Verification

```bash
claude mcp list
claude mcp get bitrise-remote-machine
```

## Claude Desktop

### Prerequisites

- Claude Desktop installed (latest version)
- [Create a Bitrise API Token](https://devcenter.bitrise.io/api/authentication):
   - Go to your [Bitrise Account Settings/Security](https://app.bitrise.io/me/account/security).
   - Navigate to the "Personal access tokens" section.
   - Copy the generated token.
- [Go](https://go.dev/) (>=1.25) installed

### Configuration File Location

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

### Setup

Add this codeblock to your `claude_desktop_config.json`:

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
        "BITRISE_TOKEN": "YOUR_BITRISE_PAT",
        "PATH": "PATH to bin directory of go:PATH to directory of git",
        "GOPATH": "your GOPATH",
        "GOCACHE": "your GOCACHE"
      }
    }
  }
}
```

### Manual Setup Steps

1. Open Claude Desktop
2. Go to Settings → Developer → Edit Config
3. Paste the code block above in your configuration file
4. If you're navigating to the configuration file outside of the app:
   - **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
5. Open the file in a text editor
6. Paste the code block above
7. Replace `YOUR_BITRISE_PAT` with your actual token
8. Save the file
9. Restart Claude Desktop

## Troubleshooting

**Authentication Failed:**
- Check token hasn't expired

**Server Not Starting / Tools Not Showing:**
- Run `claude mcp list` to view currently configured MCP servers
- Validate JSON syntax
- If using an environment variable to store your PAT, make sure you're properly sourcing your PAT using the environment variable
- Restart Claude Code and check `/mcp` command
- Delete the server by running `claude mcp remove bitrise-remote-machine` and repeating the setup process
- Make sure you're running Claude Code within the project you're currently working on to ensure the MCP configuration is properly scoped to your project
- Check logs:
  - Claude Code: Use `/mcp` command
  - Claude Desktop: `ls ~/Library/Logs/Claude/` and `cat ~/Library/Logs/Claude/mcp-server-*.log` (macOS) or `%APPDATA%\Claude\logs\` (Windows)

## Important Notes

- Configuration scopes for Claude Code:
  - `-s user`: Available across all projects
  - `-s project`: Shared via `.mcp.json` file
  - Default: `local` (current project only)
