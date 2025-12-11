package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bitrise-io/bitrise-mcp-remote-sandbox/internal/bitrise"
	"github.com/bitrise-io/bitrise-mcp-remote-sandbox/internal/tool"
	"github.com/jinzhu/configor"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const development = "development"

// BuildVersion is overwritten with go build flags.
var BuildVersion = development //nolint:gochecknoglobals

type config struct {
	// BitriseToken is the Bitrise API token used to authenticate requests.
	BitriseToken string `env:"BITRISE_TOKEN" required:"true"`
	// LogLevel is the log level for the application.
	LogLevel string `env:"LOG_LEVEL" default:"info"`
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("error: %+v", err)
	}
}

func run() error {
	var cfg config
	if err := configor.Load(&cfg); err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	logger, err := newStructuredLogger(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("initialize logger: %w", err)
	}

	toolBelt := tool.NewBelt()
	mcpServer := server.NewMCPServer(
		"bitrise",
		"2.0.0",
		server.WithRecovery(),
		server.WithToolCapabilities(false),
		server.WithLogging(),
	)
	toolBelt.RegisterAll(mcpServer)

	server.WithToolHandlerMiddleware(func(fn server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return fn(bitrise.ContextWithPAT(ctx, cfg.BitriseToken), request)
		}
	})(mcpServer)

	logger.Info("starting stdio transport")
	if err := server.ServeStdio(mcpServer); err != nil {
		return fmt.Errorf("serve stdio: %w", err)
	}
	return nil
}

func newStructuredLogger(level string) (*zap.SugaredLogger, error) {
	atom := zap.NewAtomicLevel()
	if err := atom.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("could parse log level: %w", err)
	}

	loggerConfig := zap.NewProductionConfig()
	if BuildVersion == development {
		loggerConfig = zap.NewDevelopmentConfig()
		loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		loggerConfig.DisableStacktrace = true
	}

	loggerConfig.OutputPaths = []string{"stderr"}
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	loggerConfig.Level = atom

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("new zap logger: %w", err)
	}
	return logger.Sugar(), nil
}
