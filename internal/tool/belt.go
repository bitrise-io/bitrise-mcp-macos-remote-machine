package tool

import (
	"github.com/bitrise-io/bitrise-mcp-remote-sandbox/internal/bitrise"
	"github.com/mark3labs/mcp-go/server"
)

type Belt struct {
	tools map[string]bitrise.Tool
}

func NewBelt() *Belt {
	var toolList = []bitrise.Tool{
		ListRemoteMachines,
		CreateRemoteMachine,
		DeleteRemoteMachine,
		ExecuteCommand,
		Upload,
		Download,
		OpenVNC,
		Click,
		MouseDrag,
		Screenshot,
		Scroll,
		Type,
	}
	belt := &Belt{tools: make(map[string]bitrise.Tool)}
	for _, tool := range toolList {
		belt.tools[tool.Definition.Name] = tool
	}
	return belt
}

func (b *Belt) RegisterAll(server *server.MCPServer) {
	for _, tool := range b.tools {
		server.AddTool(tool.Definition, tool.Handler)
	}
}
