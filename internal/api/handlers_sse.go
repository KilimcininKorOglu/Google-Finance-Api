package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/kilimcininkoroglu/google-finance-api/internal/decode"
	"github.com/kilimcininkoroglu/google-finance-api/internal/gfrpc"
)

var liveTickers = []string{
	"GOOGL:NASDAQ",
	"AAPL:NASDAQ",
	"MSFT:NASDAQ",
	"BTC-USD",
	"THYAO:IST",
	"USD-TRY",
	"EUR-TRY",
	"EUR-USD",
}

type liveQuote struct {
	Ticker        string  `json:"ticker"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
	Currency      string  `json:"currency"`
	Type          string  `json:"type"`
}

func (h *handlers) liveStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "streaming not supported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	sendQuotes := func() {
		quotes := h.fetchLiveQuotes(r.Context())
		if len(quotes) == 0 {
			return
		}
		data, err := json.Marshal(quotes)
		if err != nil {
			return
		}
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}

	sendQuotes()

	tick := time.NewTicker(15 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-tick.C:
			sendQuotes()
		}
	}
}

func (h *handlers) fetchLiveQuotes(ctx context.Context) []liveQuote {
	var (
		mu     sync.Mutex
		wg     sync.WaitGroup
		quotes []liveQuote
	)

	for _, ticker := range liveTickers {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()

			fetchCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
			defer cancel()

			tuple := gfrpc.TickerTuple(t)
			results, err := h.client.FetchTicker(fetchCtx, t, []gfrpc.RPCRequest{
				gfrpc.QuoteRequest(tuple),
			})
			if err != nil {
				return
			}

			raw, ok := results[gfrpc.MethodQuote]
			if !ok {
				return
			}

			q, err := decode.Quote(raw)
			if err != nil || q == nil {
				return
			}

			mu.Lock()
			quotes = append(quotes, liveQuote{
				Ticker:        q.Ticker,
				Name:          q.Name,
				Price:         q.Price,
				Change:        q.Change,
				ChangePercent: q.ChangePercent,
				Currency:      q.Currency,
				Type:          q.Type,
			})
			mu.Unlock()
		}(ticker)
	}

	wg.Wait()
	return quotes
}

func (h *handlers) sseQuotes(w http.ResponseWriter, r *http.Request) {
	quotes := h.fetchLiveQuotes(r.Context())
	if quotes == nil {
		quotes = []liveQuote{}
	}
	writeJSON(w, http.StatusOK, quotes)
}
