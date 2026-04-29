package api

import (
	"net/http"

	"github.com/kilimcininkoroglu/google-finance-api/internal/decode"
	"github.com/kilimcininkoroglu/google-finance-api/internal/gfrpc"
	"github.com/kilimcininkoroglu/google-finance-api/internal/models"
)

var chartRangeMap = map[string]int{
	"1D":  1,
	"5D":  2,
	"1M":  3,
	"6M":  4,
	"YTD": 5,
	"1Y":  6,
	"5Y":  7,
	"MAX": 8,
}

type handlers struct {
	client *gfrpc.Client
}

func (h *handlers) getQuote(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	if ticker == "" {
		writeError(w, http.StatusBadRequest, "ticker is required")
		return
	}
	if err := gfrpc.ValidateTicker(ticker); err != nil {
		writeError(w, http.StatusBadRequest, "invalid ticker: "+err.Error())
		return
	}

	tuple := gfrpc.TickerTuple(ticker)
	results, err := h.client.FetchTicker(r.Context(), ticker, []gfrpc.RPCRequest{
		gfrpc.QuoteRequest(tuple),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	raw, ok := results[gfrpc.MethodQuote]
	if !ok {
		writeError(w, http.StatusNotFound, "no quote data")
		return
	}

	quote, err := decode.Quote(raw)
	if err != nil || quote == nil {
		writeError(w, http.StatusNotFound, "could not decode quote")
		return
	}

	writeJSON(w, http.StatusOK, quote)
}

func (h *handlers) getCompany(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	if ticker == "" {
		writeError(w, http.StatusBadRequest, "ticker is required")
		return
	}
	if err := gfrpc.ValidateTicker(ticker); err != nil {
		writeError(w, http.StatusBadRequest, "invalid ticker: "+err.Error())
		return
	}

	tuple := gfrpc.TickerTuple(ticker)
	results, err := h.client.FetchTicker(r.Context(), ticker, []gfrpc.RPCRequest{
		gfrpc.CompanyRequest(tuple),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	raw, ok := results[gfrpc.MethodCompany]
	if !ok {
		writeError(w, http.StatusNotFound, "no company data")
		return
	}

	company, err := decode.Company(raw)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, company)
}

func (h *handlers) getChart(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	if ticker == "" {
		writeError(w, http.StatusBadRequest, "ticker is required")
		return
	}

	rangeStr := r.URL.Query().Get("range")
	if rangeStr == "" {
		rangeStr = "1M"
	}

	mode, ok := chartRangeMap[rangeStr]
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid range: use 1D, 5D, 1M, 6M, YTD, 1Y, 5Y, MAX")
		return
	}

	tuple := gfrpc.TickerTuple(ticker)
	results, err := h.client.FetchTicker(r.Context(), ticker, []gfrpc.RPCRequest{
		gfrpc.ChartRequest(tuple, mode),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	raw, ok := results[gfrpc.MethodChart]
	if !ok {
		writeError(w, http.StatusNotFound, "no chart data")
		return
	}

	chart, err := decode.Chart(raw)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, chart)
}

func (h *handlers) getNews(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	if ticker == "" {
		writeError(w, http.StatusBadRequest, "ticker is required")
		return
	}
	if err := gfrpc.ValidateTicker(ticker); err != nil {
		writeError(w, http.StatusBadRequest, "invalid ticker: "+err.Error())
		return
	}

	tuple := gfrpc.TickerTuple(ticker)
	isCrypto := gfrpc.IsCrypto(ticker)
	results, err := h.client.FetchTicker(r.Context(), ticker, []gfrpc.RPCRequest{
		gfrpc.NewsRequest(isCrypto, tuple),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	raw, ok := results[gfrpc.MethodNews]
	if !ok {
		writeError(w, http.StatusNotFound, "no news data")
		return
	}

	news, err := decode.News(raw)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, news)
}

func (h *handlers) getFinancials(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	if ticker == "" {
		writeError(w, http.StatusBadRequest, "ticker is required")
		return
	}
	if err := gfrpc.ValidateTicker(ticker); err != nil {
		writeError(w, http.StatusBadRequest, "invalid ticker: "+err.Error())
		return
	}

	tuple := gfrpc.TickerTuple(ticker)
	results, err := h.client.FetchTicker(r.Context(), ticker, []gfrpc.RPCRequest{
		gfrpc.FinancialsRequest(tuple),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	raw, ok := results[gfrpc.MethodFinancials]
	if !ok {
		writeError(w, http.StatusNotFound, "no financial data")
		return
	}

	periods, err := decode.Financials(raw)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	filterType := r.URL.Query().Get("type")
	switch filterType {
	case "annual":
		periods = filterPeriods(periods, true)
	case "quarterly":
		periods = filterPeriods(periods, false)
	}

	writeJSON(w, http.StatusOK, periods)
}

func (h *handlers) getRelated(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	if ticker == "" {
		writeError(w, http.StatusBadRequest, "ticker is required")
		return
	}
	if err := gfrpc.ValidateTicker(ticker); err != nil {
		writeError(w, http.StatusBadRequest, "invalid ticker: "+err.Error())
		return
	}

	tuple := gfrpc.TickerTuple(ticker)
	results, err := h.client.FetchTicker(r.Context(), ticker, []gfrpc.RPCRequest{
		gfrpc.RelatedRequest(tuple),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	raw, ok := results[gfrpc.MethodRelated]
	if !ok {
		writeError(w, http.StatusNotFound, "no related stocks data")
		return
	}

	stocks, err := decode.Related(raw)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, stocks)
}

func (h *handlers) getFull(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	if ticker == "" {
		writeError(w, http.StatusBadRequest, "ticker is required")
		return
	}

	rangeStr := r.URL.Query().Get("range")
	if rangeStr == "" {
		rangeStr = "1M"
	}
	mode, ok := chartRangeMap[rangeStr]
	if !ok {
		mode = 3
	}

	tuple := gfrpc.TickerTuple(ticker)
	isCrypto := gfrpc.IsCrypto(ticker)

	results, err := h.client.FetchTicker(r.Context(), ticker, []gfrpc.RPCRequest{
		gfrpc.QuoteRequest(tuple),
		gfrpc.CompanyRequest(tuple),
		gfrpc.ChartRequest(tuple, mode),
		gfrpc.NewsRequest(isCrypto, tuple),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	full := &models.FullQuote{}

	if raw, ok := results[gfrpc.MethodQuote]; ok {
		full.Quote, _ = decode.Quote(raw)
	}
	if raw, ok := results[gfrpc.MethodCompany]; ok {
		full.Company, _ = decode.Company(raw)
	}
	if raw, ok := results[gfrpc.MethodChart]; ok {
		full.Chart, _ = decode.Chart(raw)
	}
	if raw, ok := results[gfrpc.MethodNews]; ok {
		full.News, _ = decode.News(raw)
	}

	if full.Quote == nil {
		writeError(w, http.StatusNotFound, "no data for ticker")
		return
	}

	writeJSON(w, http.StatusOK, full)
}

func filterPeriods(periods []models.FinancialPeriod, annual bool) []models.FinancialPeriod {
	var filtered []models.FinancialPeriod
	for _, p := range periods {
		if p.IsAnnual == annual {
			filtered = append(filtered, p)
		}
	}
	return filtered
}
