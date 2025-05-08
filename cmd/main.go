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
	// Carregar variáveis do .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Printf("Bitcoin price monitor started")

	// URL do websocket da Binance
	binanceWsURL := "wss://stream.binance.com/ws/btcbrl@ticker"

	// Iniciar servidor web
	go webserver.StartWebServer()

	// Configurar endpoint de métricas
	http.Handle("/metrics", promhttp.Handler())

	// Iniciar servidor de métricas na porta 2112
	go func() {
		log.Printf("Starting metrics server on :2112")
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()

	// Iniciar WebSocket para métricas
	go websocket.MonitorPrices(
		binanceWsURL,
		metrics.UpdateMetrics,
	)

	// Manter o programa rodando
	select {}
}
