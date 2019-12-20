package alphavantage

import (
	"errors"
	"fmt"
)

const (
	OutputSizeFull    = "full"
	OutputsizeCompact = "compact"
)

type Data interface {
	GetError() error
}

type Message struct {
	ErrorMessage string `json:"Error Message"`
	Information  string `json:"Information"`
	Note         string `json:"Note"`
}

type TimeSeries struct {
	Message

	MetaData             *MetaData         `json:"Meta Data"`
	TimeSeries1min       map[string]*Trade `json:"Time Series (1min)"`
	TimeSeries5min       map[string]*Trade `json:"Time Series (5min)"`
	TimeSeries15min      map[string]*Trade `json:"Time Series (15min)"`
	TimeSeries30min      map[string]*Trade `json:"Time Series (30min)"`
	TimeSeries60min      map[string]*Trade `json:"Time Series (60min)"`
	TimeSeriesDaily      map[string]*Trade `json:"Time Series (Daily)"`
	TimeSeriesWeekly     map[string]*Trade `json:"Weekly Time Series"`
	TimeSeriesWeeklyAdj  map[string]*Trade `json:"Weekly Adjusted Time Series"`
	TimeSeriesMonthly    map[string]*Trade `json:"Monthly Time Series"`
	TimeSeriesMonthlyAdj map[string]*Trade `json:"Monthly Adjusted Time Series"`
}

func (t *TimeSeries) GetError() error {
	if t.ErrorMessage != "" {
		return fmt.Errorf("%s", t.ErrorMessage)
	}
	if t.Information != "" {
		return fmt.Errorf("%s", t.Information)
	}
	if t.Note != "" {
		return fmt.Errorf("%s", t.Note)
	}
	return nil
}

type MetaData struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
	TimeZone      string `json:"5. Time Zone"`
	TimeZone_     string `json:"4. Time Zone"`
}

type Trade struct {
	Open             float64 `json:"1. open,string"`
	High             float64 `json:"2. high,string"`
	Low              float64 `json:"3. low,string"`
	Close            float64 `json:"4. close,string"`
	Volume           float64 `json:"5. volume,string"`
	Volume_          float64 `json:"6. volume,string"` //adj
	AdjustedClose    float64 `json:"5. adjusted close,string"`
	DividendAmount   float64 `json:"7. dividend amount,string"`
	SplitCoefficient float64 `json:"8. split coefficient,string"`
}

type Quote struct {
	Message
	LastPrice *LastTrade `json:"Global Quote"`
}

type LastTrade struct {
	Symbol           string  `json:"01. symbol"`
	Open             float64 `json:"02. open,string"`
	High             float64 `json:"03. high,string"`
	Low              float64 `json:"04. low,string"`
	Price            float64 `json:"05. price,string"`
	Volume           float64 `json:"06. volume,string"`
	LatestTradingDay string  `json:"07. latest trading day"`
	PreviousClose    float64 `json:"08. previous close,string"`
	Change           float64 `json:"09. change,string"`
	ChangePercent    string  `json:"10. change percent"`
}

func (q *Quote) GetError() error {
	if q.ErrorMessage != "" {
		return fmt.Errorf("%s", q.ErrorMessage)
	}
	if q.Information != "" {
		return fmt.Errorf("%s", q.Information)
	}
	if q.Note != "" {
		return fmt.Errorf("%s", q.Note)
	}
	if q.LastPrice.Symbol == "" {
		return errors.New("No Data")
	}
	return nil
}
