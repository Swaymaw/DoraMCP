package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Swaymaw/DoraMCP/logger"
	"github.com/Swaymaw/DoraMCP/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{Name: "dora-web-explorer", Version: "1.0.0"}, nil)
	mcp.AddTool(server, tools.WebSearchTool, tools.WebSearchHandler)
	mcp.AddTool(server, &tools.FetchTool, tools.FetchHandler)

	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, nil)

	mux := http.NewServeMux()
	mux.Handle("/mcp", cors(handler))

	logger.Log.Info("Starting MCP server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Log.Error(fmt.Sprintf("Server error: %v", err))
	}
	_ = context.Background()
}
