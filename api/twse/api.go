package twse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/z-Wind/stock/api/setting"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// TWSE api
type TWSE struct {
	qps int
}

func NewTWSE() (*TWSE, error) {
	return &TWSE{1}, nil
}

func renameSymbol(symbol string) string {
	symbol = strings.ToLower(symbol)
	if !strings.Contains(symbol, ".tw") {
		symbol = fmt.Sprintf("%s.tw", symbol)
	}
	return fmt.Sprintf("tse_%s", strings.ToLower(symbol))
}

// history Monthly
func (a *TWSE) HistoryMonthly(symbol string, year, month int) ([]Trade, error) {
	//http://www.twse.com.tw/exchangeReport/STOCK_DAY?response=json&date=20181230&stockNo=0050
	data := new(QuoteHistory)
	t := time.Date(year, time.Month(month), 30, 0, 0, 0, 0, time.Local)
	paras := map[string]string{
		"response": "json",
		"date":     t.Format("20060102"),
		"stockNo":  symbol,
	}

	err := api(data, getJSONHistory, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	trades := []Trade{}
	for _, t := range data.Data {
		volume, err := strconv.ParseFloat(strings.Replace(t[1], ",", "", -1), 64)
		open, err := strconv.ParseFloat(t[3], 64)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("strconv.ParseFloat(%v)", t))
		}
		dayHigh, err := strconv.ParseFloat(t[4], 64)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("strconv.ParseFloat(%v)", t))
		}
		dayLow, err := strconv.ParseFloat(t[5], 64)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("strconv.ParseFloat(%v)", t))
		}
		tradePrice, err := strconv.ParseFloat(t[6], 64)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("strconv.ParseFloat(%v)", t))
		}
		trades = append(trades, Trade{
			Date:           t[0],
			Symbol:         symbol,
			TemporalVolume: volume,
			Open:           open,
			DayHigh:        dayHigh,
			DayLow:         dayLow,
			TradePrice:     tradePrice,
		})
	}

	return trades, nil

}

// Quote get last price
func (a *TWSE) Quote(symbol string) ([]Trade, error) {
	//http://mis.twse.com.tw/stock/api/getStockInfo.jsp?ex_ch=tse_1101.tw&json=1&delay=0&_=1539865363091
	data := new(Message)
	millis := time.Now().UnixNano() / int64(time.Millisecond)
	paras := map[string]string{
		"ex_ch": renameSymbol(symbol),
		"json":  "1",
		"delay": "0",
		"_":     fmt.Sprintf("%d", millis),
	}

	err := api(data, getJSONQuote, paras)
	if err != nil {
		return nil, errors.WithMessage(err, "api")
	}

	return data.Trades, nil
}

func getCookies() ([]*http.Cookie, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", indexURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}
	req = setting.WrapRequest(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do")
	}
	defer resp.Body.Close()

	return resp.Cookies(), nil
}
func getJSON(url string, paras map[string]string) ([]byte, error) {
	cookies, err := getCookies()
	if err != nil {
		return nil, errors.WithMessage(err, "getCookies")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}
	req = setting.WrapRequest(req)

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	q := req.URL.Query()
	for k, v := range paras {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	//log.Printf("request api url=%s", req.URL)

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do")
	}
	defer resp.Body.Close()
	defer log.Printf("twse response %s:%s", resp.Status, req.URL)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll")
	}

	return body, nil
}
func getJSONHistory(paras map[string]string) ([]byte, error) {
	return getJSON(historyURL, paras)
}
func getJSONQuote(paras map[string]string) ([]byte, error) {
	return getJSON(queryURL, paras)
}

func parseJSON(in []byte, out Data) error {
	err := json.Unmarshal(in, out)

	return errors.Wrap(err, "json.Unmarshal")
}

func api(data Data, getJSONFunc func(paras map[string]string) ([]byte, error), paras map[string]string) error {
	body, err := getJSONFunc(paras)
	if err != nil {
		return errors.WithMessage(err, "getJSON")
	}

	err = parseJSON(body, data)
	if err != nil {
		return errors.WithMessage(err, "parseJSON")
	}

	return data.GetError()
}
