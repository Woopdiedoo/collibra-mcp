package chip

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type contextKey int

const (
	callToolRequestKey contextKey = iota
	toolConfigKey
)

func SetToolConfig(ctx context.Context, toolConfig *ToolConfig) context.Context {
	return context.WithValue(ctx, toolConfigKey, toolConfig)
}

func SetCallToolRequest(ctx context.Context, toolRequest *mcp.CallToolRequest) context.Context {
	return context.WithValue(ctx, callToolRequestKey, toolRequest)
}

func GetCallToolRequest(ctx context.Context) (*mcp.CallToolRequest, error) {
	toolRequest, ok := ctx.Value(callToolRequestKey).(*mcp.CallToolRequest)
	if !ok || toolRequest == nil {
		return nil, errors.New("CallToolRequest not found in ctx")
	}
	return toolRequest, nil
}

func GetToolConfig(ctx context.Context) (*ToolConfig, error) {
	config, ok := ctx.Value(toolConfigKey).(*ToolConfig)
	if !ok || config == nil {
		return nil, errors.New("ToolConfig not found in ctx")
	}
	return config, nil
}

func GetSessionId(toolRequest *mcp.CallToolRequest) string {
	sessionId := toolRequest.GetSession().ID()
	if sessionId == "" {
		return uuid.New().String()
	}
	return sessionId
}
