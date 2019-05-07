package main

import (
	"context"
	"encoding/json"
	_ "expvar"
	"net/http"
	"os"
	"time"

	"github.com/koenbollen/kramp/sources"
	"github.com/rs/cors"
	"github.com/unrolled/secure"
	"go.uber.org/zap"
)

type response struct {
	Data  []sources.Result `json:"data"`
	Took  time.Duration    `json:"took"`
	Query string           `json:"query"`
	Count int              `json:"count"`
}

func main() {
	logger, _ := zap.NewProduction()

	source := sources.NewFromEnvironment(logger)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger := logger.With(
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		logger.Debug("request started")

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		err := r.ParseForm()
		if err != nil {
			logger.Error("failed to parse form", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		query := r.Form.Get("q")
		logger = logger.With(zap.String("query", query))

		if query == "" {
			http.Error(w, "Missing query parameter, make sure to pass a ?q= GET param", http.StatusBadRequest)
			return
		}

		logger.Debug("executing query")
		start := time.Now()
		results, err := source.Query(ctx, query)
		if err != nil {
			logger.Error("failed execute query", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		took := time.Since(start)

		if results == nil {
			results = []sources.Result{}
		}

		resp := &response{
			Query: query,
			Took:  took,
			Data:  results,
			Count: len(results),
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			logger.Error("failed encode response", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		logger.Info("request served", zap.Duration("took", took), zap.Int("count", len(results)))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
	})

	mux := http.Handler(http.DefaultServeMux)

	if os.Getenv("ENV") == "production" {
		middleware := secure.New(secure.Options{
			AllowedHosts:         []string{os.Getenv("SSL_HOST")},
			HostsProxyHeaders:    []string{"X-Forwarded-Host"},
			SSLRedirect:          true,
			SSLHost:              os.Getenv("SSL_HOST"),
			SSLProxyHeaders:      map[string]string{"X-Forwarded-Proto": "https"},
			STSIncludeSubdomains: true,
			STSPreload:           true,
			FrameDeny:            true,
			ContentTypeNosniff:   true,
			BrowserXssFilter:     true,
		})
		mux = middleware.Handler(mux)
	}

	middleware := cors.New(cors.Options{AllowedOrigins: []string{"*"}})
	mux = middleware.Handler(mux)

	addr := ":8080"
	logger.Info("service listening", zap.String("addr", addr))
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		logger.Fatal("service shutdown", zap.Error(err))
	}
}
