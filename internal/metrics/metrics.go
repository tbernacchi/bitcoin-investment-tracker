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

	bitcoinPriceChangePercent = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bitcoin_price_change_percent",
		Help: "Current price change percentage",
	})
)

// UpdateMetrics atualiza todas as métricas
func UpdateMetrics(price float64, changePercent float64, initialInvestment float64, btcAmount float64) {
	bitcoinPrice.Set(price)
	bitcoinPriceChangePercent.Set(changePercent)

	// Calcular outras métricas
	currentValue := btcAmount * price
	profitPercent := ((currentValue - initialInvestment) / initialInvestment) * 100

	bitcoinAmount.Set(btcAmount)
	bitcoinInvestmentValue.WithLabelValues("initial").Set(initialInvestment)
	bitcoinInvestmentValue.WithLabelValues("current").Set(currentValue)
	bitcoinProfitPercent.Set(profitPercent)
}
