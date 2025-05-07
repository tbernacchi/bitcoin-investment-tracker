package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

// BinanceTickerMessage represents the structure of incoming WebSocket messages
type BinanceTickerMessage struct {
	EventType    string `json:"e"` // Event type
	EventTime    int64  `json:"E"` // Event time
	Symbol       string `json:"s"` // Trading pair symbol
	PriceChange  string `json:"p"` // Price change
	PricePercent string `json:"P"` // Price change percent
	WeightedAvg  string `json:"w"` // Weighted average price
	PrevClose    string `json:"x"` // Previous day's close price
	LastPrice    string `json:"c"` // Current price
	CloseQty     string `json:"Q"` // Last quantity
	BidPrice     string `json:"b"` // Best bid price
	BidQty       string `json:"B"` // Best bid quantity
	AskPrice     string `json:"a"` // Best ask price
	AskQty       string `json:"A"` // Best ask quantity
	OpenPrice    string `json:"o"` // Open price
	HighPrice    string `json:"h"` // High price
	LowPrice     string `json:"l"` // Low price
	Volume       string `json:"v"` // Total traded volume
	QuoteVolume  string `json:"q"` // Total traded quote volume
	OpenTime     int64  `json:"O"` // Statistics open time
	CloseTime    int64  `json:"C"` // Statistics close time
	FirstTradeID int64  `json:"F"` // First trade ID
	LastTradeID  int64  `json:"L"` // Last trade ID
	TradeCount   int64  `json:"n"` // Total number of trades
}

// formatPrice formats the price in Brazilian currency format
func formatPrice(price float64) string {
	// Convert to string with 2 decimal places
	priceStr := fmt.Sprintf("%.2f", price)

	// Split into integer and decimal parts
	parts := strings.Split(priceStr, ".")

	// Format integer part with thousand separators
	intPart := parts[0]
	var formatted []string
	for i := len(intPart); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		formatted = append([]string{intPart[start:i]}, formatted...)
	}

	// Join with dots and replace decimal point with comma
	return fmt.Sprintf("R$ %s,%s", strings.Join(formatted, "."), parts[1])
}

// MonitorPrices establishes a WebSocket connection to Binance
func MonitorPrices(wsURL string, priceCallback func(float64)) {
	log.Printf("Connecting to Binance WebSocket...")

	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("WebSocket connection failed:", err)
	}
	defer c.Close()

	log.Printf("Successfully connected to Binance WebSocket")

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			return
		}

		var ticker BinanceTickerMessage
		if err := json.Unmarshal(message, &ticker); err != nil {
			log.Println("JSON parsing error:", err)
			log.Printf("Received message: %s", string(message))
			continue
		}

		// Convert string price to float64
		price, err := strconv.ParseFloat(ticker.LastPrice, 64)
		if err != nil {
			log.Println("Price conversion error:", err)
			continue
		}

		// Get 24h price change percentage
		changePercent, err := strconv.ParseFloat(ticker.PricePercent, 64)
		if err != nil {
			log.Println("Change percent conversion error:", err)
			continue
		}

		// Format the log message with price and change percentage
		var changeSymbol string
		if changePercent > 0 {
			changeSymbol = "↑" // up arrow for positive change
		} else if changePercent < 0 {
			changeSymbol = "↓" // down arrow for negative change
		} else {
			changeSymbol = "=" // equals for no change
		}

		log.Printf("BTC Price: %s (%s%.2f%%)", formatPrice(price), changeSymbol, changePercent)
		priceCallback(price)
	}
}
