package calculate

import (
	"crypto_alert/config"
)

// 检查是否满足发送的条件
func checkRsiValue(rsi float64) bool {
	if config.Cfg.Benchmark.Rsi.Enable {
		rsiLow := float64(config.Cfg.Benchmark.Rsi.Low)
		rsiTop := float64(config.Cfg.Benchmark.Rsi.Top)
		if rsi < rsiLow || rsi > rsiTop {
			return true
		}
	}
	return false
}

func GetRsi(prices []float64) float64 {
	var gains, losses float64

	period := config.Cfg.Benchmark.Rsi.Period

	for i := 1; i <= period; i++ {
		diff := prices[i] - prices[i-1]
		if diff >= 0 {
			gains += diff
		} else {
			losses -= diff
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	for i := period + 1; i < len(prices); i++ {
		diff := prices[i] - prices[i-1]
		if diff >= 0 {
			avgGain = (avgGain*float64(period-1) + diff) / float64(period)
			avgLoss = (avgLoss * float64(period-1)) / float64(period)
		} else {
			avgGain = (avgGain * float64(period-1)) / float64(period)
			avgLoss = (avgLoss*float64(period-1) - diff) / float64(period)
		}
	}

	if avgLoss == 0 {
		return 100
	}

	rs := avgGain / avgLoss
	return 100 - (100 / (1 + rs))
}
