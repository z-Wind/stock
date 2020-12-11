package stocker

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/z-Wind/gotd"
)

// TDAmeritrade stocker
type TDAmeritrade struct {
	Service *gotd.Service
}

// NewTDAmeritrade 建立 gotd Service
func NewTDAmeritrade(clientsecretPath, tokenFile string) (*TDAmeritrade, error) {
	auth := gotd.NewAuth()
	client := auth.GetClient(clientsecretPath, tokenFile)
	td, err := gotd.New(client)
	if err != nil {
		return nil, errors.Wrapf(err, "NewTDAmeritrade")
	}
	return &TDAmeritrade{Service: td}, nil
}

// NewTDAmeritradeTLS 建立 gotd Service
func NewTDAmeritradeTLS(clientsecretPath, tokenFile, TLSCertPath, TLSKeyPath string) (*TDAmeritrade, error) {
	auth := gotd.NewAuth()
	auth.SetTLS(TLSCertPath, TLSKeyPath)
	client := auth.GetClient(clientsecretPath, tokenFile)
	td, err := gotd.New(client)
	if err != nil {
		return nil, errors.Wrapf(err, "NewTDAmeritrade")
	}
	return &TDAmeritrade{Service: td}, nil
}

// Quote 得到股票價格
func (td *TDAmeritrade) Quote(ctx context.Context, symbol string) (float64, error) {
	call := td.Service.Quotes.GetQuote(symbol)
	call.Context(ctx)
	quote, err := call.Do()
	if err != nil {
		err = errors.Wrapf(err, "gotd: QuoteEndpoint.Do")

		if strings.Contains(err.Error(), "not be found") {
			return 0, ErrorNoFound{err.Error()}
		}

		return 0, ErrorFatal{err.Error()}
	}

	return float64(quote.Mark), nil
}

// PriceHistory 得到股票歷史價格
func (td *TDAmeritrade) PriceHistory(ctx context.Context, symbol string) ([]*DatePrice, error) {
	call := td.Service.PriceHistory.GetPriceHistory(symbol)
	call.Context(ctx)
	call.PeriodType(gotd.PriceHistoryPeriodTypeYear)
	call.Period(20)
	call.FrequencyType(gotd.PriceHistoryFrequencyTypeDaily)
	call.Frequency(1)
	p, err := call.Do()
	if err != nil {
		err = errors.Wrapf(err, "gotd: Daily.Do")

		if strings.Contains(err.Error(), "not be found") {
			return nil, ErrorNoFound{err.Error()}
		}
		return nil, ErrorFatal{err.Error()}
	}

	timeSeries := make([]*DatePrice, len(p.Candles))
	for i, trade := range p.Candles {
		t := DatePrice{
			Date:   Time(trade.Datetime),
			Open:   trade.Open,
			High:   trade.High,
			Low:    trade.Low,
			Close:  trade.Close,
			Volume: trade.Volume,
		}
		timeSeries[i] = &t
	}

	if len(timeSeries) == 0 {
		return nil, ErrorFatal{"gotd: Empty List"}
	}

	return timeSeries, nil
}

// PriceAdjHistory 得到股票歷史 Adj 價格
func (td *TDAmeritrade) PriceAdjHistory(ctx context.Context, symbol string) ([]*DatePrice, error) {
	return nil, ErrorNoSupport{"TDAmeritrade does not support PriceAdjHistory"}
}
