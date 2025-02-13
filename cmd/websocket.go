package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"TradesAggregator/pkg/poloniex"
)

func wsConnect() {
	conn, _, err := websocket.DefaultDialer.Dial(poloniex.PublicAPI, nil)
	if err != nil {
		log.Fatal("Dial() error:", err)
	}
	defer conn.Close()

	//make quoted string from slice
	var builder strings.Builder
	for i, s := range poloniex.TradedSymbols {
		builder.WriteString(`"` + s + `"`)
		if i < len(poloniex.TradedSymbols)-1 {
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

	// candle aggregator timer (1-minute interval)
	candler := time.NewTicker(60 * time.Second)
	defer candler.Stop()

	// gracefull interruption
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// just a trick to create the first candle for each pair
	// in real life we ​​don't create candles at runtime
	// instead we save trades (tick data) to the database
	// and generate candles on demand
	candle := make(map[string]poloniex.Kline)
	for _, s := range poloniex.TradedSymbols {
		candle[s] = poloniex.Kline{
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

			//make candle
		case <-candler.C:
			fmt.Println("Candler", candle)
			fmt.Println(candle)
			//clean map
			candle = make(map[string]poloniex.Kline)

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

				// candle[recent.Data[0].Symbol] =
				p, _ := recent.GetPrice()
				s, _ := recent.GetSymbol()
				a, _ := recent.GetAmount()

				fmt.Println("--", s, p, a)

			}
		}

	}
}
