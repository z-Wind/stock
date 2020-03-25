package stocker

import (
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
func (td *TDAmeritrade) Quote(symbol string) (float64, error) {
	call := td.Service.Quotes.GetQuote(symbol)
	quote, err := call.Do()
	if err != nil {
		if strings.Contains(err.Error(), "not be found") {
			return 0, ErrorNoFound{err.Error()}
		}

		return 0, errors.Wrapf(err, "gotd: QuoteEndpoint.Do")
	}

	return float64(quote.Mark), nil
}

// PriceHistory 得到股票歷史價格
func (td *TDAmeritrade) PriceHistory(symbol string) ([]*DatePrice, error) {
	call := td.Service.PriceHistory.GetPriceHistory(symbol)
	call.PeriodType(gotd.PriceHistoryPeriodTypeYear)
	call.Period(20)
	call.FrequencyType(gotd.PriceHistoryFrequencyTypeDaily)
	call.Frequency(1)
	p, err := call.Do()
	if err != nil {
		if strings.Contains(err.Error(), "not be found") {
			return nil, ErrorNoFound{err.Error()}
		}
		return nil, errors.WithMessage(err, "gotd: Daily.Do")
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

	return timeSeries, nil
}

// PriceAdjHistory 得到股票歷史 Adj 價格
func (td *TDAmeritrade) PriceAdjHistory(symbol string) ([]*DatePrice, error) {
	return nil, ErrorNoSupport{"TDAmeritrade does not support PriceAdjHistory"}
}