package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Investment metrics
	bitcoinInvestmentValue = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bitcoin_investment_value",
		Help: "Investment values in BRL",
	}, []string{"type", "bitcoin_amount_held", "bitcoin_current_price_brl", "bitcoin_investment_value", "bitcoin_investment_profit_percent", "bitcoin_price_change_percent"})

	bitcoinProfitPercent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bitcoin_investment_profit_percent",
		Help: "Current profit/loss percentage",
	}, []string{"bitcoin_amount_held", "bitcoin_current_price_brl", "bitcoin_investment_value", "bitcoin_investment_profit_percent", "bitcoin_price_change_percent"})

	bitcoinAmount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bitcoin_amount_held",
		Help: "Amount of Bitcoin held",
	}, []string{"bitcoin_amount_held", "bitcoin_current_price_brl", "bitcoin_investment_value", "bitcoin_investment_profit_percent", "bitcoin_price_change_percent"})

	bitcoinPrice = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bitcoin_current_price_brl",
		Help: "Current Bitcoin price in BRL",
	}, []string{"bitcoin_amount_held", "bitcoin_current_price_brl", "bitcoin_investment_value", "bitcoin_investment_profit_percent", "bitcoin_price_change_percent"})

	bitcoinPriceChangePercent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bitcoin_price_change_percent",
		Help: "Current price change percentage",
	}, []string{"bitcoin_amount_held", "bitcoin_current_price_brl", "bitcoin_investment_value", "bitcoin_investment_profit_percent", "bitcoin_price_change_percent"})
)

// UpdateMetrics updates all metrics
func UpdateMetrics(price float64, changePercent float64, initialInvestment float64, btcAmount float64) {
	// Calculate values
	currentValue := btcAmount * price
	profitPercent := ((currentValue - initialInvestment) / initialInvestment) * 100

	// Create labels map
	labels := prometheus.Labels{
		"bitcoin_amount_held":               fmt.Sprintf("%.5f", btcAmount),
		"bitcoin_current_price_brl":         fmt.Sprintf("%.2f", price),
		"bitcoin_investment_value":          fmt.Sprintf("%.2f", currentValue),
		"bitcoin_investment_profit_percent": fmt.Sprintf("%.2f", profitPercent),
		"bitcoin_price_change_percent":      fmt.Sprintf("%.2f", changePercent),
	}

	// Update metrics with labels
	bitcoinPrice.With(labels).Set(price)
	bitcoinPriceChangePercent.With(labels).Set(changePercent)
	bitcoinAmount.With(labels).Set(btcAmount)
	bitcoinInvestmentValue.With(prometheus.Labels{
		"type":                              "initial",
		"bitcoin_amount_held":               fmt.Sprintf("%.5f", btcAmount),
		"bitcoin_current_price_brl":         fmt.Sprintf("%.2f", price),
		"bitcoin_investment_value":          fmt.Sprintf("%.2f", currentValue),
		"bitcoin_investment_profit_percent": fmt.Sprintf("%.2f", profitPercent),
		"bitcoin_price_change_percent":      fmt.Sprintf("%.2f", changePercent),
	}).Set(initialInvestment)
	bitcoinInvestmentValue.With(prometheus.Labels{
		"type":                              "current",
		"bitcoin_amount_held":               fmt.Sprintf("%.5f", btcAmount),
		"bitcoin_current_price_brl":         fmt.Sprintf("%.2f", price),
		"bitcoin_investment_value":          fmt.Sprintf("%.2f", currentValue),
		"bitcoin_investment_profit_percent": fmt.Sprintf("%.2f", profitPercent),
		"bitcoin_price_change_percent":      fmt.Sprintf("%.2f", changePercent),
	}).Set(currentValue)
	bitcoinProfitPercent.With(labels).Set(profitPercent)
}
