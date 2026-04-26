package gfrpc

import "strings"

func TickerTuple(ticker string) []any {
	if strings.Contains(ticker, "-") && !strings.Contains(ticker, ":") {
		parts := strings.SplitN(ticker, "-", 2)
		if len(parts) != 2 {
			return []any{nil, nil, []string{ticker}}
		}
		return []any{nil, nil, []string{parts[0], parts[1]}}
	}

	parts := strings.SplitN(ticker, ":", 2)
	if len(parts) != 2 {
		return []any{nil, []string{ticker}}
	}
	return []any{nil, []string{parts[0], parts[1]}}
}

func IsCrypto(ticker string) bool {
	return strings.Contains(ticker, "-") && !strings.Contains(ticker, ":")
}
