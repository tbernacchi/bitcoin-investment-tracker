package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tbernacchi/bitcoin-investment-tracker/internal/metrics"
	"github.com/tbernacchi/bitcoin-investment-tracker/internal/websocket"
)

func main() {
	// URL do websocket da Binance
	binanceWsURL := "wss://stream.binance.com/ws/btcbrl@ticker"

	// Starting HTTP server
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()

	// Starting WebSocket connection to Binance
	go websocket.MonitorPrices(
		binanceWsURL,
		metrics.UpdateMetrics,
	)

	// Keep the program running
	select {}
}
