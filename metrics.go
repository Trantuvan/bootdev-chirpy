package main

import (
	"html/template"
	"net/http"
)

const metricsHtml string = "templates/metrics.html"

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf8")
	w.WriteHeader(http.StatusOK)
	html := template.Must(template.ParseFiles(metricsHtml))
	err := html.Execute(w, cfg.fileserverHits.Load())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		w.Header().Add("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		next.ServeHTTP(w, r)
	})
}
