package decode

import (
	"encoding/json"
	"fmt"

	"github.com/kilimcininkoroglu/google-finance-api/internal/models"
)

func MarketIndices(raw json.RawMessage) ([]models.MarketIndex, error) {
	arr, err := unmarshalNested(raw)
	if err != nil {
		return nil, err
	}

	var indices []models.MarketIndex
	findQuotableItems(arr, func(a []any) {
		mi := extractMarketItem(a)
		indices = append(indices, models.MarketIndex{
			Ticker:        mi.ticker,
			Name:          mi.name,
			Price:         mi.price,
			Change:        mi.change,
			ChangePercent: mi.changePercent,
		})
	})

	return indices, nil
}

func isQuotableItem(a []any) bool {
	if len(a) < 6 {
		return false
	}
	tickerArr, ok1 := a[1].([]any)
	priceArr, ok2 := a[5].([]any)
	if !ok1 || !ok2 || len(tickerArr) < 2 || len(priceArr) < 1 {
		return false
	}
	_, isStr := tickerArr[0].(string)
	return isStr
}

func findQuotableItems(obj []any, fn func([]any)) {
	if isQuotableItem(obj) {
		fn(obj)
		return
	}
	for _, item := range obj {
		if child, ok := item.([]any); ok {
			findQuotableItems(child, fn)
		}
	}
}

func MarketMovers(raw json.RawMessage) ([]models.MarketMover, error) {
	arr, err := unmarshalNested(raw)
	if err != nil {
		return nil, err
	}

	var movers []models.MarketMover
	findQuotableItems(arr, func(a []any) {
		mi := extractMarketItem(a)
		movers = append(movers, models.MarketMover{
			Ticker:        mi.ticker,
			Name:          mi.name,
			Price:         mi.price,
			Change:        mi.change,
			ChangePercent: mi.changePercent,
		})
	})

	return movers, nil
}

func Trending(raw json.RawMessage) ([]models.MarketMover, error) {
	return MarketMovers(raw)
}

func Earnings(raw json.RawMessage) ([]models.EarningsEvent, error) {
	arr, err := unmarshalNested(raw)
	if err != nil {
		return nil, err
	}

	items := atSlice(arr, 0)
	if items == nil {
		return nil, nil
	}

	var events []models.EarningsEvent
	for _, item := range items {
		a, ok := item.([]any)
		if !ok || len(a) < 4 {
			continue
		}

		ev := models.EarningsEvent{
			Name: atString(a, 3),
		}

		tickerArr := atSlice(a, 1)
		if tickerArr != nil && len(tickerArr) >= 2 {
			ev.Ticker = atString(tickerArr, 0)
			ev.Exchange = atString(tickerArr, 1)
		}

		if len(a) > 4 {
			dateInfo := atSlice(a, 4)
			if dateInfo != nil && len(dateInfo) >= 2 {
				year := int(atFloat(dateInfo, 0))
				quarter := int(atFloat(dateInfo, 1))
				ev.Date = fmt.Sprintf("%d-Q%d", year, quarter)
			}
		}

		events = append(events, ev)
	}

	return events, nil
}

func TopHeadline(raw json.RawMessage) (*models.Headline, error) {
	arr, err := unmarshalNested(raw)
	if err != nil {
		return nil, err
	}

	root := atSlice(arr, 0)
	if root == nil {
		return nil, nil
	}

	return &models.Headline{
		Title:  atString(root, 1),
		URL:    atString(root, 0),
		Source: atString(root, 2),
	}, nil
}

func Related(raw json.RawMessage) ([]models.RelatedStock, error) {
	arr, err := unmarshalNested(raw)
	if err != nil {
		return nil, err
	}

	var stocks []models.RelatedStock
	findQuotableItems(arr, func(a []any) {
		mi := extractMarketItem(a)
		stocks = append(stocks, models.RelatedStock{
			Ticker:        mi.ticker,
			Name:          mi.name,
			Price:         mi.price,
			Change:        mi.change,
			ChangePercent: mi.changePercent,
		})
	})

	return stocks, nil
}

type marketItem struct {
	ticker, name                   string
	price, change, changePercent float64
}

func extractMarketItem(a []any) marketItem {
	mi := marketItem{
		name: atString(a, 2),
	}

	tickerArr := atSlice(a, 1)
	if tickerArr != nil && len(tickerArr) >= 2 {
		mi.ticker = atString(tickerArr, 0) + ":" + atString(tickerArr, 1)
	}

	priceArr := atSlice(a, 5)
	if priceArr != nil {
		mi.price = atFloat(priceArr, 0)
		mi.change = atFloat(priceArr, 1)
		mi.changePercent = atFloat(priceArr, 2)
	}

	return mi
}
