// calculator.go
package calculator

import (
	"os"
	"strconv"

	"github.com/tbernacchi/bitcoin-investment-tracker/internal/metrics"
)

type Investment struct {
	AmountBRL     float64
	CurrentBTCUSD float64
	CurrentUSDBRL float64
}

func (i *Investment) Calculate() (btcAmount, valueBRL float64) {
	// Convert BRL to USD
	amountUSD := i.AmountBRL / i.CurrentUSDBRL

	// Calculate BTC amount
	btcAmount = amountUSD / i.CurrentBTCUSD

	// Calculate current value in BRL
	valueBRL = btcAmount * i.CurrentBTCUSD * i.CurrentUSDBRL

	return btcAmount, valueBRL
}

func ShowInvestmentInfo(btcPriceBRL float64, usdBrlRate float64, changePercent float64) {
	investmentStr := os.Getenv("INVESTMENT_BRL")
	btcAmountStr := os.Getenv("BTC_AMOUNT")
	if investmentStr == "" || btcAmountStr == "" {
		return
	}

	investment, err := strconv.ParseFloat(investmentStr, 64)
	if err != nil {
		return
	}

	btcAmount, err := strconv.ParseFloat(btcAmountStr, 64)
	if err != nil {
		return
	}

	// Calculate values
	currentValueBRL := btcAmount * btcPriceBRL
	profit := currentValueBRL - investment
	_ = (profit / investment) * 100

	// Update metrics with the calculated values
	metrics.UpdateMetrics(
		btcPriceBRL,   // price
		changePercent, // changePercent
		investment,    // initialInvestment
		btcAmount,     // btcAmount
	)
}
