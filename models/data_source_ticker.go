package models

import (
	"github.com/uadmin/uadmin"
)

type OHLCDataSourceTicker struct {
	uadmin.Model
	OHLCDataSource   OHLCDataSource `uadmin:"required"`
	OHLCDataSourceID uint
	Ticker           Ticker `uadmin:"required"`
	TickerID         uint
	Symbol           string
}

func (OHLCDataSourceTicker) HideInDashboard() bool {
	return true
}
