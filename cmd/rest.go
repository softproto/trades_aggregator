package main

import (
	"TradesAggregator/pkg/poloniex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func fetchKlines(endpoint, resource, symbol, interval string, startTime, endTime int64) ([]poloniex.Kline, error) {
	kline := []poloniex.Kline{}
	url := fmt.Sprintf("%s%s%s?interval=%s&startTime=%d&endTime=%d", endpoint, symbol, resource, interval, startTime, endTime)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//some magic to prepare the string
	b := strings.ReplaceAll(string(body), "\"", "")
	parts := strings.Split(b[4:len(b)-4], " ], [ ")

	//cant avoid this
	for _, part := range parts {
		c := strings.Split(part, ", ")
		k := poloniex.Kline{
			Pair:      symbol,
			TimeFrame: interval,
			O:         toFloat64(c[2]),
			H:         toFloat64(c[1]),
			L:         toFloat64(c[0]),
			C:         toFloat64(c[3]),
			UtcBegin:  toInt64(c[12]),
			UtcEnd:    toInt64(c[13]),
			VolumeBS: poloniex.VBS{
				BuyBase:   toFloat64(c[6]),
				SellBase:  toFloat64(c[4]) - toFloat64(c[6]),
				BuyQuote:  toFloat64(c[7]),
				SellQuote: toFloat64(c[5]) - toFloat64(c[7]),
			},
		}
		kline = append(kline, k)

	}

	return kline, nil
}

func toInt64(s string) int64 {
	r, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Println("ParseInt() error: ", err)
		return 0
	}
	return r
}

func toFloat64(s string) float64 {
	r, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Println("ParseFloat() error: ", err)
		return 0
	}
	return r
}
