package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"TradesAggregator/pkg/poloniex"

	"TradesAggregator/internal/store"
)

func wsConnect(wg *sync.WaitGroup) {
	defer wg.Done()
	conn, _, err := websocket.DefaultDialer.Dial(poloniex.WssAPI, nil)
	if err != nil {
		log.Fatal("Dial() error:", err)
	}
	defer conn.Close()

	//make quoted string from slice
	var builder strings.Builder
	for i, s := range poloniex.Symbols {
		builder.WriteString(`"` + s + `"`)
		if i < len(poloniex.Symbols)-1 {
			builder.WriteString(", ")
		}
	}
	symbols := builder.String()
	message := fmt.Sprintf(`{"event":"subscribe","channel":["%s"], "symbols":[%s]}`, poloniex.TradesChannel, symbols)
	// message := `{"event":"subscribe","channel":["trades"],"symbols":["BTC_USDT","TRX_USDT","ETH_USDT","DOGE_USDT","BCH_USDT"]}`
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Fatal("WriteMessage() error:", err)
	}

	message = `{"event":"ping"}`
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Fatal("WriteMessage() error:", err)
	}

	// ping-message timer (needs 30 sec, but...)
	pinger := time.NewTicker(20 * time.Second)
	defer pinger.Stop()

	// candle aggregator timers
	candlerMinute1 := time.NewTicker(1 * time.Minute)
	defer candlerMinute1.Stop()

	candlerMinute15 := time.NewTicker(15 * time.Minute)
	defer candlerMinute15.Stop()

	candlerHour1 := time.NewTicker(1 * time.Hour)
	defer candlerHour1.Stop()

	candlerDay1 := time.NewTicker(24 * time.Hour)
	defer candlerDay1.Stop()

	// gracefull interruption
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// just a trick to create the first candles for each pair
	// in real life we ​​don't create candles at runtime
	// instead we save trades (tick data) to the database
	// and generate candles on demand
	candles_minute_1 := make(map[string]poloniex.Kline)
	for _, s := range poloniex.Symbols {
		candles_minute_1[s] = poloniex.Kline{
			L:        999999,
			UtcBegin: time.Now().Unix(),
		}
	}
	candles_minute_15 := make(map[string]poloniex.Kline)
	for _, s := range poloniex.Symbols {
		candles_minute_15[s] = poloniex.Kline{
			L:        999999,
			UtcBegin: time.Now().Unix(),
		}
	}

	candles_hour_1 := make(map[string]poloniex.Kline)
	for _, s := range poloniex.Symbols {
		candles_hour_1[s] = poloniex.Kline{
			L:        999999,
			UtcBegin: time.Now().Unix(),
		}
	}

	candles_day_1 := make(map[string]poloniex.Kline)
	for _, s := range poloniex.Symbols {
		candles_day_1[s] = poloniex.Kline{
			L:        999999,
			UtcBegin: time.Now().Unix(),
		}
	}

	for {
		select {
		//send ping message
		case <-pinger.C:
			message := `{"event":"ping"}`
			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("WriteMessage() error:", err)
				return
			}

			//make candle candles_minute_1
		case <-candlerMinute1.C:
			log.Println("candles_minute_1")

			//make 1-minute candles
			for _, s := range poloniex.Symbols {
				c := candles_minute_1[s]
				c.TimeFrame = "candles_minute_1"
				c.UtcEnd = time.Now().Unix()
				//store candle
				store.Single(c)
				//seting up new candle
				c.UtcBegin = time.Now().Unix()
				c.O = c.C
				c.H = c.C
				c.L = c.C
				c.VolumeBS = poloniex.VBS{}
			}

			//make candle candles_minute_15
		case <-candlerMinute15.C:
			log.Println("candles_minute_15")
			//make 15-minute candles
			for _, s := range poloniex.Symbols {
				c := candles_minute_15[s]
				c.TimeFrame = "candles_minute_15"
				c.UtcEnd = time.Now().Unix()
				//store candle
				store.Single(c)
				//seting up new candle
				c.UtcBegin = time.Now().Unix()
				c.O = c.C
				c.H = c.C
				c.L = c.C
				c.VolumeBS = poloniex.VBS{}
			}

		case <-candlerHour1.C:
			log.Println("candles_hour_1")
			//make 1-hour candles
			for _, s := range poloniex.Symbols {
				c := candles_hour_1[s]
				c.TimeFrame = "candles_hour_1"
				c.UtcEnd = time.Now().Unix()
				//store candle
				store.Single(c)
				//seting up new candle
				c.UtcBegin = time.Now().Unix()
				c.O = c.C
				c.H = c.C
				c.L = c.C
				c.VolumeBS = poloniex.VBS{}
			}

		case <-candlerDay1.C:
			log.Println("candles_day_1")
			//make 1-hour candles
			for _, s := range poloniex.Symbols {
				c := candles_day_1[s]
				c.TimeFrame = "candles_day_1"
				c.UtcEnd = time.Now().Unix()
				//store candle
				store.Single(c)
				//seting up new candle
				c.UtcBegin = time.Now().Unix()
				c.O = c.C
				c.H = c.C
				c.L = c.C
				c.VolumeBS = poloniex.VBS{}
			}

			//close socket
		case <-interrupt:
			log.Println("Interrupt received, closing connection...")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("WriteMessage() error:", err)
				return
			}
			time.Sleep(time.Second)
			return

			//retrive data
		default:
			_, m, err := conn.ReadMessage()
			if err != nil {
				log.Println("ReadMessage() error:", err)
				return
			}

			var rt poloniex.RecentTrade
			if err := json.Unmarshal(m, &rt); err != nil {
				log.Println("Unmarshal() error:", err)
				continue
			}
			if rt.Channel == "trades" && len(rt.Data) > 0 {
				recent := rt.Data[0]

				c := candles_minute_1[recent.Symbol]
				c.Pair = recent.Symbol
				c.H = recent.HighPrice(c.H)
				c.L = recent.LowPrice(c.L)
				c.C = recent.HighPrice(0)

				c.VolumeBS.BuyBase += recent.BuyBase()
				c.VolumeBS.SellBase += recent.SellBase()

				c.VolumeBS.BuyQuote += recent.BuyQuote()
				c.VolumeBS.SellQuote += recent.SellQuote()

				candles_minute_1[recent.Symbol] = c

			}
		}

	}
}
