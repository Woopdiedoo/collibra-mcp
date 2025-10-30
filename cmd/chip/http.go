package main

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	chip "github.com/collibra/chip/app"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func newCollibraClientFactory(config *chip.Config, transport http.RoundTripper) func(mcpRequest *mcp.CallToolRequest) *http.Client {
	return func(mcpRequest *mcp.CallToolRequest) *http.Client {
		client := &collibraClient{
			config:     config,
			mcpRequest: mcpRequest,
			next:       transport,
		}
		return &http.Client{
			Transport: client,
		}
	}
}

type collibraClient struct {
	config     *chip.Config
	mcpRequest *mcp.CallToolRequest
	next       http.RoundTripper
}

func (c *collibraClient) RoundTrip(request *http.Request) (*http.Response, error) {
	reqClone := request.Clone(request.Context())
	if c.config.Api.Url == "" {
		return nil, fmt.Errorf("API URL is not configured")
	}
	if c.config.Api.Username != "" && c.config.Api.Password != "" {
		reqClone.SetBasicAuth(c.config.Api.Username, c.config.Api.Password)
	} else {
		chip.CopyHeader(c.mcpRequest, reqClone, "Authorization")
	}
	reqClone.Header.Set("X-MCP-Session-Id", chip.GetSessionId(c.mcpRequest))
	baseURL, err := url.Parse(c.config.Api.Url)
	if err != nil {
		return nil, fmt.Errorf("invalid API URL configuration: %w", err)
	}
	reqClone.URL.Scheme = baseURL.Scheme
	reqClone.URL.Host = baseURL.Host
	reqClone.URL.Path = path.Join(baseURL.Path, request.URL.Path)
	return c.next.RoundTrip(reqClone)
}
