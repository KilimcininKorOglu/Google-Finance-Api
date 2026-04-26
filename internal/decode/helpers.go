package decode

import (
	"encoding/json"
	"fmt"
)

func unmarshalNested(raw json.RawMessage) ([]any, error) {
	var arr []any
	if err := json.Unmarshal(raw, &arr); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return arr, nil
}

func at(arr []any, idx int) any {
	if idx < 0 || idx >= len(arr) {
		return nil
	}
	return arr[idx]
}

func atSlice(arr []any, idx int) []any {
	v := at(arr, idx)
	if v == nil {
		return nil
	}
	if s, ok := v.([]any); ok {
		return s
	}
	return nil
}

func atFloat(arr []any, idx int) float64 {
	v := at(arr, idx)
	if v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return n
	case json.Number:
		f, _ := n.Float64()
		return f
	default:
		return 0
	}
}

func atInt64(arr []any, idx int) int64 {
	v := at(arr, idx)
	if v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int64(n)
	case json.Number:
		i, _ := n.Int64()
		return i
	default:
		return 0
	}
}

func atString(arr []any, idx int) string {
	v := at(arr, idx)
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func atStringSlice(arr []any, idx int) []string {
	s := atSlice(arr, idx)
	if s == nil {
		return nil
	}
	result := make([]string, 0, len(s))
	for _, v := range s {
		if str, ok := v.(string); ok {
			result = append(result, str)
		}
	}
	return result
}
