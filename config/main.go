package config

import (
	_ "embed"

	"github.com/airtonix/bank-downloaders/schemas"
)

func init() {
	schemas.Initialize()
}
