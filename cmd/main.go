package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tbernacchi/bitcoin-investment-tracker/internal/metrics"
	"github.com/tbernacchi/bitcoin-investment-tracker/internal/webserver"
	"github.com/tbernacchi/bitcoin-investment-tracker/internal/websocket"
)

// updateMetricsWrapper wraps the metrics.UpdateMetrics function to match the expected interface
func updateMetricsWrapper(price float64, changePercent float64) {
	// Get values from environment variables
	initialInvestment, err := strconv.ParseFloat(os.Getenv("INVESTMENT_BRL"), 64)
	if err != nil {
		log.Printf("Error parsing INVESTMENT_BRL: %v, using default 20000", err)
		initialInvestment = 20000
	}

	btcAmount, err := strconv.ParseFloat(os.Getenv("BTC_AMOUNT"), 64)
	if err != nil {
		log.Printf("Error parsing BTC_AMOUNT: %v, using default 0.1", err)
		btcAmount = 0.1
	}

	metrics.UpdateMetrics(price, changePercent, initialInvestment, btcAmount)
}

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
		updateMetricsWrapper,
	)

	// Keep the program running
	select {}
}
