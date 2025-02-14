package store

import (
	"TradesAggregator/pkg/poloniex"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func Single(Kline poloniex.Kline) {

	query := `
		INSERT INTO kline (pair, timeframe, open, high, low, close, utcbegin, utcend, buybase, sellbase, buyquote, sellquote)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := db.Exec(query,
		Kline.Pair, Kline.TimeFrame, Kline.O, Kline.H, Kline.L, Kline.C,
		time.Unix(Kline.UtcBegin, 0), time.Unix(Kline.UtcEnd, 0),
		Kline.VolumeBS.BuyBase, Kline.VolumeBS.SellBase, Kline.VolumeBS.BuyQuote, Kline.VolumeBS.SellQuote)
	if err != nil {
		log.Println("Error saving kline:", err)
	}

}

func Batch(Klines []poloniex.Kline) {

	for _, Kline := range Klines {
		Single(Kline)
	}
}

const (
	host     = "localhost"
	port     = 5432
	user     = "posgresql"
	password = "posgresql"
	dbname   = "kline"
)

var db *sql.DB

func InitDB() {
	p := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", p)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: connected")
}

func CreateTable() {

	createKlineTable := `
    CREATE TABLE IF NOT EXISTS kline (
        id SERIAL PRIMARY KEY,
        pair VARCHAR(32),
		timeframe VARCHAR(32),
        open FLOAT8,
        high FLOAT8,
        low FLOAT8,
        close FLOAT8,
		utcbegin TIMESTAMP,
		utcend TIMESTAMP,
        buybase FLOAT8,
		sellbase FLOAT8,
		buyquote FLOAT8,
		sellquote FLOAT8
    );`

	_, err := db.Exec(createKlineTable)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("kline table created")
}
