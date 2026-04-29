package gfrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const maxResponseSize = 10 << 20 // 10 MB

const batchExecuteURL = "https://www.google.com/finance/_/GoogleFinanceUi/data/batchexecute"

var defaultHeaders = map[string]string{
	"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
	"Cookie":       "CONSENT=YES+cb",
	"Origin":       "https://www.google.com",
	"Referer":      "https://www.google.com/finance/",
}

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) Execute(ctx context.Context, sourcePath string, requests []RPCRequest) (map[string]json.RawMessage, error) {
	rpcIDs := uniqueRPCIDs(requests)
	reqURL := fmt.Sprintf("%s?rpcids=%s&source-path=%s&hl=en&gl=us&rt=c",
		batchExecuteURL, strings.Join(rpcIDs, ","), url.QueryEscape(sourcePath))

	body := BuildBody(requests)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, strings.NewReader(body))
	if err != nil {
		log.Printf("gfrpc: create request error: %v", err)
		return nil, fmt.Errorf("upstream service unavailable")
	}

	for k, v := range defaultHeaders {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("gfrpc: execute error: %v", err)
		return nil, fmt.Errorf("upstream service unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("gfrpc: unexpected status %d for %s", resp.StatusCode, sourcePath)
		return nil, fmt.Errorf("upstream service error")
	}

	rawBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		log.Printf("gfrpc: read error: %v", err)
		return nil, fmt.Errorf("upstream service error")
	}

	results, err := ParseResponse(string(rawBody))
	if err != nil {
		log.Printf("gfrpc: parse error: %v", err)
		return nil, fmt.Errorf("data not available")
	}

	resultMap := make(map[string]json.RawMessage, len(results))
	for _, r := range results {
		resultMap[r.ID] = r.Data
	}

	return resultMap, nil
}

func (c *Client) FetchTicker(ctx context.Context, ticker string, requests []RPCRequest) (map[string]json.RawMessage, error) {
	sourcePath := fmt.Sprintf("/finance/quote/%s", ticker)
	return c.Execute(ctx, sourcePath, requests)
}

func (c *Client) FetchMarket(ctx context.Context, requests []RPCRequest) (map[string]json.RawMessage, error) {
	return c.Execute(ctx, "/finance/", requests)
}

func uniqueRPCIDs(requests []RPCRequest) []string {
	seen := make(map[string]bool, len(requests))
	var ids []string
	for _, r := range requests {
		if !seen[r.ID] {
			seen[r.ID] = true
			ids = append(ids, r.ID)
		}
	}
	return ids
}
