package calculate

import (
	"crypto_alert/base"
	"crypto_alert/config"

	"context"
	"crypto_alert/utils/notify"
	"fmt"
	"time"
)

// 进行macd
func Start(ctx context.Context, cycle string) {
	for _, symbolInfo := range base.SymbolList {
		symbol := symbolInfo.Symbol
		Msg := ""

		// 获取K线数据
		klines, err := base.GetContractKlines(symbol, cycle)
		if err != nil {
			fmt.Println("错误:", err)
			continue
		}
		if len(klines) < config.Cfg.Benchmark.Klines {
			fmt.Printf("数据不足: %s 有 %d 条, 需要 %d 条\n", symbol, len(klines), config.Cfg.Benchmark.Klines)
			continue
		}
		lastKline := klines[config.Cfg.Benchmark.Klines-1:]
		if (time.Now().Unix() - lastKline[0].OpenTime) > 60*10 {
			lastKlineTime := time.Unix(lastKline[0].OpenTime/1000, 0).Format("2006-01-02 15:04:05")
			fmt.Println("最后一根k线距当前时间超过10分钟，无效数据:", lastKlineTime)
			continue
		}

		// 计算涨跌幅
		latestKline := klines[len(klines)-1]
		symbolInfo.Change = (latestKline.Close - latestKline.Open) / latestKline.Open * 100

		// 处理收线价格
		closes := base.ClosePrice(klines)

		// 计算MACD (快线7，慢线25，信号线9)
		macd, signalLine, histogram := calculateMACD(closes)

		// 计算交叉
		crossType, klineIndex := detectCrosses(klines, macd, signalLine)
		rsiValue := GetRsi(closes)

		if klineIndex != 0 && checkRsiValue(rsiValue) {
			crossTime := time.Unix(klines[klineIndex].OpenTime/1000, 0)
			takerBuyRatio := (klines[klineIndex].TakerBuyVolume / klines[klineIndex].Volume) * 100
			symbolInfo.Rsi = rsiValue
			symbolInfo.Rate = getFundingRate(symbolInfo.Symbol)
			symbolInfo.Price = klines[klineIndex].Close
			symbolInfo.Volume = klines[klineIndex].Volume
			symbolInfo.TakerBuyVolume = klines[klineIndex].TakerBuyVolume
			symbolInfo.TakerBuyRatio = takerBuyRatio
			symbolInfo.CrossType = crossType

			if symbolInfo.CrossTime != crossTime {
				symbolInfo.CrossTime = crossTime
				Msg = happenCrossFmt(symbolInfo, cycle)
			}

		}

		// 3根K线变化
		reversal := detectReversal(klines, histogram)
		if reversal != "" && Msg == "" && config.Cfg.Benchmark.Detect == "true" {
			symbolInfo.Shape = reversal
			symbolInfo.CrossTime = time.Now()
			symbolInfo.Rsi = rsiValue
			Msg = trendFmt(symbolInfo, cycle)
		}

		if Msg == "" {
			continue
		}

		notify.SendTelegramMessage(cycle, Msg)
	}
}

// ticker
func MacdTicker(ctx context.Context, cycle string) {
	duration := CycleDurationFmt(cycle)

	// 计算距离下一次整点的时间
	now := time.Now()
	nextTick := now.Truncate(duration).Add(duration)
	waitTime := nextTick.Sub(now)

	fmt.Printf("[%s] 任务将在 %s (等待 %v) 后开始\n", cycle, nextTick.Format("15:04:05"), waitTime)

	// 等待第一次执行
	timer := time.NewTimer(waitTime)
	select {
	case <-ctx.Done():
		timer.Stop()
		fmt.Printf("周期任务 %s 收到退出信号\n", cycle)
		return
	case <-timer.C:
		fmt.Printf("周期性任务: %s 首次执行中...\n", cycle)
		go Start(ctx, cycle)
	}

	// 启动周期性 Ticker
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("周期任务 %s 收到退出信号\n", cycle)
			return

		case <-ticker.C:
			fmt.Printf("周期性任务: %s 执行中...\n", cycle)
			go Start(ctx, cycle)
		}
	}

}
