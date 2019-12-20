package alphavantage

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/z-Wind/stock/api/setting"

	"github.com/pkg/errors"
)

// AlphaVantage api
type AlphaVantage struct {
	key string
	qps int
}

func NewAlphaVantage(key string) (*AlphaVantage, error) {
	return &AlphaVantage{key: key, qps: 1}, nil
}

// Intraday https://www.alphavantage.co/documentation/#intraday
func (a *AlphaVantage) Intraday(symbol, interval, outputsize string) (*TimeSeries, error) {
	data := new(TimeSeries)
	paras := map[string]string{
		"function":   "TIME_SERIES_INTRADAY",
		"symbol":     symbol,
		"interval":   interval,
		"outputsize": outputsize,
		"datatype":   "json",
		"apikey":     a.key,
	}

	err := api(data, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return data, nil
}

// Daily https://www.alphavantage.co/documentation/#daily
func (a *AlphaVantage) Daily(symbol, outputsize string) (*TimeSeries, error) {
	data := new(TimeSeries)
	paras := map[string]string{
		"function":   "TIME_SERIES_DAILY",
		"symbol":     symbol,
		"outputsize": outputsize,
		"datatype":   "json",
		"apikey":     a.key,
	}

	err := api(data, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return data, nil
}

// DailyAdj https://www.alphavantage.co/documentation/#dailyadj
func (a *AlphaVantage) DailyAdj(symbol, outputsize string) (*TimeSeries, error) {
	data := new(TimeSeries)
	paras := map[string]string{
		"function":   "TIME_SERIES_DAILY_ADJUSTED",
		"symbol":     symbol,
		"outputsize": outputsize,
		"datatype":   "json",
		"apikey":     a.key,
	}

	err := api(data, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return data, nil
}

// Weekly https://www.alphavantage.co/documentation/#weekly
func (a *AlphaVantage) Weekly(symbol, outputsize string) (*TimeSeries, error) {
	data := new(TimeSeries)
	paras := map[string]string{
		"function":   "TIME_SERIES_WEEKLY",
		"symbol":     symbol,
		"outputsize": outputsize,
		"datatype":   "json",
		"apikey":     a.key,
	}

	err := api(data, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return data, nil
}

// WeeklyAdj https://www.alphavantage.co/documentation/#weeklyadj
func (a *AlphaVantage) WeeklyAdj(symbol, outputsize string) (*TimeSeries, error) {
	data := new(TimeSeries)
	paras := map[string]string{
		"function":   "TIME_SERIES_WEEKLY_ADJUSTED",
		"symbol":     symbol,
		"outputsize": outputsize,
		"datatype":   "json",
		"apikey":     a.key,
	}

	err := api(data, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return data, nil
}

// Monthly https://www.alphavantage.co/documentation/#monthly
func (a *AlphaVantage) Monthly(symbol, outputsize string) (*TimeSeries, error) {
	data := new(TimeSeries)
	paras := map[string]string{
		"function":   "TIME_SERIES_MONTHLY",
		"symbol":     symbol,
		"outputsize": outputsize,
		"datatype":   "json",
		"apikey":     a.key,
	}

	err := api(data, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return data, nil
}

// MonthlyAdj https://www.alphavantage.co/documentation/#monthlyadj
func (a *AlphaVantage) MonthlyAdj(symbol, outputsize string) (*TimeSeries, error) {
	data := new(TimeSeries)
	paras := map[string]string{
		"function":   "TIME_SERIES_MONTHLY_ADJUSTED",
		"symbol":     symbol,
		"outputsize": outputsize,
		"datatype":   "json",
		"apikey":     a.key,
	}

	err := api(data, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return data, nil
}

// LatestPrice https://www.alphavantage.co/documentation/#latestprice
func (a *AlphaVantage) LatestPrice(symbol string) (*LastTrade, error) {
	data := new(Quote)
	paras := map[string]string{
		"function": "GLOBAL_QUOTE",
		"symbol":   symbol,
		"datatype": "json",
		"apikey":   a.key,
	}

	err := api(data, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return data.LastPrice, nil
}

func api(data Data, paras map[string]string) error {
	url := "https://www.alphavantage.co/query"
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Wrap(err, "api")
	}
	req = setting.WrapRequest(req)

	q := req.URL.Query()
	for k, v := range paras {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	//log.Printf("request api url=%s", req.URL)

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "client.Do")
	}
	defer resp.Body.Close()
	defer log.Printf("alphavantage response %s:%s", resp.Status, req.URL)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "ioutil.ReadAll")
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return errors.Wrap(err, "json.Unmarshal")
	}

	return data.GetError()
}
