package models

import (
	"github.com/uadmin/uadmin"
)

type Ticker struct {
	uadmin.Model
	Name        string `uadmin:"required"`
	Symbol      string
	Type        TickerType `uadmin:"required"`
	Description string     `uadmin:"html"`
}
