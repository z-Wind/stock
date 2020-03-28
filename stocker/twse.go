package stocker

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/z-Wind/twse"
)

// TWSE stocker
type TWSE struct {
	Service *twse.Service
}

// NewTWSE 建立 twse Service
func NewTWSE() (*TWSE, error) {
	client := twse.GetClient()
	twse, err := twse.New(client)
	if err != nil {
		return nil, errors.Wrapf(err, "NewTWSE")
	}
	return &TWSE{Service: twse}, nil
}

func (t *TWSE) isOTC(symbol string) bool {
	pattern := regexp.MustCompile(`(?i)\.two$`)

	return pattern.MatchString(symbol)
}

func (t *TWSE) isListed(symbol string) bool {
	pattern := regexp.MustCompile(`(?i)\.tw$`)

	return pattern.MatchString(symbol)
}

// Quote 得到股票價格
func (t *TWSE) Quote(symbol string) (float64, error) {
	var quote *twse.StockInfo
	var err error
	switch {
	case t.isListed(symbol):
		id := strings.SplitN(symbol, ".", 2)[0]
		call := t.Service.Quotes.GetStockInfoTWSE(id)
		quote, err = call.Do()
		if err != nil {
			return 0, errors.Wrapf(err, "twse: GetStockInfoTWSE.Do")
		}
	case t.isOTC(symbol):
		id := strings.SplitN(symbol, ".", 2)[0]
		call := t.Service.Quotes.GetStockInfoOTC(id)
		quote, err = call.Do()
		if err != nil {
			return 0, errors.Wrapf(err, "twse: GetStockInfoOTC.Do")
		}
	default:
		return 0, ErrorNoSupport{fmt.Sprintf("%s is not supported by TWSE Quote", symbol)}
	}

	price := quote.TradePrice
	if quote.TradePrice == 0 {
		price = quote.YesterdayPrice
	}
	return float64(price), nil
}

func (t *TWSE) priceHistoryTWSE(symbol string) ([]*DatePrice, error) {
	id := strings.SplitN(symbol, ".", 2)[0]
	call := t.Service.Timeseries.MonthlyTWSE(id, time.Now())
	p, err := call.Do()
	if err != nil {
		return nil, errors.Wrapf(err, "twse: MonthlyTWSE.Do")
	}

	timeSeries := make([]*DatePrice, len(p.TimeSeries))
	for i, trade := range p.TimeSeries {
		t := DatePrice{
			Date:   Time(trade.Time),
			Open:   float64(trade.Open),
			High:   float64(trade.High),
			Low:    float64(trade.Low),
			Close:  float64(trade.Close),
			Volume: float64(trade.Volume),
		}
		timeSeries[i] = &t
	}

	return timeSeries, nil
}
func (t *TWSE) priceHistoryOTC(symbol string) ([]*DatePrice, error) {
	id := strings.SplitN(symbol, ".", 2)[0]
	call := t.Service.Timeseries.MonthlyOTC(id, time.Now())
	p, err := call.Do()
	if err != nil {
		return nil, errors.Wrapf(err, "twse: MonthlyOTC.Do")
	}

	timeSeries := make([]*DatePrice, len(p.TimeSeries))
	for i, trade := range p.TimeSeries {
		t := DatePrice{
			Date:   Time(trade.Time),
			Open:   float64(trade.Open),
			High:   float64(trade.High),
			Low:    float64(trade.Low),
			Close:  float64(trade.Close),
			Volume: float64(trade.Volume),
		}
		timeSeries[i] = &t
	}

	return timeSeries, nil
}

// PriceHistory 得到股票歷史價格
func (t *TWSE) PriceHistory(symbol string) ([]*DatePrice, error) {
	switch {
	case t.isListed(symbol):
		return t.priceHistoryTWSE(symbol)
	case t.isOTC(symbol):
		return t.priceHistoryOTC(symbol)
	default:
		return nil, ErrorNoSupport{fmt.Sprintf("%s is not supported by TWSE PriceHistory", symbol)}
	}
}

// PriceAdjHistory 得到股票歷史 Adj 價格
func (t *TWSE) PriceAdjHistory(symbol string) ([]*DatePrice, error) {
	return nil, ErrorNoSupport{"TWSE does not support PriceAdjHistory"}
}
