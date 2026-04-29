package api

import (
	"bytes"
	"io/fs"
	"net/http"
	"os"
)

const baseURLPlaceholder = "https://finance.hermestech.uk"

func webHandler(content fs.FS) http.HandlerFunc {
	raw, _ := fs.ReadFile(content, "index.html")
	envBase := os.Getenv("BASE_URL")

	return func(w http.ResponseWriter, r *http.Request) {
		baseURL := envBase
		if baseURL == "" {
			scheme := "https"
			if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
				scheme = proto
			}
			baseURL = scheme + "://" + r.Host
		}
		data := bytes.ReplaceAll(raw, []byte(baseURLPlaceholder), []byte(baseURL))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Link", `</openapi.json>; rel="describedby"; type="application/json"`)
		w.Write(data)
	}
}

func staticFileHandler(content fs.FS, filename, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := fs.ReadFile(content, filename)
		if err != nil {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		w.Header().Set("Content-Type", contentType)
		w.Write(data)
	}
}

func openAPIHandler(content fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := fs.ReadFile(content, "openapi.json")
		if err != nil {
			writeError(w, http.StatusNotFound, "openapi spec not found")
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(data)
	}
}
