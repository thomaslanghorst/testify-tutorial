package calculations

import (
	"errors"
	"testify-tutorial/stocks"
	"time"
)

type PriceIncreaseCalculator interface {
	PriceIncrease() (float64, error)
}

type priceIncreaseCalculator struct {
	PriceProvider stocks.PriceProvider
}

func NewPriceIncreaseCalculator(pp stocks.PriceProvider) PriceIncreaseCalculator {
	return &priceIncreaseCalculator{
		PriceProvider: pp,
	}
}

func (pic *priceIncreaseCalculator) PriceIncrease() (float64, error) {

	prices, err := pic.PriceProvider.List(time.Now())
	if err != nil {
		return 0.0, err
	}

	if len(prices) < 2 {
		return 0.0, errors.New("not enough data")
	}

	return (prices[0].Price/prices[1].Price - 1.0) * 100.0, nil
}
