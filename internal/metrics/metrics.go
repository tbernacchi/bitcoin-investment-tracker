package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Métricas de investimento
	bitcoinInvestmentValue = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bitcoin_investment_value_brl",
		Help: "Current investment value in BRL",
	})

	bitcoinInitialInvestment = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bitcoin_initial_investment_brl",
		Help: "Initial investment value in BRL",
	})

	bitcoinProfitPercent = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bitcoin_profit_percent",
		Help: "Current profit/loss percentage",
	})

	bitcoinAmount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bitcoin_amount_btc",
		Help: "Amount of Bitcoin held",
	})

	bitcoinPrice = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bitcoin_price_brl",
		Help: "Current Bitcoin price in BRL",
	})

	bitcoinPriceChangePercent = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bitcoin_price_change_percent_24h",
		Help: "24h price change percentage",
	})
)

// UpdateMetrics atualiza todas as métricas
func UpdateMetrics(price float64, changePercent float64, initialInvestment float64, btcAmount float64) {
	bitcoinPrice.Set(price)
	bitcoinPriceChangePercent.Set(changePercent)
	bitcoinAmount.Set(btcAmount)
	bitcoinInitialInvestment.Set(initialInvestment)

	// Calcular valor atual e lucro
	currentValue := btcAmount * price
	profitPercent := ((currentValue - initialInvestment) / initialInvestment) * 100

	bitcoinInvestmentValue.Set(currentValue)
	bitcoinProfitPercent.Set(profitPercent)
}
