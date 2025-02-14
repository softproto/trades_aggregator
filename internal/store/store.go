package store

import (
	"TradesAggregator/pkg/poloniex"
	"log"
)

func Single(c poloniex.Kline) {
	log.Println(c)
}

func Batch(c []poloniex.Kline) {

	for _, b := range c {
		log.Println(b)
	}
}
