package main

import (
	"log"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tbernacchi/bitcoin-investment-tracker/internal/calculator"
	"github.com/tbernacchi/bitcoin-investment-tracker/internal/webserver"
	"github.com/tbernacchi/bitcoin-investment-tracker/internal/websocket"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Printf("Bitcoin price monitor started")

	// Binance WebSocket URL for BTC/BRL pair
	// Using the ticker stream which provides 24hr rolling window price change statistics
	binanceWsURL := "wss://stream.binance.com/ws/btcbrl@ticker"

	// Channel to keep main goroutine alive
	done := make(chan bool)

	// Start web server in a separate goroutine
	go webserver.StartWebServer()

	// Start WebSocket connection to Binance in a separate goroutine
	go websocket.MonitorPrices(binanceWsURL, func(price float64) {
		calculator.ShowInvestmentInfo(price, 1.0)
	})

	// Endpoint to expose metrics
	http.Handle("/metrics", promhttp.Handler())

	// Start the web server http on port 2112
	log.Fatal(http.ListenAndServe(":2112", nil))

	// Keep the program running
	<-done
}
