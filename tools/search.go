package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Swaymaw/DoraMCP/logger"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var WebSearchTool = &mcp.Tool{Name: "web-search", Description: "get web results "}

type searchInput struct {
	Query  string `json:"query" jsonschema:"query for the web search term"`
	MaxNum int    `json:"num,omitempty" jsonschmea:"max number of search terms to return"`
}

type searchResult struct {
	URL      string `json:"url"`
	Category string `json:"category"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type searchOutput struct {
	Results []searchResult `json:"results"`
}

func searxNGSearch(query string, max_count int) (searchOutput, error) {
	endpoint := fmt.Sprintf("http://127.0.0.1:9999/search?q=%s&format=json", url.QueryEscape(query))
	resp, err := http.Get(endpoint)

	if err != nil {
		logger.Log.Error("Cannot connect to DDG API -", err)
		return searchOutput{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error("Cannot read DDG API Response - %v", err)
		return searchOutput{}, err
	}

	var response searchOutput
	if err := json.Unmarshal(body, &response); err != nil {
		logger.Log.Error("Cannot parse SearxNG API Response - %v", err)
		return searchOutput{}, err
	}

	results := response.Results

	if max_count < len(results) {
		results = results[:max_count]
	}
	return response, nil
}

func WebSearchHandler(ctx context.Context, req *mcp.CallToolRequest, args searchInput) (*mcp.CallToolResult, searchOutput, error) {
	q := args.Query
	num := args.MaxNum

	response, err := searxNGSearch(q, num)

	if err != nil {
		logger.Log.Error("Search Failed for query: %s with Error - %v", q, err)
		return nil, searchOutput{}, nil
	}

	logger.Log.Info("Successfully Responded for query: " + q)

	return nil, response, nil
}
