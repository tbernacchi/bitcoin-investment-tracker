package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Métricas de investimento
	bitcoinInvestmentValue = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bitcoin_investment_value",
		Help: "Investment values in BRL",
	}, []string{"type"}) // type pode ser "initial", "current"

	bitcoinProfitPercent = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bitcoin_investment_profit_percent",
		Help: "Current profit/loss percentage",
	})

	bitcoinAmount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bitcoin_amount_held",
		Help: "Amount of Bitcoin held",
	})

	bitcoinPrice = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bitcoin_current_price_brl",
		Help: "Current Bitcoin price in BRL",
	})
)

// UpdateMetrics updates all investment metrics
func UpdateMetrics(initialInvestment, currentValue, btcAmount, btcPrice float64) {
	// Investment values
	bitcoinInvestmentValue.WithLabelValues("initial").Set(initialInvestment)
	bitcoinInvestmentValue.WithLabelValues("current").Set(currentValue)

	// Calculate and update the profit/loss percentage
	profitPercent := ((currentValue - initialInvestment) / initialInvestment) * 100
	bitcoinProfitPercent.Set(profitPercent)

	// Bitcoin amount
	bitcoinAmount.Set(btcAmount)

	// Current Bitcoin price
	bitcoinPrice.Set(btcPrice)
}
