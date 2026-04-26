package decode

import (
	"encoding/json"

	"github.com/kilimcininkoroglu/google-finance-api/internal/models"
)

var typeNames = map[int]string{
	0: "stock",
	1: "index",
	3: "crypto",
	5: "etf",
}

func Quote(raw json.RawMessage) (*models.Quote, error) {
	arr, err := unmarshalNested(raw)
	if err != nil {
		return nil, err
	}

	root := drillDown(arr, 0, 0, 0)
	if root == nil {
		return nil, nil
	}

	typeCode := int(atFloat(root, 3))
	typeName, ok := typeNames[typeCode]
	if !ok {
		typeName = "unknown"
	}

	isCrypto := typeCode == 3

	var ticker, exchange string
	if isCrypto {
		ticker = atString(root, 21)
	} else {
		tickerArr := atSlice(root, 1)
		if tickerArr != nil {
			ticker = atString(tickerArr, 0)
			exchange = atString(tickerArr, 1)
		}
	}

	priceArr := atSlice(root, 5)

	q := &models.Quote{
		Ticker:        ticker,
		Exchange:      exchange,
		Name:          atString(root, 2),
		Type:          typeName,
		Currency:      atString(root, 4),
		Timezone:      atString(root, 12),
		PreviousClose: atFloat(root, 7),
	}

	if priceArr != nil {
		q.Price = atFloat(priceArr, 0)
		q.Change = atFloat(priceArr, 1)
		q.ChangePercent = atFloat(priceArr, 2)
	}

	afterHoursArr := atSlice(root, 16)
	if afterHoursArr != nil && len(afterHoursArr) > 0 {
		q.AfterHours = &models.AfterHours{
			Price:         atFloat(afterHoursArr, 0),
			Change:        atFloat(afterHoursArr, 1),
			ChangePercent: atFloat(afterHoursArr, 2),
		}
	}

	return q, nil
}

func drillDown(arr []any, indices ...int) []any {
	current := arr
	for _, idx := range indices {
		next := atSlice(current, idx)
		if next == nil {
			return nil
		}
		current = next
	}
	return current
}
