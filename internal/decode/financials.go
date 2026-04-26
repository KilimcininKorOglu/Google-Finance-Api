package decode

import (
	"encoding/json"
	"fmt"

	"github.com/kilimcininkoroglu/google-finance-api/internal/models"
)

func Financials(raw json.RawMessage) ([]models.FinancialPeriod, error) {
	arr, err := unmarshalNested(raw)
	if err != nil {
		return nil, err
	}

	root := drillDown(arr, 0, 0)
	if root == nil {
		return nil, nil
	}

	var periods []models.FinancialPeriod

	quarterly := atSlice(root, 0)
	for _, entry := range quarterly {
		e, ok := entry.([]any)
		if !ok || len(e) < 3 {
			continue
		}
		year := int(atFloat(e, 0))
		quarter := int(atFloat(e, 1))
		data := atSlice(e, 2)
		if data == nil {
			continue
		}
		fp := decodeFinancialData(data, false)
		fp.FiscalEnd = fmt.Sprintf("%d-Q%d", year, quarter)
		periods = append(periods, fp)
	}

	annual := atSlice(root, 1)
	for _, entry := range annual {
		e, ok := entry.([]any)
		if !ok || len(e) < 2 {
			continue
		}
		year := int(atFloat(e, 0))
		data := atSlice(e, 1)
		if data == nil {
			continue
		}
		fp := decodeFinancialData(data, true)
		fp.FiscalEnd = fmt.Sprintf("%d", year)
		periods = append(periods, fp)
	}

	return periods, nil
}

func decodeFinancialData(d []any, isAnnual bool) models.FinancialPeriod {
	fp := models.FinancialPeriod{
		IsAnnual:        isAnnual,
		Revenue:         atFloat(d, 0),
		NetIncome:       atFloat(d, 1),
		EPS:             atFloat(d, 2),
		OperatingMargin: atFloat(d, 3),
		OperatingIncome: atFloat(d, 4),
	}

	if len(d) > 7 {
		fp.EBITDA = atFloat(d, 7)
	}
	if len(d) > 8 {
		fp.SharesOutstanding = atFloat(d, 8)
	}
	if len(d) > 9 {
		fp.EPSDiluted = atFloat(d, 9)
	}
	if len(d) > 11 {
		fp.RevenueGrowthYoY = atFloat(d, 11)
	}
	if len(d) > 16 {
		fp.Currency = atString(d, 16)
	}
	if len(d) > 17 {
		dateArr := atSlice(d, 17)
		if dateArr != nil && len(dateArr) >= 3 {
			y := int(atFloat(dateArr, 0))
			m := int(atFloat(dateArr, 1))
			dd := int(atFloat(dateArr, 2))
			fp.FiscalEnd = fmt.Sprintf("%d-%02d-%02d", y, m, dd)
		}
	}
	if len(d) > 18 {
		fp.PERatio = atFloat(d, 18)
	}
	if len(d) > 19 {
		fp.TotalAssets = atFloat(d, 19)
	}
	if len(d) > 20 {
		fp.TotalLiabilities = atFloat(d, 20)
	}
	if len(d) > 21 {
		fp.TotalEquity = atFloat(d, 21)
	}
	if len(d) > 24 {
		fp.OperatingCashFlow = atFloat(d, 24)
	}
	if len(d) > 29 {
		fp.ProfitMargin = atFloat(d, 29)
	}
	if len(d) > 32 {
		fp.FreeCashFlow = atFloat(d, 32)
	}
	if len(d) > 34 {
		fp.CapEx = atFloat(d, 34)
	}

	return fp
}
