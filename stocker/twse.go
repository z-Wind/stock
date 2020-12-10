package stocker

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/z-Wind/twse"
	"golang.org/x/time/rate"
)

// TWSE stocker
type TWSE struct {
	Service *twse.Service

	limit *rate.Limiter
}

// NewTWSE 建立 twse Service
func NewTWSE() (*TWSE, error) {
	client := twse.GetClient()
	twse, err := twse.New(client)
	if err != nil {
		return nil, errors.Wrapf(err, "NewTWSE")
	}
	return &TWSE{Service: twse, limit: rate.NewLimiter(rate.Every(time.Second*5/3), 1)}, nil
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
func (t *TWSE) Quote(ctx context.Context, symbol string) (float64, error) {
	if err := t.limit.Wait(ctx); err != nil {
		return 0, ErrorFatal{errors.Wrapf(err, "twse: limit.Wait").Error()}
	}
	select {
	case <-ctx.Done():
		return 0, ErrorFatal{errors.Wrapf(ctx.Err(), "twse").Error()}
	default:
	}

	var quote *twse.StockInfo
	var err error
	switch {
	case t.isListed(symbol):
		id := strings.SplitN(symbol, ".", 2)[0]
		call := t.Service.Quotes.GetStockInfoTWSE(id)
		call.Context(ctx)
		quote, err = call.Do()
		if err != nil {
			return 0, ErrorFatal{errors.Wrapf(err, "twse: GetStockInfoTWSE.Do").Error()}
		}
	case t.isOTC(symbol):
		id := strings.SplitN(symbol, ".", 2)[0]
		call := t.Service.Quotes.GetStockInfoOTC(id)
		call.Context(ctx)
		quote, err = call.Do()
		if err != nil {
			return 0, ErrorFatal{errors.Wrapf(err, "twse: GetStockInfoOTC.Do").Error()}
		}
	default:
		return 0, ErrorNoSupport{fmt.Sprintf("twse: Quote: %s is not supported", symbol)}
	}

	price := quote.TradePrice
	if quote.TradePrice == 0 {
		price = quote.YesterdayPrice
	}
	return float64(price), nil
}

func (t *TWSE) priceHistoryTWSE(ctx context.Context, symbol string) ([]*DatePrice, error) {
	id := strings.SplitN(symbol, ".", 2)[0]
	call := t.Service.Timeseries.MonthlyTWSE(id, time.Now())
	call.Context(ctx)
	p, err := call.Do()
	if err != nil {
		return nil, ErrorFatal{errors.Wrapf(err, "twse: MonthlyTWSE.Do").Error()}
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

	if len(timeSeries) == 0 {
		return nil, ErrorFatal{"Empty List"}
	}

	return timeSeries, nil
}
func (t *TWSE) priceHistoryOTC(ctx context.Context, symbol string) ([]*DatePrice, error) {
	id := strings.SplitN(symbol, ".", 2)[0]
	call := t.Service.Timeseries.MonthlyOTC(id, time.Now())
	call.Context(ctx)
	p, err := call.Do()
	if err != nil {
		return nil, ErrorFatal{errors.Wrapf(err, "twse: MonthlyOTC.Do").Error()}
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

	if len(timeSeries) == 0 {
		return nil, ErrorFatal{"Empty List"}
	}

	return timeSeries, nil
}

// PriceHistory 得到股票歷史價格
func (t *TWSE) PriceHistory(ctx context.Context, symbol string) ([]*DatePrice, error) {
	if err := t.limit.Wait(ctx); err != nil {
		return nil, ErrorFatal{errors.Wrapf(err, "twse: limit.Wait").Error()}
	}
	select {
	case <-ctx.Done():
		return nil, ErrorFatal{errors.Wrapf(ctx.Err(), "twse").Error()}
	default:
	}

	switch {
	case t.isListed(symbol):
		return t.priceHistoryTWSE(ctx, symbol)
	case t.isOTC(symbol):
		return t.priceHistoryOTC(ctx, symbol)
	default:
		return nil, ErrorNoSupport{fmt.Sprintf("twse: PriceHistory: %s is not supported", symbol)}
	}
}

// PriceAdjHistory 得到股票歷史 Adj 價格
func (t *TWSE) PriceAdjHistory(ctx context.Context, symbol string) ([]*DatePrice, error) {
	return nil, ErrorNoSupport{"TWSE does not support PriceAdjHistory"}
}
