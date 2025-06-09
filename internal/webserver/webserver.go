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
    <title>Bitcoin Investment Tracker</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .card {
            background: white;
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .value {
            font-size: 24px;
            font-weight: bold;
            color: #333;
        }
        .profit-positive {
            color: #28a745;
        }
        .profit-negative {
            color: #dc3545;
        }
        .last-update {
            color: #666;
            font-size: 12px;
            text-align: right;
        }
    </style>
    <script>
        function refreshPage() {
            location.reload();
        }
        setInterval(refreshPage, 10000); // Atualiza a cada 10 segundos
    </script>
</head>
<body>
    <h1>Bitcoin Investment Tracker</h1>
    
    <div class="card">
        <h2>Investment Summary</h2>
        <p>Initial Investment: <span class="value">{{.Investment}}</span></p>
        <p>Bitcoin Amount: <span class="value">{{.BTCAmount}}</span> BTC</p>
        <p>Current BTC Price: <span class="value">{{.CurrentBTCPrice}} ({{formatUSD .CurrentBTCPriceUSD}})</span></p>
        <p>Current Value: <span class="value">{{.CurrentValue}}</span></p>
        <p>Profit/Loss: <span class="value {{if gt .ProfitPercent 0.0}}profit-positive{{else}}profit-negative{{end}}">
            {{.Profit}} ({{printf "%.2f" .ProfitPercent}}%)
        </span></p>
    </div>

    <p class="last-update">Last update: {{.LastUpdate}}</p>
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
