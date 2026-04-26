package api

import (
	"net/http"
	"strconv"

	"github.com/kilimcininkoroglu/google-finance-api/internal/decode"
	"github.com/kilimcininkoroglu/google-finance-api/internal/gfrpc"
)

func (h *handlers) getMarketIndices(w http.ResponseWriter, r *http.Request) {
	results, err := h.client.FetchMarket(r.Context(), []gfrpc.RPCRequest{
		gfrpc.MarketIndicesRequest(),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	raw, ok := results[gfrpc.MethodMarketIndices]
	if !ok {
		writeError(w, http.StatusNotFound, "no indices data")
		return
	}

	indices, err := decode.MarketIndices(raw)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, indices)
}

func (h *handlers) getMarketMovers(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		category = "most-active"
	}

	count := 10
	if c := r.URL.Query().Get("count"); c != "" {
		if parsed, err := strconv.Atoi(c); err == nil && parsed > 0 {
			count = parsed
		}
	}

	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	results, err := h.client.FetchMarket(r.Context(), []gfrpc.RPCRequest{
		gfrpc.MarketMoversRequest([]string{category}, count, offset),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	raw, ok := results[gfrpc.MethodMarketMovers]
	if !ok {
		writeError(w, http.StatusNotFound, "no movers data")
		return
	}

	movers, err := decode.MarketMovers(raw)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, movers)
}

func (h *handlers) getTrending(w http.ResponseWriter, r *http.Request) {
	results, err := h.client.FetchMarket(r.Context(), []gfrpc.RPCRequest{
		gfrpc.TrendingRequest(),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	raw, ok := results[gfrpc.MethodTrending]
	if !ok {
		writeError(w, http.StatusNotFound, "no trending data")
		return
	}

	trending, err := decode.Trending(raw)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, trending)
}

func (h *handlers) getEarnings(w http.ResponseWriter, r *http.Request) {
	results, err := h.client.FetchMarket(r.Context(), []gfrpc.RPCRequest{
		gfrpc.EarningsRequest(),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	raw, ok := results[gfrpc.MethodEarnings]
	if !ok {
		writeError(w, http.StatusNotFound, "no earnings data")
		return
	}

	events, err := decode.Earnings(raw)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, events)
}

func (h *handlers) getHeadlines(w http.ResponseWriter, r *http.Request) {
	results, err := h.client.FetchMarket(r.Context(), []gfrpc.RPCRequest{
		gfrpc.TopHeadlineRequest(),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	raw, ok := results[gfrpc.MethodTopHeadline]
	if !ok {
		writeError(w, http.StatusNotFound, "no headline data")
		return
	}

	headline, err := decode.TopHeadline(raw)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, headline)
}
