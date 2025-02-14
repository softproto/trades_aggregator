package main

import (
	"TradesAggregator/internal/store"
	"TradesAggregator/pkg/poloniex"
	"log"
	"sync"
	"time"
)

func main() {
	log.Println("main() started")

	log.Println("fetch recent Klines")
	startTime := time.Date(2024, time.December, 1, 0, 0, 0, 0, time.UTC).Unix()
	for _, symbol := range poloniex.Symbols {
		for _, interval := range poloniex.Intervals {
			Klines, _ := fetchKlines(poloniex.RestAPI, poloniex.CandlesResource, symbol, interval, startTime, 0)
			store.Batch(Klines)
		}
	}
	var wg sync.WaitGroup

	log.Println("collect realtime data")
	wg.Add(1)
	go wsConnect(&wg)

	wg.Wait()
	log.Println("main() exited")

}
