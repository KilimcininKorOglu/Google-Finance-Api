package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kilimcininkoroglu/google-finance-api/internal/decode"
	"github.com/kilimcininkoroglu/google-finance-api/internal/gfrpc"
)

var (
	sseConnections atomic.Int32
	maxSSEConns    int32 = 50
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

type liveHub struct {
	client *gfrpc.Client
	mu     sync.RWMutex
	quotes []liveQuote
	notify chan struct{}
}

func newLiveHub(client *gfrpc.Client) *liveHub {
	return &liveHub{
		client: client,
		notify: make(chan struct{}),
	}
}

func (hub *liveHub) run(ctx context.Context) {
	hub.fetch()

	tick := time.NewTicker(15 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			hub.fetch()
		}
	}
}

func (hub *liveHub) fetch() {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	var (
		mu     sync.Mutex
		wg     sync.WaitGroup
		quotes []liveQuote
	)

	for _, ticker := range liveTickers {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()

			fetchCtx, fetchCancel := context.WithTimeout(ctx, 8*time.Second)
			defer fetchCancel()

			tuple := gfrpc.TickerTuple(t)
			results, err := hub.client.FetchTicker(fetchCtx, t, []gfrpc.RPCRequest{
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

	if len(quotes) > 0 {
		hub.mu.Lock()
		hub.quotes = quotes
		hub.mu.Unlock()

		ch := hub.notify
		hub.notify = make(chan struct{})
		close(ch)
	}
}

func (hub *liveHub) snapshot() []liveQuote {
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	result := make([]liveQuote, len(hub.quotes))
	copy(result, hub.quotes)
	return result
}

func (h *handlers) liveStream(w http.ResponseWriter, r *http.Request) {
	if sseConnections.Load() >= maxSSEConns {
		writeError(w, http.StatusServiceUnavailable, "too many live connections")
		return
	}
	sseConnections.Add(1)
	defer sseConnections.Add(-1)

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "streaming not supported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	send := func() {
		quotes := h.hub.snapshot()
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

	send()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-h.hub.notify:
			send()
		}
	}
}

func (h *handlers) sseQuotes(w http.ResponseWriter, r *http.Request) {
	quotes := h.hub.snapshot()
	if quotes == nil {
		quotes = []liveQuote{}
	}
	writeJSON(w, http.StatusOK, quotes)
}
