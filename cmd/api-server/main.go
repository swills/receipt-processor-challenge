package main

import (
	"log/slog"
	"net/http"
	"time"

	receiptserver "github.com/swills/receipt-processor-challenge/internal/receipt/server"
)

func main() {
	server := receiptserver.NewServer()

	mux := http.NewServeMux()

	handler := receiptserver.HandlerFromMux(server, mux)

	httpServer := &http.Server{
		Handler:     handler,
		Addr:        "0.0.0.0:8888",
		ReadTimeout: 5 * time.Second,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		slog.Error("error listening", "err", err)
	}
}
