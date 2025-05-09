package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tbernacchi/bitcoin-investment-tracker/internal/metrics"
	"github.com/tbernacchi/bitcoin-investment-tracker/internal/webserver"
	"github.com/tbernacchi/bitcoin-investment-tracker/internal/websocket"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Printf("Bitcoin price monitor started")

	// Binance websocket URL
	binanceWsURL := "wss://stream.binance.com/ws/btcbrl@ticker"

	// Start web server
	go webserver.StartWebServer()

	// Configure metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Start metrics server on port 2112
	go func() {
		log.Printf("Starting metrics server on :2112")
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()

	// Start WebSocket for metrics
	go websocket.MonitorPrices(
		binanceWsURL,
		metrics.UpdateMetrics,
	)

	// Keep the program running
	select {}
}
