# Bitrise Remote Machine MCP Server

MCP Server for Bitrise Remote Machines, enabling AI assistants to create and manage remote macOS virtual machines for executing commands, running builds, and transferring files.

## Features

- **VM Lifecycle Management**: Create, list, and delete remote macOS virtual machines.
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
| `bitrise_remote_machine_list` | List all running VMs |
| `bitrise_remote_machine_create` | Create a new macOS VM for remote execution |
| `bitrise_remote_machine_delete` | Terminate and delete a VM |

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

- **One VM at a time**: Users can only have one remote machine running
- **Auto-expiration**: VMs automatically terminate after 1 hour if not manually deleted
- **Boot time**: First command after creation may take longer while VM boots
- **Always check first**: Call `bitrise_remote_machine_list` before creating a new VM to reuse existing machines

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
