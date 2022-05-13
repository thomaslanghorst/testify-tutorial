# testify-tutorial

This repository contains the code for my testify tutorial [youtube video series](https://www.youtube.com/watch?v=Su6zn1_blw0&ab_channel=ThomasLanghorst). It shows how to use [testify](https://github.com/stretchr/testify) together with [mockery](https://github.com/mockery/mockery) to build awesome golang tests.

## Description

The tests are written for the `calculations` package and test the [PriceIncreaseCalculator](./calculations/priceIncrease.go) functionality which depends on [PriceProvider](./stocks/stocks.go) interface. 

## Prepare Postgres in Docker
If you want to run integration tests you need to have a running postgres instance. The constants that are being used are the default ones from [dockerhub](https://hub.docker.com/_/postgres)

```go
const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "mysecretpassword"
	dbName     = "postgres"
)
```

All you need to do is execute the following command and you should be good to go

```
docker run --name pg-testify -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 postgres
```

## Running the tests

```go
# run all tests
go test ./calculations

# only run unit tests
go test ./calculations -run UnitTestSuite -v

# only run integration tests
# running postgres instance required
go test ./calculations -run IntTestSuite -v
````

