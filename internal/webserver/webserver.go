package webserver

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tbernacchi/bitcoin-investment-tracker/internal/formatter"
)

func StartWebServer() {
	// Create a handler that logs access
	loggedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handleHome(w, r)

		// Format the time in milliseconds
		duration := time.Since(start).Milliseconds()

		// Get the real IP considering proxy headers
		ip := r.Header.Get("X-Real-IP")
		if ip == "" {
			ip = r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = r.RemoteAddr
			}
		}

		log.Printf(
			"[%s] %s from %s - %dms",
			r.Method,
			r.RequestURI,
			ip,
			duration,
		)
	})

	http.Handle("/", loggedHandler)

	log.Printf("Starting web server on http://localhost:8080")
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("Web server error: %v", err)
		}
	}()
}

type PageData struct {
	Investment         string
	BTCAmount          float64
	CurrentBTCPrice    string
	CurrentBTCPriceUSD float64
	CurrentValue       string
	Profit             string
	ProfitPercent      float64
	LastUpdate         string
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Bitcoin Tracker</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f8f9fa;
            color: #212529;
        }
        .price-card {
            background: white;
            border-radius: 12px;
            padding: 24px;
            margin: 20px 0;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
            text-align: center;
        }
        .price {
            font-size: 48px;
            font-weight: bold;
            margin: 10px 0;
        }
        .change {
            font-size: 20px;
            margin: 10px 0;
        }
        .positive { color: #198754; }
        .negative { color: #dc3545; }
        .neutral { color: #6c757d; }
        .small-text {
            font-size: 14px;
            color: #6c757d;
        }
        .header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 20px;
        }
        .last-update {
            font-size: 12px;
            color: #6c757d;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>Bitcoin Tracker</h1>
        <span class="last-update" id="lastUpdate">{{.LastUpdate}}</span>
    </div>

    <div class="price-card">
        <div class="small-text">Bitcoin Price</div>
        <div class="price" id="btcPrice">{{.CurrentBTCPrice}}</div>
        <div class="change" id="priceChange">
            <span id="changeValue"></span>
        </div>
    </div>

    <div class="price-card">
        <div class="small-text">Your Investment</div>
        <div>Amount: <strong>{{.BTCAmount}} BTC</strong></div>
        <div>Value: <strong id="currentValue">{{.CurrentValue}}</strong></div>
        <div class="change" id="profitLoss">
            <span id="profitValue">{{.Profit}} ({{printf "%.2f" .ProfitPercent}}%)</span>
        </div>
    </div>

    <script>
    const ws = new WebSocket('ws://' + window.location.host + '/ws');
    
    ws.onmessage = function(event) {
        const data = JSON.parse(event.data);
        
        // Update price
        document.getElementById('btcPrice').textContent = data.price;
        
        // Update change
        const changeElement = document.getElementById('changeValue');
        const changeClass = data.changePercent > 0 ? 'positive' : (data.changePercent < 0 ? 'negative' : 'neutral');
        const changeSymbol = data.changePercent > 0 ? '↑' : (data.changePercent < 0 ? '↓' : '=');
        changeElement.textContent = changeSymbol + ' ' + data.changePercent.toFixed(2) + '%';
        changeElement.className = changeClass;
        
        // Update investment value
        document.getElementById('currentValue').textContent = data.currentValue;
        
        // Update profit/loss
        const profitElement = document.getElementById('profitValue');
        profitElement.textContent = data.profit + ' (' + data.profitPercent.toFixed(2) + '%)';
        profitElement.className = data.profitPercent > 0 ? 'positive' : 'negative';
        
        // Update timestamp
        document.getElementById('lastUpdate').textContent = new Date().toLocaleTimeString();
    };
    
    ws.onclose = function() {
        console.log('WebSocket connection closed. Reconnecting...');
        setTimeout(function() {
            window.location.reload();
        }, 5000);
    };
    </script>
</body>
</html>
`

	// Get updated data
	btcPrice, err := GetBTCPrice(os.Getenv("MERCADO_BITCOIN_API_URL"))
	if err != nil {
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}

	// Get USD/BRL rate from Binance
	usdBrlRate, err := getUSDRate()
	if err != nil {
		http.Error(w, "Error fetching USD rate", http.StatusInternalServerError)
		return
	}

	investment, _ := strconv.ParseFloat(os.Getenv("INVESTMENT_BRL"), 64)
	btcAmount, _ := strconv.ParseFloat(os.Getenv("BTC_AMOUNT"), 64)

	currentValue := btcAmount * btcPrice
	profit := currentValue - investment
	profitPercent := (profit / investment) * 100
	btcPriceUSD := btcPrice / usdBrlRate

	data := PageData{
		Investment:         formatter.FormatBRL(investment),
		BTCAmount:          btcAmount,
		CurrentBTCPrice:    formatter.FormatBRL(btcPrice),
		CurrentBTCPriceUSD: btcPriceUSD,
		CurrentValue:       formatter.FormatBRL(currentValue),
		Profit:             formatter.FormatBRL(profit),
		ProfitPercent:      profitPercent,
		LastUpdate:         time.Now().Format("02/01/2006 15:04:05"),
	}

	funcMap := template.FuncMap{
		"formatUSD": formatter.FormatUSD,
	}

	t := template.Must(template.New("home").Funcs(funcMap).Parse(tmpl))
	t.Execute(w, data)
}

func getUSDRate() (float64, error) {
	binanceUrl := os.Getenv("BINANCE_USDT_API_URL")
	if binanceUrl == "" {
		return 0, fmt.Errorf("Missing BINANCE_USDT_API_URL in .env")
	}

	resp, err := http.Get(binanceUrl)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Price string `json:"price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return strconv.ParseFloat(result.Price, 64)
}

// Function to get the BTC price
func GetBTCPrice(apiUrl string) (float64, error) {
	resp, err := http.Get(apiUrl)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var tickerResp struct {
		Ticker struct {
			Last string `json:"last"`
		} `json:"ticker"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tickerResp); err != nil {
		return 0, err
	}

	return strconv.ParseFloat(tickerResp.Ticker.Last, 64)
}
