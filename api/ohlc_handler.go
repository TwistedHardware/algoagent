package api

import (
	"fmt"
	"github.com/twistedhardware/algoagent/models"
	"github.com/uadmin/uadmin"
	"net/http"
	"strings"
)

func OHLCHandler(w http.ResponseWriter, r *http.Request) {
	symbol := strings.ToUpper(r.FormValue("s"))
	if symbol == "" {
		uadmin.ReturnJSON(w, r, map[string]interface{}{
			"status":  "error",
			"err_msg": "No symbol",
		})
		return
	}
	step := r.FormValue("step")
	if step == "" {
		step = "86400"
	}
	from := r.FormValue("from")
	if from == "" {
		from = "-3months"
	}
	to := r.FormValue("to")
	if to == "" {
		to = "now"
	}

	filename := fmt.Sprintf("data/ohlc/%s_%s.rrd", symbol, step)

	ohlc := models.GetOHLC(filename, from, to, step)

	uadmin.ReturnJSON(w, r, map[string]interface{}{
		"status":  "ok",
		"results": ohlc,
	})
}
