package models

import (
	"encoding/json"
	"github.com/uadmin/rrd"
	"github.com/uadmin/uadmin"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var rrdOHLCTemplates map[time.Duration]rrd.RRD

func init() {
	rrdOHLCTemplates = map[time.Duration]rrd.RRD{}

	d, err := os.Open("templates/data")
	if err != nil {
		uadmin.Trail(uadmin.WARNING, "No RRD OHLC templates folder. %s", err)
		return
	}

	fList, err := d.Readdir(-1)
	if err != nil {
		uadmin.Trail(uadmin.WARNING, "Error reading OHLC templates folder. %s", err)
		return
	}

	for _, f := range fList {
		if f.IsDir() {
			continue
		}
		if strings.HasPrefix(f.Name(), "ohlc_") {
			durationRaw := strings.TrimPrefix(f.Name(), "ohlc_")
			durationRaw = strings.TrimSuffix(durationRaw, ".json")
			uadmin.Trail(uadmin.DEBUG, durationRaw)

			buf, err := ioutil.ReadFile("templates/data/" + f.Name())
			if err != nil {
				uadmin.Trail(uadmin.ERROR, "Error reading ohlc template (%s). %s", f.Name(), err)
				continue
			}
			d, err := strconv.ParseInt(durationRaw, 10, 64)
			if err != nil {
				uadmin.Trail(uadmin.ERROR, "Error parsing duration of ohlc template (%s). %s", f.Name(), err)
				continue
			}
			tmpl := rrd.RRD{}
			tmpl.Start = "-40y"
			err = json.Unmarshal(buf, &tmpl)
			if err != nil {
				uadmin.Trail(uadmin.ERROR, "Error parsing ohlc template (%s). %s", f.Name(), err)
				continue
			}
			rrdOHLCTemplates[time.Duration(d)*time.Second] = tmpl
		}
	}
	uadmin.Trail(uadmin.DEBUG, "%#v", rrdOHLCTemplates)
}

type OHLC struct {
	Date          time.Time
	Duration      time.Duration
	Open          float64
	High          float64
	Low           float64
	Close         float64
	Volume        float64
	AdjustedClose float64
	Dividend      float64
	Split         float64
}

func (o *OHLC) Save(filename string) {
	err := rrd.CreateRRD(filename, rrdOHLCTemplates[o.Duration])
	if err != nil {
		uadmin.Trail(uadmin.ERROR, "Error creating ohlc file. %s", err)
		return
	}
	err = rrd.UpdateRRDWithDate(filename, 8, &o.Date, o.Open, o.High, o.Low, o.Close, o.Volume, o.AdjustedClose, o.Dividend, o.Split)
	if err != nil {
		uadmin.Trail(uadmin.ERROR, "Unable to update ohlc data. %s", err)
		return
	}
}

func GetOHLC(filename string, from interface{}, to interface{}, step interface{}) []OHLC {
	ohlcList := []OHLC{}
	dateList, valueList, err := rrd.FetchRRD(filename, from, to, step)
	if err != nil {
		uadmin.Trail(uadmin.ERROR, "Unable to fetch OHLC. %s", err)
		return ohlcList
	}

	for i := range dateList {
		ohlcList = append(ohlcList, OHLC{
			Date:          dateList[i],
			Open:          valueList[i][0],
			High:          valueList[i][1],
			Low:           valueList[i][2],
			Close:         valueList[i][3],
			Volume:        valueList[i][4],
			AdjustedClose: valueList[i][5],
			Dividend:      valueList[i][6],
			Split:         valueList[i][7],
		})
	}

	return ohlcList
}

type OHLCByDate []OHLC

func (o OHLCByDate) Len() int           { return len(o) }
func (o OHLCByDate) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o OHLCByDate) Less(i, j int) bool { return o[i].Date.Before(o[j].Date) }
