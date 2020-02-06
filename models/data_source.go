package models

import (
	"encoding/csv"
	"fmt"
	"github.com/uadmin/uadmin"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type OHLCDataSource struct {
	uadmin.Model
	Name              string             `uadmin:"required"`
	URL               string             `uadmin:"required;list_exclude"`
	Username          string             `uadmin:"list_exclude"`
	Password          string             `uadmin:"password;list_exclude"`
	Format            DataSourceFormat   `uadmin:"list_exclude"`
	Duration          DataSourceDuration `uadmin:"required"`
	DateFormat        string             `uadmin:"list_exclude"`
	DateName          string             `uadmin:"list_exclude"`
	OpenName          string             `uadmin:"list_exclude"`
	HighName          string             `uadmin:"list_exclude"`
	LowName           string             `uadmin:"list_exclude"`
	CloseName         string             `uadmin:"list_exclude"`
	VolumeName        string             `uadmin:"list_exclude"`
	AdjustedCloseName string             `uadmin:"list_exclude"`
	DividendName      string             `uadmin:"list_exclude"`
	SplitName         string             `uadmin:"list_exclude"`
	Fetch             string             `uadmin:"link"`
}

func (d *OHLCDataSource) Save() {
	uadmin.Save(d)

	if d.Fetch != fmt.Sprintf("%sapi/d/ohlcdatasource/method/FetchData/%d?$next=$back", uadmin.RootURL, d.ID) {
		d.Fetch = fmt.Sprintf("%sapi/d/ohlcdatasource/method/FetchData/%d?$next=$back", uadmin.RootURL, d.ID)
		uadmin.Save(d)
	}
}

func (d *OHLCDataSource) FetchData() {
	tickers := []OHLCDataSourceTicker{}
	uadmin.Filter(&tickers, "ohlc_data_source_id = ?", d.ID)

	d.URL = d.URL
	d.URL = strings.ReplaceAll(d.URL, "{DataSource.Username}", d.Username)
	d.URL = strings.ReplaceAll(d.URL, "{DataSource.Password}", d.Password)

	for _, t := range tickers {
		uadmin.Preload(&t)

		if t.Symbol == "" {
			t.Symbol = t.Ticker.Symbol
		}
		URL := strings.ReplaceAll(d.URL, "{Ticker.Symbol}", t.Symbol)

		uadmin.Trail(uadmin.DEBUG, URL)
		res, err := http.Get(URL)
		if err != nil {
			uadmin.Trail(uadmin.WARNING, "Error fetching data. %s", err)
			continue
		}

		ohlcList, _ := d.parseResults(res.Body)

		// Sort data from old to new
		sort.Sort(OHLCByDate(ohlcList))

		// Store data to rrd
		for _, ohlc := range ohlcList {
			uadmin.Trail(uadmin.DEBUG, "Saving %d", ohlc.Date.Unix())
			ohlc.Save(fmt.Sprintf("data/ohlc/%s_%d.rrd", t.Symbol, d.Duration))
		}
		//uadmin.Trail(uadmin.DEBUG, ohlcList)
	}
}

func (d *OHLCDataSource) parseResults(source io.Reader) ([]OHLC, error) {
	ohlcList := []OHLC{}
	if d.Format == d.Format.JSON() {
	}
	if d.Format == d.Format.CSVWithHeader() {
		colMap := map[string]int{}
		r := csv.NewReader(source)

		for i := 0; ; i++ {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return ohlcList, err
			}

			// If header parse it into the map
			if i == 0 {
				for index, v := range record {
					colMap[v] = index
				}
				continue
			}
			ohlc := OHLC{}
			ohlc.Date, _ = time.Parse(d.DateFormat, record[colMap[d.DateName]])
			ohlc.Open, _ = strconv.ParseFloat(record[colMap[d.OpenName]], 64)
			ohlc.High, _ = strconv.ParseFloat(record[colMap[d.HighName]], 64)
			ohlc.Low, _ = strconv.ParseFloat(record[colMap[d.LowName]], 64)
			ohlc.Close, _ = strconv.ParseFloat(record[colMap[d.CloseName]], 64)
			ohlc.Volume, _ = strconv.ParseFloat(record[colMap[d.VolumeName]], 64)
			ohlc.AdjustedClose, _ = strconv.ParseFloat(record[colMap[d.AdjustedCloseName]], 64)
			ohlc.Dividend, _ = strconv.ParseFloat(record[colMap[d.DividendName]], 64)
			ohlc.Split, _ = strconv.ParseFloat(record[colMap[d.SplitName]], 64)
			ohlc.Duration = time.Duration(int64(d.Duration)) * time.Second
			ohlcList = append(ohlcList, ohlc)
		}
	}
	return ohlcList, nil
}
