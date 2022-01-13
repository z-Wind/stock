package stocker

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	yfinance "github.com/z-Wind/yahoofinance"
)

// YahooFinance stocker
type YahooFinance struct {
	Service *yfinance.Service
}

// NewYahooFinance 建立 yahoofinance Service
func NewYahooFinance() (*YahooFinance, error) {
	client := yfinance.GetClient()
	yfinance, err := yfinance.New(client)
	if err != nil {
		return nil, errors.Wrapf(err, "NewYahooFinance")
	}
	return &YahooFinance{Service: yfinance}, nil
}

// Quote 得到股票價格
func (yf *YahooFinance) Quote(ctx context.Context, symbol string) (float64, error) {
	call := yf.Service.Quote.RegularMarketPrice(symbol)
	call.Context(ctx)
	quote, err := call.Do()
	if err != nil {
		err = errors.Wrapf(err, "yfinance: Quote.Do")

		if strings.Contains(err.Error(), "not be found") {
			return 0, ErrorNoFound{err.Error()}
		}
		return 0, ErrorFatal{err.Error()}
	}

	return quote.Chart.Result[0].Meta.RegularMarketPrice, nil
}

// PriceHistory 得到股票歷史價格
func (yf *YahooFinance) PriceHistory(ctx context.Context, symbol string) ([]*DatePrice, error) {
	call := yf.Service.History.Between(symbol, time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC), time.Now())
	call.Context(ctx)
	p, err := call.Do()
	if err != nil {
		err = errors.Wrapf(err, "yfinance: Daily.Do")

		if strings.Contains(err.Error(), "Not Found") {
			return nil, ErrorNoFound{err.Error()}
		}
		return nil, ErrorFatal{err.Error()}
	}

	gmtoffset := p.Chart.Result[0].Meta.Gmtoffset
	timestamps := p.Chart.Result[0].Timestamp
	Open := p.Chart.Result[0].Indicators.Quote[0].Open
	High := p.Chart.Result[0].Indicators.Quote[0].High
	Low := p.Chart.Result[0].Indicators.Quote[0].Low
	Close := p.Chart.Result[0].Indicators.Quote[0].Close
	//CloseAdj := p.Chart.Result[0].Indicators.Adjclose[0].Value
	Volume := p.Chart.Result[0].Indicators.Quote[0].Volume

	timeSeries := make([]*DatePrice, len(timestamps))

	for i := range timestamps {
		t := DatePrice{
			Date:   Time(time.Unix(timestamps[i]+gmtoffset, 0)),
			Open:   Open[i],
			High:   High[i],
			Low:    Low[i],
			Close:  Close[i],
			Volume: Volume[i],
		}
		timeSeries[i] = &t
	}

	if len(timeSeries) == 0 {
		return nil, ErrorFatal{"yfinance: Empty List"}
	}

	return timeSeries, nil
}

// PriceAdjHistory 得到股票歷史 Adj 價格
func (yf *YahooFinance) PriceAdjHistory(ctx context.Context, symbol string) ([]*DatePrice, error) {
	call := yf.Service.History.Between(symbol, time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC), time.Now())
	call.Context(ctx)
	p, err := call.Do()
	if err != nil {
		err = errors.Wrapf(err, "yfinance: Daily.Do")

		if strings.Contains(err.Error(), "Not Found") {
			return nil, ErrorNoFound{err.Error()}
		}
		return nil, ErrorFatal{err.Error()}
	}

	gmtoffset := p.Chart.Result[0].Meta.Gmtoffset
	timestamps := p.Chart.Result[0].Timestamp
	Open := p.Chart.Result[0].Indicators.Quote[0].Open
	High := p.Chart.Result[0].Indicators.Quote[0].High
	Low := p.Chart.Result[0].Indicators.Quote[0].Low
	Close := p.Chart.Result[0].Indicators.Quote[0].Close
	CloseAdj, err := yf.toAdjHistory(timestamps, Close, p.Chart.Result[0].Events)
	if err != nil {
		return nil, ErrorFatal{"yfinance: toAdjHistory"}
	}
	Volume := p.Chart.Result[0].Indicators.Quote[0].Volume

	timeSeries := make([]*DatePrice, len(timestamps))

	for i := range timestamps {
		t := DatePrice{
			Date:     Time(time.Unix(timestamps[i]+gmtoffset, 0)),
			Open:     Open[i],
			High:     High[i],
			Low:      Low[i],
			Close:    Close[i],
			CloseAdj: CloseAdj[i],
			Volume:   Volume[i],
		}
		timeSeries[i] = &t
	}

	if len(timeSeries) == 0 {
		return nil, ErrorFatal{"yfinance: Empty List"}
	}

	return timeSeries, nil
}

// History to Adj History 將歷史價格轉換為歷史 Adj 價格
func (yf *YahooFinance) toAdjHistory(times []int64, history []float64, event yfinance.Events) ([]float64, error) {
	if len(times) != len(history) {
		return nil, ErrorFatal{"yfinance: the lengh of times and history are different"}
	}

	// splits := event.Splits
	divs := event.Dividends

	adj := make([]float64, len(history))
	ratio := 1.0
	pre_div := 0.0
	for i := len(times) - 1; i >= 0; i -= 1 {
		if pre_div != 0 {
			ratio *= (1.0 - pre_div/history[i])
			pre_div = 0.0
		}
		adj[i] = history[i] * ratio

		date := strconv.FormatInt(times[i], 10)
		// in yahoo, Close price adjusted for splits.
		// if split, ok := splits[date]; ok {
		// 	ratio *= (float64(split.Denominator) / float64(split.Numerator))
		// }
		if div, ok := divs[date]; ok {
			pre_div = div.Amount
		}
	}

	return adj, nil
}
