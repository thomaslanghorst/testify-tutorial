package calculations

import (
	"database/sql"
	"errors"
	"fmt"
	"testify-tutorial/mocks"
	"testify-tutorial/stocks"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "mysecretpassword"
	dbName     = "postgres"
)

// Integration tests

type IntTestSuite struct {
	suite.Suite
	db         *sql.DB
	calculator PriceIncreaseCalculator
}

func TestIntTestSuite(t *testing.T) {
	suite.Run(t, &IntTestSuite{})
}

func (its *IntTestSuite) SetupSuite() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		its.FailNowf("unable to connect to database", err.Error())
	}

	setupDatabase(its, db)

	pp := stocks.NewPriceProvider(db)
	calculator := NewPriceIncreaseCalculator(pp)

	its.db = db
	its.calculator = calculator
}

func (its *IntTestSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestCalculate_Error" {
		return
	}
	seedTestTable(its, its.db) // ts -> price=1, ts+1min -> price=2
}

func (its *IntTestSuite) TearDownSuite() {
	tearDownDatabase(its)
}

func (its *IntTestSuite) TearDownTest() {
	cleanTable(its)
}

func (its *IntTestSuite) TestCalculate_Error() {

	actual, err := its.calculator.PriceIncrease()

	its.EqualError(err, "not enough data")
	its.Equal(0.0, actual)

}

func (its *IntTestSuite) TestCalculate() {

	actual, err := its.calculator.PriceIncrease()

	its.Nil(err)
	its.Equal(100.0, actual)

}

// Helper functions

func setupDatabase(its *IntTestSuite, db *sql.DB) {
	its.T().Log("setting up database")

	_, err := db.Exec(`CREATE DATABASE stockprices_test`)
	if err != nil {
		its.FailNowf("unable to create database", err.Error())
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS stockprices (
		timestamp TIMESTAMPTZ PRIMARY KEY,
		price DECIMAL NOT NULL
	)`)

	if err != nil {
		its.FailNowf("unable to create table", err.Error())
	}

}

func seedTestTable(its *IntTestSuite, db *sql.DB) {
	its.T().Log("seeding test table")

	for i := 1; i <= 2; i++ {
		_, err := db.Exec("INSERT INTO stockprices (timestamp, price) VALUES ($1,$2)", time.Now().Add(time.Duration(i)*time.Minute), float64(i))
		if err != nil {
			its.FailNowf("unable to seed table", err.Error())
		}
	}
}

func cleanTable(its *IntTestSuite) {
	its.T().Log("cleaning database")

	_, err := its.db.Exec(`DELETE FROM stockprices`)
	if err != nil {
		its.FailNowf("unable to clean table", err.Error())
	}
}

func tearDownDatabase(its *IntTestSuite) {
	its.T().Log("tearing down database")

	_, err := its.db.Exec(`DROP TABLE stockprices`)
	if err != nil {
		its.FailNowf("unable to drop table", err.Error())
	}

	_, err = its.db.Exec(`DROP DATABASE stockprices_test`)
	if err != nil {
		its.FailNowf("unable to drop database", err.Error())
	}

	err = its.db.Close()
	if err != nil {
		its.FailNowf("unable to close database", err.Error())
	}
}

// Unit tests

type UnitTestSuite struct {
	suite.Suite
	calculator        PriceIncreaseCalculator
	priceProviderMock *mocks.PriceProvider
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &UnitTestSuite{})
}

func (uts *UnitTestSuite) SetupTest() {
	priceProviderMock := mocks.PriceProvider{}
	calculator := NewPriceIncreaseCalculator(&priceProviderMock)

	uts.calculator = calculator
	uts.priceProviderMock = &priceProviderMock
}

func (uts *UnitTestSuite) TestCalculate() {
	uts.priceProviderMock.On("List", mock.Anything).Return([]*stocks.PriceData{}, nil)

	actual, err := uts.calculator.PriceIncrease()

	uts.Equal(0.0, actual)
	uts.EqualError(err, "not enough data")
}

func (uts *UnitTestSuite) TestCalculate_ErrorFromPriceProvider() {
	expectedError := errors.New("oh my god")

	uts.priceProviderMock.On("List", mock.Anything).Return([]*stocks.PriceData{}, expectedError)

	actual, err := uts.calculator.PriceIncrease()

	uts.Equal(0.0, actual)
	uts.Equal(expectedError, err)

}
