package decode

import (
	"encoding/json"
	"fmt"

	"github.com/kilimcininkoroglu/google-finance-api/internal/models"
)

func Chart(raw json.RawMessage) (*models.ChartData, error) {
	arr, err := unmarshalNested(raw)
	if err != nil {
		return nil, err
	}

	chartRoot := drillDown(arr, 0, 0)
	if chartRoot == nil {
		return nil, nil
	}

	chart := &models.ChartData{
		PreviousClose: atFloat(chartRoot, 6),
	}

	periods := atSlice(chartRoot, 3)
	if periods == nil {
		return chart, nil
	}

	for _, p := range periods {
		period, ok := p.([]any)
		if !ok {
			continue
		}

		points := atSlice(period, 1)
		if points == nil {
			continue
		}

		for _, pt := range points {
			point, ok := pt.([]any)
			if !ok {
				continue
			}

			dateArr := atSlice(point, 0)
			priceArr := atSlice(point, 1)
			if dateArr == nil || priceArr == nil {
				continue
			}

			y := int(atFloat(dateArr, 0))
			m := int(atFloat(dateArr, 1))
			d := int(atFloat(dateArr, 2))

			cp := models.ChartPoint{
				Date:  fmt.Sprintf("%d-%02d-%02d", y, m, d),
				Price: atFloat(priceArr, 0),
			}

			if v := at(point, 2); v != nil {
				vol := int64(atFloat(point, 2))
				cp.Volume = &vol
			}

			chart.Points = append(chart.Points, cp)
		}
	}

	return chart, nil
}
