package gfrpc

import (
	"fmt"
	"regexp"
	"strings"
)

var validTicker = regexp.MustCompile(`^[A-Za-z0-9._]{1,20}:[A-Za-z0-9_]{1,20}$|^[A-Za-z0-9]{1,10}-[A-Za-z0-9]{1,10}$`)

func ValidateTicker(ticker string) error {
	if len(ticker) > 40 {
		return fmt.Errorf("ticker too long")
	}
	if !validTicker.MatchString(ticker) {
		return fmt.Errorf("invalid ticker format")
	}
	return nil
}

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
