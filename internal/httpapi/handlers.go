package httpapi

import (
	"context"
	"encoding/json"
	"net/http"

	"config-audit/internal/app"
)

const maxRequestBodyBytes = 5 << 20

func NewHandler(analyzer app.Analyzer) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("POST /v1/analyze", analyzeHandler(analyzer))
	return mux
}

func healthHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write([]byte(`{"status":"ok"}` + "\n"))
}

func analyzeHandler(analyzer app.Analyzer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		source := request.URL.Query().Get("filename")
		if source == "" {
			source = "request-body"
		}

		body := http.MaxBytesReader(writer, request.Body, maxRequestBodyBytes)
		defer body.Close()

		findings, err := analyzer.AnalyzeReader(context.Background(), body, source)
		if err != nil {
			writeError(writer, http.StatusBadRequest, err)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(findings); err != nil {
			writeError(writer, http.StatusInternalServerError, err)
			return
		}
	}
}

func writeError(writer http.ResponseWriter, status int, err error) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(map[string]string{
		"error": err.Error(),
	})
}
