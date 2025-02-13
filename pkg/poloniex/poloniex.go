package poloniex

import (
	"fmt"
	"log"
	"strconv"
)

var TradedSymbols = []string{"BTC_USDT", "TRX_USDT", "ETH_USDT", "DOGE_USDT", "BCH_USDT"}

const TradesChannel = "trades"

const PublicAPI = "wss://ws.poloniex.com/ws/public"

type Kline struct {
	Pair      string  // название пары в Bitsgap
	TimeFrame string  // период формирования свечи (1m, 1h, 1d
	O         float64 // open - цена открытия
	H         float64 // high - максимальная цена
	L         float64 // low - минимальная цена
	C         float64 // close - цена закрытия
	UtcBegin  int64   // время unix начала формирования свечки
	UtcEnd    int64   // время unix окончания формирования свечки
	VolumeBS  VBS
}

type VBS struct {
	BuyBase   float64 // объём покупок в базовой валюте
	SellBase  float64 // объём продаж в базовой валюте
	BuyQuote  float64 // объём покупок в котируемой валюте
	SellQuote float64 // объём продаж в котируемой валюте
}

type RecentTrade struct {
	Channel string  `json:"channel"`
	Data    []Trade `json:"data"`
	// Data    []struct {
	// 	Symbol     string `json:"symbol"`
	// 	Amount     string `json:"amount"`
	// 	Quantity   string `json:"quantity"`
	// 	TakerSide  string `json:"takerSide"`
	// 	CreateTime int64  `json:"createTime"`
	// 	Price      string `json:"price"`
	// 	ID         string `json:"id"`
	// 	Ts         int64  `json:"ts"`
	// } `json:"data"`
}

type Trade struct {
	Symbol     string `json:"symbol"`
	Amount     string `json:"amount"`
	Quantity   string `json:"quantity"`
	TakerSide  string `json:"takerSide"`
	CreateTime int64  `json:"createTime"`
	Price      string `json:"price"`
	ID         string `json:"id"`
	Ts         int64  `json:"ts"`
}

func (t *Trade) GetPrice() (float64, error) {

	n, err := strconv.ParseFloat(t.Price, 64)
	if err != nil {
		log.Println("ParseFloat() error:", err)
		return 0, err
	}

	return n, nil
}

func (t *Trade) GetAmount() (float64, error) {

	n, err := strconv.ParseFloat(t.Amount, 64)
	if err != nil {
		log.Println("ParseFloat() error:", err)
		return 0, err
	}

	return n, nil
}

func (t *Trade) GetQuantity() (float64, error) {

	n, err := strconv.ParseFloat(t.Quantity, 64)
	if err != nil {
		log.Println("ParseFloat() error:", err)
		return 0.0, err
	}

	return n, nil
}

func (t *Trade) GetSymbol() (string, error) {
	var err error

	if t.Symbol == "" {
		err = fmt.Errorf("invalid Symbol")
	}

	return t.Symbol, err
}

func (t *Trade) Timestamp() (int64, error) {
	var err error

	if t.Ts <= 0 {
		err = fmt.Errorf("invalid Timestamp")
	}

	return t.Ts, err
}
