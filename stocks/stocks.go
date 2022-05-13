package stocks

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	dateFormat = "2006-01-02"
)

type PriceData struct {
	Timestamp time.Time
	Price     float64
}

type PriceProvider interface {
	Latest() (*PriceData, error)
	List(date time.Time) ([]*PriceData, error)
}

type priceProvider struct {
	db *sql.DB
}

func NewPriceProvider(db *sql.DB) PriceProvider {
	return &priceProvider{
		db: db,
	}
}

func (p *priceProvider) Latest() (*PriceData, error) {
	var priceData PriceData

	err := p.db.QueryRow("SELECT * FROM stockprices ORDER BY timestamp DESC limit 1").Scan(&priceData.Timestamp, &priceData.Price)
	if err != nil {
		return &priceData, fmt.Errorf("unable to query table. Error %s", err.Error())
	}

	return &priceData, nil

}

func (p *priceProvider) List(date time.Time) ([]*PriceData, error) {
	priceData := make([]*PriceData, 0)

	var rows *sql.Rows
	var err error

	rows, err = p.db.Query("SELECT * FROM stockprices where timestamp::date = $1 ORDER BY timestamp DESC", date.Format(dateFormat))

	if err != nil {
		return priceData, fmt.Errorf("unable to prepare SELECT statement. Error %s", err.Error())
	}

	var timestamp time.Time
	var price float64

	for rows.Next() {
		err = rows.Scan(&timestamp, &price)
		if err != nil {
			return priceData, fmt.Errorf("unable to query table. Error %s", err.Error())
		}

		priceData = append(priceData, &PriceData{
			Timestamp: timestamp,
			Price:     price,
		})
	}

	return priceData, nil

}
