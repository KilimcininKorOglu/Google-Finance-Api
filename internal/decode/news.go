package decode

import (
	"encoding/json"

	"github.com/kilimcininkoroglu/google-finance-api/internal/models"
)

func News(raw json.RawMessage) ([]models.NewsItem, error) {
	arr, err := unmarshalNested(raw)
	if err != nil {
		return nil, err
	}

	items := atSlice(arr, 0)
	if items == nil {
		return nil, nil
	}

	var news []models.NewsItem
	for _, item := range items {
		a, ok := item.([]any)
		if !ok || len(a) < 3 {
			continue
		}

		title := atString(a, 1)
		if title == "" {
			continue
		}

		n := models.NewsItem{
			URL:    atString(a, 0),
			Title:  title,
			Source: atString(a, 2),
		}

		if len(a) > 4 {
			n.Timestamp = atInt64(a, 4)
		}

		news = append(news, n)
	}

	return news, nil
}
