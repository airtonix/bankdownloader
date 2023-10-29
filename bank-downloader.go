package main

import (
	"github.com/airtonix/bank-downloaders/cmd"
	"github.com/shopspring/decimal"
)

func main() {
	decimal.MarshalJSONWithoutQuotes = true
	cmd.Execute()
}
