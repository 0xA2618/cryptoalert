package calculate

import (
	"crypto_alert/base"
	"fmt"
	"time"
)

// 格式化为普通或“万”单位
func formatWithWan(val float64) string {
	if val >= 10000 || val <= 10000 {
		return fmt.Sprintf("%.2f万", val/10000)
	}
	return fmt.Sprintf("%.2f", val)
}

// 发生交叉 消息 格式化
func happenCrossFmt(info *base.SymbolInfo, cycle string) string {

	return fmt.Sprintf(
		"     ---- 【  %s %s  】 ---- \n"+
			"价格: %.4f\n"+
			"涨跌: %.2f%%\n"+
			"MACD: %s\n"+
			"RSI: %.2f\n"+
			"费率: %.4f%%\n"+
			"时间: %s",
		info.Symbol,
		cycle,
		info.Price,
		info.Change,
		info.CrossType,
		info.Rsi,
		info.Rate,
		info.CrossTime.Format("2006-01-02 15:04:05"),
	)
}

// 趋势变化 消息 格式化
func trendFmt(info *base.SymbolInfo, cycle string) string {
	return fmt.Sprintf("    --- 【 %s %s 趋势】 --- \n"+
		"价格: %.4f\n"+
		"涨跌: %.2f%%\n"+
		"上次交叉: %s\n"+
		"趋势: %s\n"+
		"RSI: %.2f\n"+
		"时间:  %s",
		info.Symbol,
		cycle,
		info.Price,
		info.Change,
		info.CrossType,
		info.Shape,
		info.Rsi,
		info.CrossTime.Format("2006-01-02 15:04:05"),
	)

}

// 判断使用小时周期还是分钟周期
func CycleDurationFmt(cycle string) time.Duration {

	var duration time.Duration
	switch cycle {
	case "5m":
		duration = 5 * time.Minute
	case "30m":
		duration = 30 * time.Minute
	case "1h":
		duration = 1 * time.Hour
	case "4h":
		duration = 4 * time.Hour
	default:
		duration = 30 * time.Minute
	}

	// defaultTime := 30 * time.Minute
	// if len(cycle) == 0 {
	// 	fmt.Println("CheckCycleDuration is empty use default time 30m")
	// 	return defaultTime
	// }

	// if string(cycle[len(cycle)-1]) == "h" {
	// 	num, _ := strconv.Atoi(cycle[1:])
	// 	return time.Duration(num) * time.Hour
	// } else if string(cycle[0]) == "m" {
	// 	num, _ := strconv.Atoi(cycle[1:])
	// 	return time.Duration(num) * time.Minute
	// }

	return duration

}
