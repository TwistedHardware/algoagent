package main

import (
	"github.com/twistedhardware/algoagent/api"
	"github.com/twistedhardware/algoagent/models"
	"github.com/uadmin/uadmin"
	"net/http"
)

func main() {
	uadmin.Register(
		models.Ticker{},
		models.OHLCDataSource{},
		models.OHLCDataSourceTicker{},
	)

	uadmin.RegisterInlines(models.OHLCDataSource{}, map[string]string{
		"OHLCDataSourceTicker": "OHLCDataSourceID",
	})

	uadmin.RootURL = "/console/"
	uadmin.SiteName = "Algo Agent"

	http.HandleFunc("/", uadmin.Handler(api.APIHandler))

	uadmin.StartServer()
}
