package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/collibra/chip/pkg/auth"
	"github.com/collibra/chip/pkg/chip"
	"github.com/collibra/chip/pkg/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	config := Init()

	slog.Info(fmt.Sprintf("Starting Collibra MCP server (version: %s)...", chip.Version))

	if config.Api.Url == "" {
		slog.Error("Missing Api url")
		os.Exit(1)
	}

	// Handle SSO authentication if enabled
	if config.Api.SSOAuth {
		if err := handleSSOAuthentication(config); err != nil {
			slog.Error(fmt.Sprintf("SSO authentication failed: %v", err))
			os.Exit(1)
		}
	}

	if config.Api.Username != "" && config.Api.Password != "" {
		slog.Warn("Using a single basic auth header for all requests is not recommended as it will result in all actions being attributed to the same account. Consider setting an appropriate basic auth header for each request.")
	}

	client := newCollibraClient(config)
	server := chip.NewServer(chip.WithToolMiddleware(chip.ToolMiddlewareFunc(setCollibraHost(config.Api.Url))))
	toolConfig := &chip.ToolConfig{
		EnabledTools:  config.Mcp.EnabledTools,
		DisabledTools: config.Mcp.DisabledTools,
	}
	tools.RegisterAll(server, client, toolConfig)

	if config.Mcp.Mode == "stdio" {
		runStdioServer(server)
	} else if strings.HasPrefix(config.Mcp.Mode, "http") {
		runHttpServer(config.Mcp.Mode, server, config.Mcp.Http.Port)
	} else {
		slog.Error(fmt.Sprintf("Invalid server mode: '%s'", config.Mcp.Mode))
		os.Exit(1)
	}
}

// handleSSOAuthentication performs browser-based SSO authentication
func handleSSOAuthentication(config *Config) error {
	cache := auth.NewSessionCache(config.Api.SSOCachePath)

	// Try to load cached session first
	cached, err := cache.Load(config.Api.Url)
	if err != nil {
		slog.Warn(fmt.Sprintf("Failed to load cached session: %v", err))
	}

	if cached != nil {
		// Use cached session
		config.Api.Cookie = fmt.Sprintf("JSESSIONID=%s", cached.Cookie)
		slog.Info(fmt.Sprintf("Using cached SSO session (expires: %s)", cached.ExpiresAt.Format(time.RFC3339)))
		return nil
	}

	// No valid cached session, perform SSO authentication
	timeout := time.Duration(config.Api.SSOTimeout) * time.Second
	result, err := auth.AuthenticateWithSSO(context.Background(), config.Api.Url, timeout)
	if err != nil {
		return err
	}

	// Set the cookie in config
	config.Api.Cookie = fmt.Sprintf("JSESSIONID=%s", result.Cookie)

	// Cache the session for future use
	if err := cache.Save(result, config.Api.Url); err != nil {
		slog.Warn(fmt.Sprintf("Failed to cache session: %v", err))
	}

	return nil
}

func runStdioServer(server *chip.Server) {
	slog.Info("Listening on stdio")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		slog.Error(fmt.Sprintf("Failed to run stdio server: %v", err))
		os.Exit(1)
	}
}

func runHttpServer(mode string, server *chip.Server, port int) {
	var handler http.Handler

	switch mode {
	case "http", "http-streamable":
		slog.Info("Using streamable http handler")
		handler = mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
			return &server.Server
		}, &mcp.StreamableHTTPOptions{
			Stateless: true,
		})
	case "http-sse":
		slog.Info("Using SSE http handler")
		handler = mcp.NewSSEHandler(func(req *http.Request) *mcp.Server {
			return &server.Server
		}, &mcp.SSEOptions{})
	default:
		slog.Error(fmt.Sprintf("Invalid HTTP mode: %s (must be 'http', 'http-sse' or 'http-streamable')", mode))
		os.Exit(1)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		Handler: handler,
	}

	slog.Warn("HTTP server is only listening on localhost for security reasons.")
	slog.Info(fmt.Sprintf("Listening on localhost:%d", port))
	if err := httpServer.ListenAndServe(); err != nil {
		slog.Error(fmt.Sprintf("Failed to start HTTP server: %v", err))
		os.Exit(1)
	}
}

func setCollibraHost(collibraHost string) func(ctx context.Context, toolRequest *mcp.CallToolRequest, next chip.CallToolFunc) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, toolRequest *mcp.CallToolRequest, next chip.CallToolFunc) (*mcp.CallToolResult, error) {
		ctx = chip.SetCollibraHost(ctx, collibraHost)
		slog.InfoContext(ctx, fmt.Sprintf("Calling tool: %s", toolRequest.Params.Name), "tool_name", toolRequest.Params.Name)
		return next(ctx, toolRequest)
	}
}
