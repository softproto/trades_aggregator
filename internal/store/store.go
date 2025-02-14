package store

import (
	"TradesAggregator/pkg/poloniex"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func Single(c poloniex.Kline) {
	// log.Println(c)
}

func Batch(c []poloniex.Kline) {

	// for _, b := range c {
	// 	log.Println(b)
	// }
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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
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
    CREATE TABLE IF NOT EXISTS candles (
        id SERIAL PRIMARY KEY,
        market VARCHAR(20),
        open FLOAT8,
        high FLOAT8,
        low FLOAT8,
        close FLOAT8,
        volume FLOAT8,
        timestamp TIMESTAMP
    );`

	_, err := db.Exec(createKlineTable)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Kline table created")
}
