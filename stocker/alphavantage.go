package stocker

import (
	"github.com/pkg/errors"
	"github.com/z-Wind/alphavantage"
)

// Alphavantage stocker
type Alphavantage struct {
	Service *alphavantage.Service
}

// NewAlphavantage 建立 alphavantage Service
func NewAlphavantage(apikey string) (*Alphavantage, error) {
	client := alphavantage.GetClient(apikey)
	av, err := alphavantage.New(client)
	if err != nil {
		return nil, errors.Wrapf(err, "NewAlphavantage")
	}
	return &Alphavantage{Service: av}, nil
}

// Quote 得到股票價格
func (av *Alphavantage) Quote(symbol string) (float64, error) {
	call := av.Service.TimeSeries.QuoteEndpoint(symbol)
	quote, err := call.Do()
	if err != nil {
		return 0, errors.Wrapf(err, "alphavantage: QuoteEndpoint.Do")
	}

	return float64(quote.Price), nil
}

// PriceHistory 得到股票歷史價格
func (av *Alphavantage) PriceHistory(symbol string) ([]*DatePrice, error) {
	call := av.Service.TimeSeries.Daily(symbol)
	call = call.Outputsize(alphavantage.OutputSizeFull)
	p, err := call.Do()
	if err != nil {
		return nil, errors.WithMessage(err, "alphavantage: Daily.Do")
	}

	timeSeries := make([]*DatePrice, len(p.TimeSeries))
	for i, trade := range p.TimeSeries {
		t := DatePrice{
			Date:   Time(trade.Time),
			Open:   trade.Open,
			High:   trade.High,
			Low:    trade.Low,
			Close:  trade.Close,
			Volume: trade.Volume,
		}
		timeSeries[i] = &t
	}

	return timeSeries, nil
}

// PriceAdjHistory 得到股票歷史 Adj 價格
func (av *Alphavantage) PriceAdjHistory(symbol string) ([]*DatePrice, error) {
	call := av.Service.TimeSeries.DailyAdj(symbol)
	call = call.Outputsize(alphavantage.OutputSizeFull)
	p, err := call.Do()
	if err != nil {
		return nil, errors.WithMessage(err, "alphavantage: DailyAdj.Do")
	}

	timeSeries := make([]*DatePrice, len(p.TimeSeries))
	for i, trade := range p.TimeSeries {
		t := DatePrice{
			Date:   Time(trade.Time),
			Open:   trade.Open,
			High:   trade.High,
			Low:    trade.Low,
			Close:  trade.Close,
			Volume: trade.Volume,
		}
		timeSeries[i] = &t
	}

	return timeSeries, nil
}
