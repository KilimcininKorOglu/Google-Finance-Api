package gfrpc

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type RPCRequest struct {
	ID     string
	Params []any
}

type RPCResult struct {
	ID   string
	Data json.RawMessage
}

func BuildBody(requests []RPCRequest) string {
	arr := make([][]any, len(requests))
	for i, r := range requests {
		paramJSON, _ := json.Marshal(r.Params)
		arr[i] = []any{r.ID, string(paramJSON), nil, fmt.Sprintf("%d", i+1)}
	}
	outer, _ := json.Marshal([]any{arr})
	return "f.req=" + url.QueryEscape(string(outer))
}

var hexLinePattern = regexp.MustCompile(`^[0-9a-fA-F]+$`)

func ParseResponse(raw string) ([]RPCResult, error) {
	stripped := strings.TrimPrefix(raw, ")]}'\n\n")
	stripped = strings.TrimPrefix(stripped, ")]}'\n")

	var results []RPCResult
	lines := strings.Split(stripped, "\n")

	for i := 0; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if !hexLinePattern.MatchString(trimmed) || i+1 >= len(lines) {
			continue
		}

		var entries []json.RawMessage
		if err := json.Unmarshal([]byte(lines[i+1]), &entries); err != nil {
			i++
			continue
		}

		for _, entry := range entries {
			var row []json.RawMessage
			if err := json.Unmarshal(entry, &row); err != nil || len(row) < 3 {
				continue
			}

			var marker string
			if err := json.Unmarshal(row[0], &marker); err != nil || marker != "wrb.fr" {
				continue
			}

			var rpcID string
			if err := json.Unmarshal(row[1], &rpcID); err != nil {
				continue
			}

			var dataStr string
			if err := json.Unmarshal(row[2], &dataStr); err != nil {
				results = append(results, RPCResult{ID: rpcID, Data: row[2]})
				continue
			}

			results = append(results, RPCResult{ID: rpcID, Data: json.RawMessage(dataStr)})
		}

		i++
	}

	return results, nil
}
