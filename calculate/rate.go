package calculate

import (
	"crypto_alert/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
)

type FundingRate struct {
	Symbol      string  `json:"symbol"`
	FundingRate float64 `json:"fundingRate,string"`
	FundingTime int64   `json:"fundingTime"`
}

type PremiumIndexResponse struct {
	LastFundingRate string `json:"lastFundingRate"`
}

func getFundingRate(symbol string) float64 {
	url := fmt.Sprintf("https://fapi.binance.com/fapi/v1/premiumIndex?symbol=%s", symbol)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("failed to request: ", err)
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("bad status: %s, body: %s", resp.Status, string(body))
		return 0
	}

	var result PremiumIndexResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("failed to decode response: %v", err)
		return 0
	}

	// 转成 float64
	rateFloat, err := strconv.ParseFloat(result.LastFundingRate, 64)
	if err != nil {
		fmt.Printf("failed to parse funding rate: %v", err)
		return 0
	}

	return rateFloat * 100
}

func GetRate(symbol string) float64 {
	url := fmt.Sprintf(config.Cfg.Api.Binance.FApi.Rate, symbol)
	resp, err := http.Get(url)
	if err != nil {
		log.Error("get rate failed: %s\nurl: %s\nerror: %s\n", url, err.Error())
		return 0
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("get rate parser request body error:", err.Error())
		return 0
	}

	// 这个接口返回的是数组
	var result []FundingRate
	if err := json.Unmarshal(body, &result); err != nil {
		log.Error("rate Unmarshal failed: ", err.Error())
		return 0
	}

	if len(result) == 0 {
		log.Error("rate response is empty")
		return 0
	}

	// 返回最新的 FundingRate
	return result[0].FundingRate
}
