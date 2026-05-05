package httpapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"config-audit/internal/app"
)

func Run(ctx context.Context, addr string, analyzer app.Analyzer) error {
	server := &http.Server{
		Addr:              addr,
		Handler:           NewHandler(analyzer),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("run http server: %w", err)
	}

	return nil
}
