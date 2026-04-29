package api

import (
	"context"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/kilimcininkoroglu/google-finance-api/internal/gfrpc"
)

func NewServer(ctx context.Context, client *gfrpc.Client, port string, webFS fs.FS) *http.Server {
	hub := newLiveHub(client)
	go hub.run(ctx)

	h := &handlers{client: client, hub: hub}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", webHandler(webFS))
	mux.HandleFunc("GET /openapi.json", openAPIHandler(webFS))
	mux.HandleFunc("GET /robots.txt", staticFileHandler(webFS, "robots.txt", "text/plain; charset=utf-8", "public, max-age=3600"))
	mux.HandleFunc("GET /sitemap.xml", staticFileHandler(webFS, "sitemap.xml", "application/xml; charset=utf-8", "public, max-age=3600"))
	mux.HandleFunc("GET /llms.txt", staticFileHandler(webFS, "llms.txt", "text/plain; charset=utf-8", "public, max-age=3600"))

	mux.HandleFunc("GET /v1/quote/{ticker}", h.getQuote)
	mux.HandleFunc("GET /v1/company/{ticker}", h.getCompany)
	mux.HandleFunc("GET /v1/chart/{ticker}", h.getChart)
	mux.HandleFunc("GET /v1/news/{ticker}", h.getNews)
	mux.HandleFunc("GET /v1/financials/{ticker}", h.getFinancials)
	mux.HandleFunc("GET /v1/related/{ticker}", h.getRelated)
	mux.HandleFunc("GET /v1/full/{ticker}", h.getFull)

	mux.HandleFunc("GET /v1/market/indices", h.getMarketIndices)
	mux.HandleFunc("GET /v1/market/movers", h.getMarketMovers)
	mux.HandleFunc("GET /v1/market/trending", h.getTrending)
	mux.HandleFunc("GET /v1/market/earnings", h.getEarnings)
	mux.HandleFunc("GET /v1/market/headlines", h.getHeadlines)

	mux.HandleFunc("GET /v1/live", h.liveStream)
	mux.HandleFunc("GET /v1/live/snapshot", h.sseQuotes)

	version := "unknown"
	if v, err := fs.ReadFile(webFS, "VERSION"); err == nil {
		version = strings.TrimSpace(string(v))
	}

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "version": version})
	})

	var handler http.Handler = mux
	handler = loggingMiddleware(handler)
	handler = recoveryMiddleware(handler)
	handler = corsMiddleware(handler)

	return &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}
