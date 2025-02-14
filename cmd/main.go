package main

import (
	"TradesAggregator/internal/store"
	"TradesAggregator/pkg/poloniex"
)

func main() {
	//wsConnect()
	candles, _ := fetchCandles(poloniex.RestAPI, poloniex.CandlesResource, "BTC_USDT", "MINUTE_1")
	// for _, k := range candles {
	// 	fmt.Println(k)
	// }
	store.Batch(candles)
}
