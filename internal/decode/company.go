package decode

import (
	"encoding/json"

	"github.com/kilimcininkoroglu/google-finance-api/internal/models"
)

func Company(raw json.RawMessage) (*models.CompanyInfo, error) {
	arr, err := unmarshalNested(raw)
	if err != nil {
		return nil, err
	}

	info := drillDown(arr, 0, 0)
	if info == nil {
		return nil, nil
	}

	desc := atString(info, 2)
	if desc == "" {
		return nil, nil
	}

	return &models.CompanyInfo{
		Description:      desc,
		CEO:              atString(info, 5),
		Employees:        atInt64(info, 6),
		MarketCap:        atFloat(info, 7),
		Open:             atFloat(info, 9),
		High:             atFloat(info, 10),
		Low:              atFloat(info, 11),
		FiftyTwoWeekHigh: atFloat(info, 12),
		FiftyTwoWeekLow:  atFloat(info, 13),
		PERatio:          atFloat(info, 16),
		Volume:           atInt64(info, 18),
		Sector:           atString(info, 71),
	}, nil
}
