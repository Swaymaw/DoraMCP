package tools

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Swaymaw/DoraMCP/logger"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type fetchInput struct {
	Url string `json:"url" jsonschema:"url to be fetched"`
}

var FetchTool = mcp.Tool{
	Name:        "fetch",
	Description: "Fetches a URL and returns the response body.",
}

type fetchOutput struct {
	Body string `json:"body"`
}

func fetchUrl(url string) (string, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Mimic a real browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	// Handle gzip/br encoding transparently
	var reader io.Reader
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return "", err
		}
	default:
		reader = resp.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func FetchHandler(ctx context.Context, req *mcp.CallToolRequest, args fetchInput) (*mcp.CallToolResult, fetchOutput, error) {
	body, err := fetchUrl(args.Url)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to fetch URL: %v", err))
		return nil, fetchOutput{}, err
	}
	logger.Log.Debug("Response body: " + body)
	logger.Log.Info("Successfully fetched URL: " + args.Url)
	return nil, fetchOutput{Body: body}, nil
}
