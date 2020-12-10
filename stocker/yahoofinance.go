package stocker

import (
	"context"
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
	call := yf.Service.History.Period(symbol, "max", "1d")
	call.Context(ctx)
	p, err := call.Do()
	if err != nil {
		err = errors.Wrapf(err, "yfinance: Daily.Do")

		if strings.Contains(err.Error(), "Not Found") {
			return nil, ErrorNoFound{err.Error()}
		}
		return nil, ErrorFatal{err.Error()}
	}

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
			Date:   Time(time.Unix(timestamps[i], 0)),
			Open:   Open[i],
			High:   High[i],
			Low:    Low[i],
			Close:  Close[i],
			Volume: Volume[i],
		}
		timeSeries[i] = &t
	}

	if len(timeSeries) == 0 {
		return nil, ErrorFatal{"Empty List"}
	}

	return timeSeries, nil
}

// PriceAdjHistory 得到股票歷史 Adj 價格
func (yf *YahooFinance) PriceAdjHistory(ctx context.Context, symbol string) ([]*DatePrice, error) {
	call := yf.Service.History.Period(symbol, "max", "1d")
	call.Context(ctx)
	p, err := call.Do()
	if err != nil {
		err = errors.Wrapf(err, "yfinance: Daily.Do")

		if strings.Contains(err.Error(), "Not Found") {
			return nil, ErrorNoFound{err.Error()}
		}
		return nil, ErrorFatal{err.Error()}
	}

	timestamps := p.Chart.Result[0].Timestamp
	Open := p.Chart.Result[0].Indicators.Quote[0].Open
	High := p.Chart.Result[0].Indicators.Quote[0].High
	Low := p.Chart.Result[0].Indicators.Quote[0].Low
	Close := p.Chart.Result[0].Indicators.Quote[0].Close
	CloseAdj := p.Chart.Result[0].Indicators.Adjclose[0].Value
	Volume := p.Chart.Result[0].Indicators.Quote[0].Volume

	timeSeries := make([]*DatePrice, len(timestamps))

	for i := range timestamps {
		t := DatePrice{
			Date:     Time(time.Unix(timestamps[i], 0)),
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
		return nil, ErrorFatal{"Empty List"}
	}

	return timeSeries, nil
}
