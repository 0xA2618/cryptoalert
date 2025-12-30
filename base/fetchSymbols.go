package base

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type SymbolInfo struct {
	Symbol         string
	Price          float64
	Volume         float64
	TakerBuyVolume float64
	TakerBuyRatio  float64
	Rsi            float64
	Rate           float64
	CrossType      string
	CrossTime      time.Time
	Shape          string
	Change         float64
}

var SymbolList []*SymbolInfo

type SymbolPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type SymbolChange struct {
	Symbol      string
	ClosePrice  float64
	CycleChange map[string]float64 // 周期 => 涨跌幅
}

type ChangeInfo struct {
	Change     float64
	ClosePrice float64
}

func FetchBinanceSymbols() {
	resp, err := http.Get("https://fapi.binance.com/fapi/v1/ticker/price")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var prices []SymbolPrice
	json.Unmarshal(body, &prices)

	for _, p := range prices {
		fmt.Println(p.Symbol, p.Price)
		price, _ := strconv.ParseFloat(p.Price, 64)
		SymbolList = append(SymbolList, &SymbolInfo{Symbol: p.Symbol, Price: price})
	}
}

// 获取指定 symbol 和周期的涨跌幅
func FetchChange(symbol, interval string) *ChangeInfo {
	url := fmt.Sprintf("https://fapi.binance.com/fapi/v1/klines?symbol=%s&interval=%s&limit=2", symbol, interval)
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var klines [][]interface{}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &klines); err != nil || len(klines) < 2 {
		return nil
	}

	open, _ := strconv.ParseFloat(klines[0][1].(string), 64)
	closeVal, _ := strconv.ParseFloat(klines[1][4].(string), 64)

	change := (closeVal - open) / open * 100
	return &ChangeInfo{Change: change, ClosePrice: closeVal}
}

func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
