package calculate

import (
	"crypto_alert/base"
	"crypto_alert/config"
	"fmt"
)

// 计算EMA
func calculateEMA(prices []float64, period int) []float64 {
	ema := make([]float64, len(prices))
	k := 2.0 / float64(period+1)
	for i, price := range prices {
		if i == 0 {
			ema[i] = price
		} else {
			ema[i] = price*k + ema[i-1]*(1-k)
		}
	}
	return ema
}

// 计算MACD
func calculateMACD(closes []float64) ([]float64, []float64, []float64) {
	emaFast := calculateEMA(closes, config.Cfg.Benchmark.Macd.FastPeriod)
	emaSlow := calculateEMA(closes, config.Cfg.Benchmark.Macd.SlowPeriod)
	macd := make([]float64, len(closes))
	for i := range closes {
		macd[i] = emaFast[i] - emaSlow[i]
	}
	signalLine := calculateEMA(macd, config.Cfg.Benchmark.Macd.Window)
	histogram := make([]float64, len(closes))
	for i := range closes {
		histogram[i] = macd[i] - signalLine[i]
	}
	return macd, signalLine, histogram
}

// 检测柱状图变化
func detectReversal(klines []base.KLine, histogram []float64) string {
	if len(histogram) < 3 {
		fmt.Println("数据不足，无法判断")
		return ""
	}

	// 取最后三根k线
	last3kline := klines[len(klines)-3:]

	// 初始化 maxHigh 和 minLow
	k3MaxPrice := last3kline[0].High
	k3MinPrice := last3kline[0].Low

	for _, k := range last3kline {
		if k.High > k3MaxPrice {
			k3MaxPrice = k.High
		}
		if k.Low < k3MinPrice {
			k3MinPrice = k.Low
		}
	}

	// 取最后三根柱子
	last3 := histogram[len(histogram)-3:]

	// fmt.Printf("最近三根柱子: %.5f, %.5f, %.5f\n", last3[0], last3[1], last3[2])

	// 判断动能增强 or 减弱
	delta1 := last3[1] - last3[0]
	delta2 := last3[2] - last3[1]

	if last3[2] > 0 {
		// 多头
		if delta1 > 0 && delta2 < 0 {
			return "近3根K线多头动能减弱"
		}
	} else {
		// 空头
		if delta1 < 0 && delta2 > 0 {
			return "近3根K线空头动能减弱"
		}
	}
	return ""
}

// 检测金叉和死叉
func detectCrosses(klines []base.KLine, macd, signalLine []float64) (string, int) {
	var gold, die string
	var index int

	for i := 1; i < len(klines); i++ {
		prevMacd := macd[i-1]
		prevSignal := signalLine[i-1]
		currMacd := macd[i]
		currSignal := signalLine[i]

		if klines[i].Volume > 0 {
			// 金叉逻辑
			if prevMacd <= prevSignal && currMacd > currSignal {
				gold = "金叉"
				die = ""
				index = i

			}

			// 死叉逻辑
			if prevMacd >= prevSignal && currMacd < currSignal {
				die = "死叉"
				gold = ""
				index = i
			}
		}
	}

	if die == "" {
		return gold, index
	}

	if gold == "" {
		return die, index
	}
	return "", 0

}
