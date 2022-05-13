package main

import (
	"database/sql"
	"fmt"
	"testify-tutorial/calculations"
	"testify-tutorial/stocks"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "mysecretpassword"
	dbName     = "postgres"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	createTable(db)
	seedTable(db) // needs only be executed once

	pp := stocks.NewPriceProvider(db)
	calculator := calculations.NewPriceIncreaseCalculator(pp)

	increase, err := calculator.PriceIncrease()
	if err != nil {
		panic(err)
	}

	fmt.Println(increase)
}

func createTable(db *sql.DB) {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS stockprices (
		timestamp TIMESTAMPTZ PRIMARY KEY,
		price DECIMAL NOT NULL
	)`)

	if err != nil {
		log.Fatalf("unable to prepare create table statement. Error %s", err.Error())
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("unable to execute create table statement. Error %s", err.Error())
	}
}

func seedTable(db *sql.DB) {
	log.Info("Seeding stockprices table")

	for i := 1; i <= 5; i++ {
		_, err := db.Exec("INSERT INTO stockprices (timestamp, price) VALUES ($1,$2)", time.Now().Add(time.Duration(-i)*time.Minute), float64((6-i)*5))
		if err != nil {
			log.Fatalf("unable to insert into stockprices table. Error %s", err.Error())
		}
	}

}
