package poloniex

import (
	"fmt"
	"log"
	"strconv"
)

var Symbols = []string{"BTC_USDT", "TRX_USDT", "ETH_USDT", "DOGE_USDT", "BCH_USDT"}
var Intervals = []string{"MINUTE_1", "MINUTE_15", "HOUR_1", "DAY_1"}

const WssAPI = "wss://ws.poloniex.com/ws/public"
const TradesChannel = "trades"

const RestAPI = "https://api.poloniex.com/markets/"
const CandlesResource = "/candles"

type Candle struct {
	Low              string //lowest price over the interval
	High             string //highest price over the interval
	Open             string //price at the start time
	Close            string //price at the end time
	Amount           string //quote units traded over the interval
	Quantity         string //base units traded over the interval
	BuyTakerAmount   string //quote units traded over the interval filled by market buy orders
	BuyTakerQuantity string //base units traded over the interval filled by market buy orders
	TradeCount       int64  //count of trades
	Ts               int64  //time the record was pushed
	WeightedAverage  string //weighted average over the interval
	Interval         string //the selected interval
	StartTime        int64  //start time of interval
	CoseTime         int64  //close time of interval
}

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

func (t *Trade) getPrice() (float64, error) {

	n, err := strconv.ParseFloat(t.Price, 64)
	if err != nil {
		log.Println("ParseFloat() error:", err)
		return 0, err
	}

	return n, nil
}

// reurn max of v and price
func (t *Trade) HighPrice(v float64) float64 {
	p, err := t.getPrice()
	if err != nil {
		if p > v {
			return p
		}
		log.Println("Price value invalid:", err)
	}
	return v
}

// reurn min of v and price
func (t *Trade) LowPrice(v float64) float64 {
	p, err := t.getPrice()
	if err != nil {
		if p < v {
			return p
		}
		log.Println("Price value invalid:", err)
	}
	return v
}

func (t *Trade) getAmount() (float64, error) {

	n, err := strconv.ParseFloat(t.Amount, 64)
	if err != nil {
		log.Println("ParseFloat() error:", err)
		return 0, err
	}

	return n, nil
}

func (t *Trade) getQuantity() (float64, error) {

	n, err := strconv.ParseFloat(t.Quantity, 64)
	if err != nil {
		log.Println("ParseFloat() error:", err)
		return 0.0, err
	}

	return n, nil
}

func (t *Trade) BuyBase() float64 {
	a, err := t.getAmount()
	if err != nil {
		log.Println("amount value invalid:", err)
	}

	if t.TakerSide == "sell" {
		return 0
	}

	return a
}

func (t *Trade) SellBase() float64 {
	a, err := t.getAmount()
	if err != nil {
		log.Println("amount value invalid:", err)
	}

	if t.TakerSide == "sell" {
		return a
	}

	return 0
}

func (t *Trade) BuyQuote() float64 {
	q, err := t.getQuantity()
	if err != nil {
		log.Println("amount value invalid:", err)
	}

	if t.TakerSide == "sell" {
		return 0
	}

	return q

}

func (t *Trade) SellQuote() float64 {
	q, err := t.getQuantity()
	if err != nil {
		log.Println("amount value invalid:", err)
	}

	if t.TakerSide == "sell" {
		return q
	}

	return 0

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
