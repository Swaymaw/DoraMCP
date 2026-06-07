package main

import (
	"context"
	"net/http"

	"github.com/Swaymaw/DoraMCP/logger"
	"github.com/Swaymaw/DoraMCP/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{Name: "dora-web-explorer", Version: "1.0.0"}, nil)
	mcp.AddTool(server, tools.WebSearchTool, tools.WebSearchHandler)
	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, nil)

	mux := http.NewServeMux()
	mux.Handle("/mcp", handler)

	logger.Log.Info("Starting MCP server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Log.Error("Server error: %v", err)
	}
	_ = context.Background()
}
