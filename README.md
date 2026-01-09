# Bitrise Remote Machine MCP Server

MCP Server for Bitrise Remote Machines, enabling AI assistants to create and manage remote macOS virtual machines for executing commands, running builds, and transferring files.

## Features

- **VM Lifecycle Management**: Create, list, start, stop, and delete remote macOS virtual machines. Users can have up to 5 machines with persistent disk state.
- **Command Execution**: Run shell commands on the VM (Xcode, Git, Homebrew, etc.).
- **File Transfer**: Upload local files/folders to the VM and download build artifacts.
- **GUI Automation**: Interact with the VM's graphical interface via screenshots, mouse clicks, keyboard input, and scrolling.
- **VNC Access**: Connect to the VM with a VNC client for full remote desktop access.

## Installation

- **[VS Code](/docs/install-vscode.md)** - Installation for VS Code IDE
- **[GitHub Copilot in other IDEs](/docs/install-other-copilot-ides.md)** - Installation for JetBrains, Visual Studio, Eclipse, and Xcode with GitHub Copilot
- **[Claude Applications](/docs/install-claude.md)** - Installation guide for Claude Desktop and Claude Code CLI
- **[Cursor](/docs/install-cursor.md)** - Installation guide for Cursor IDE
- **[Windsurf](/docs/install-windsurf.md)** - Installation guide for Windsurf IDE
- **[Gemini CLI](/docs/install-gemini-cli.md)** - Installation guide for Gemini CLI

## Available Tools

### VM Lifecycle

| Tool | Description |
|------|-------------|
| `bitrise_remote_machine_list` | List all VMs with their status and description |
| `bitrise_remote_machine_create` | Create a new macOS VM (with optional description) |
| `bitrise_remote_machine_start` | Start a terminated or failed VM |
| `bitrise_remote_machine_stop` | Stop a running VM (preserves disk state) |
| `bitrise_remote_machine_delete` | Permanently delete a VM |
| `bitrise_remote_machine_update_description` | Update a VM's description |

### Command & File Operations

| Tool | Description |
|------|-------------|
| `bitrise_remote_machine_execute` | Run shell commands on the VM using `bash -c` |
| `bitrise_remote_machine_upload` | Upload local files/folders to the VM |
| `bitrise_remote_machine_download` | Download files/folders from the VM |

### GUI Interaction

| Tool | Description |
|------|-------------|
| `bitrise_remote_machine_screenshot` | Capture the current VM display (1024x768 resolution) |
| `bitrise_remote_machine_click` | Simulate mouse clicks at specified coordinates (left/right/middle) |
| `bitrise_remote_machine_mouse_drag` | Simulate mouse drag operations between two points |
| `bitrise_remote_machine_type` | Simulate keyboard input (supports control characters: \n, \t, \b, \e) |
| `bitrise_remote_machine_scroll` | Scroll within the VM GUI (up/down in line units) |

### Remote Access

| Tool | Description |
|------|-------------|
| `bitrise_remote_machine_open_vnc` | Get VNC credentials for graphical remote desktop access |

## Usage Notes

### VM Management

- **Multi-machine support**: Users can have up to 5 machines total (in any state)
- **One running at a time**: Only one machine can be in running or pending state
- **Persistent disk state**: Stopped (terminated) machines preserve their disk state and can be restarted
- **Auto-expiration**: VMs automatically terminate after 1 hour if not manually stopped
- **Boot time**: Creating or starting a VM takes 30-60 seconds
- **Always check first**: Call `bitrise_remote_machine_list` before creating to reuse existing machines
- **Use descriptions**: Always provide descriptions when creating machines to identify them later

### Command Execution

- **Bash commands**: Commands run via `bash -c`, supporting pipes, redirects, and command chaining
- **Terminating commands required**: Commands must exit cleanly; avoid backgrounding, infinite loops, or interactive commands
- **No file transfers**: Do not use execute for file transfers; use the dedicated upload/download tools instead

### File Transfer

- **Upload**: Local files/folders are automatically compressed to tar.gz and extracted on the VM
- **Download**: Files/folders are extracted from tar.gz automatically on your local machine

### Screen Resolution

- **Standard resolution**: All GUI operations use 1024x768 pixel resolution on the remote machine
- **Coordinate system**: Coordinates for clicks/drags are absolute (0-1023 for x, 0-767 for y)
