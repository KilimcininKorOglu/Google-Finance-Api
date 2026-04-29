package api

import (
	"io/fs"
	"net/http"
)

func webHandler(content fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := fs.ReadFile(content, "index.html")
		if err != nil {
			writeError(w, http.StatusInternalServerError, "page not found")
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
