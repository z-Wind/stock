package types

import (
	"encoding/json"
	"fmt"
)

type Quotes map[string]*Quote

// func (q *Quotes) Get(key string) *Quote {
// 	return map[string]*Quote(*q)[key]
// }

// func (q *Quotes) Len() int {
// 	return len(map[string]*Quote(*q))
// }

// func (q *Quotes) Iterate() map[string]*Quote {
// 	return map[string]*Quote(*q)
// }

type Quote struct {
	AssetType string `json:"assetType"`

	// MutualFund
	// Future
	// FutureOptions
	// Index
	// Option
	// Forex
	// ETF
	// Equity
	Data interface{}
}

type _Quote Quote

func (q *Quote) UnmarshalJSON(bs []byte) (err error) {
	quote := _Quote{}

	err = json.Unmarshal(bs, &quote)
	if err != nil {
		return err
	}

	switch quote.AssetType {
	case "MutualFund":
		quote.Data = &MutualFundQ{}
	case "Future":
		quote.Data = &Future{}
	case "FutureOptions":
		quote.Data = &FutureOptions{}
	case "Index":
		quote.Data = &Index{}
	case "Option":
		quote.Data = &OptionQ{}
	case "Forex":
		quote.Data = &Forex{}
	case "ETF":
		quote.Data = &ETF{}
	case "Equity":
		quote.Data = &EquityQ{}
	default:
		return fmt.Errorf("Not support type %s", quote.AssetType)
	}
	err = json.Unmarshal(bs, quote.Data)
	if err != nil {
		return err
	}

	*q = Quote(quote)

	return nil
}

type MutualFundQ struct {
	Symbol          string  `json:"symbol"`
	Description     string  `json:"description"`
	ClosePrice      float64 `json:"closePrice"`
	NetChange       float64 `json:"netChange"`
	TotalVolume     float64 `json:"totalVolume"`
	TradeTimeInLong float64 `json:"tradeTimeInLong"`
	Exchange        string  `json:"exchange"`
	ExchangeName    string  `json:"exchangeName"`
	Digits          float64 `json:"digits"`
	Wk52High        float64 `json:"52WkHigh"`
	Wk52Low         float64 `json:"52WkLow"`
	NAV             float64 `json:"nAV"`
	PeRatio         float64 `json:"peRatio"`
	DivAmount       float64 `json:"divAmount"`
	DivYield        float64 `json:"divYield"`
	DivDate         string  `json:"divDate"`
	SecurityStatus  string  `json:"securityStatus"`
}

type Future struct {
	Symbol                string  `json:"symbol"`
	BidPriceInDouble      float64 `json:"bidPriceInDouble"`
	AskPriceInDouble      float64 `json:"askPriceInDouble"`
	LastPriceInDouble     float64 `json:"lastPriceInDouble"`
	BidID                 string  `json:"bidId"`
	AskID                 string  `json:"askId"`
	HighPriceInDouble     float64 `json:"highPriceInDouble"`
	LowPriceInDouble      float64 `json:"lowPriceInDouble"`
	ClosePriceInDouble    float64 `json:"closePriceInDouble"`
	Exchange              string  `json:"exchange"`
	Description           string  `json:"description"`
	LastID                string  `json:"lastId"`
	OpenPriceInDouble     float64 `json:"openPriceInDouble"`
	ChangeInDouble        float64 `json:"changeInDouble"`
	FuturePercentChange   float64 `json:"futurePercentChange"`
	ExchangeName          string  `json:"exchangeName"`
	SecurityStatus        string  `json:"securityStatus"`
	OpenInterest          float64 `json:"openInterest"`
	Mark                  float64 `json:"mark"`
	Tick                  float64 `json:"tick"`
	TickAmount            float64 `json:"tickAmount"`
	Product               string  `json:"product"`
	FuturePriceFormat     string  `json:"futurePriceFormat"`
	FutureTradingHours    string  `json:"futureTradingHours"`
	FutureIsTradable      bool    `json:"futureIsTradable"`
	FutureMultiplier      float64 `json:"futureMultiplier"`
	FutureIsActive        bool    `json:"futureIsActive"`
	FutureSettlementPrice float64 `json:"futureSettlementPrice"`
	FutureActiveSymbol    string  `json:"futureActiveSymbol"`
	FutureExpirationDate  string  `json:"futureExpirationDate"`
}

type FutureOptions struct {
	Symbol                      string  `json:"symbol"`
	BidPriceInDouble            float64 `json:"bidPriceInDouble"`
	AskPriceInDouble            float64 `json:"askPriceInDouble"`
	LastPriceInDouble           float64 `json:"lastPriceInDouble"`
	HighPriceInDouble           float64 `json:"highPriceInDouble"`
	LowPriceInDouble            float64 `json:"lowPriceInDouble"`
	ClosePriceInDouble          float64 `json:"closePriceInDouble"`
	Description                 string  `json:"description"`
	OpenPriceInDouble           float64 `json:"openPriceInDouble"`
	NetChangeInDouble           float64 `json:"netChangeInDouble"`
	OpenInterest                float64 `json:"openInterest"`
	ExchangeName                string  `json:"exchangeName"`
	SecurityStatus              string  `json:"securityStatus"`
	Volatility                  float64 `json:"volatility"`
	MoneyIntrinsicValueInDouble float64 `json:"moneyIntrinsicValueInDouble"`
	MultiplierInDouble          float64 `json:"multiplierInDouble"`
	Digits                      float64 `json:"digits"`
	StrikePriceInDouble         float64 `json:"strikePriceInDouble"`
	ContractType                string  `json:"contractType"`
	Underlying                  string  `json:"underlying"`
	TimeValueInDouble           float64 `json:"timeValueInDouble"`
	DeltaInDouble               float64 `json:"deltaInDouble"`
	GammaInDouble               float64 `json:"gammaInDouble"`
	ThetaInDouble               float64 `json:"thetaInDouble"`
	VegaInDouble                float64 `json:"vegaInDouble"`
	RhoInDouble                 float64 `json:"rhoInDouble"`
	Mark                        float64 `json:"mark"`
	Tick                        float64 `json:"tick"`
	TickAmount                  float64 `json:"tickAmount"`
	FutureIsTradable            bool    `json:"futureIsTradable"`
	FutureTradingHours          string  `json:"futureTradingHours"`
	FuturePercentChange         float64 `json:"futurePercentChange"`
	FutureIsActive              bool    `json:"futureIsActive"`
	FutureExpirationDate        float64 `json:"futureExpirationDate"`
	ExpirationType              string  `json:"expirationType"`
	ExerciseType                string  `json:"exerciseType"`
	InTheMoney                  bool    `json:"inTheMoney"`
}

type Index struct {
	Symbol          string  `json:"symbol"`
	Description     string  `json:"description"`
	LastPrice       float64 `json:"lastPrice"`
	OpenPrice       float64 `json:"openPrice"`
	HighPrice       float64 `json:"highPrice"`
	LowPrice        float64 `json:"lowPrice"`
	ClosePrice      float64 `json:"closePrice"`
	NetChange       float64 `json:"netChange"`
	TotalVolume     float64 `json:"totalVolume"`
	TradeTimeInLong float64 `json:"tradeTimeInLong"`
	Exchange        string  `json:"exchange"`
	ExchangeName    string  `json:"exchangeName"`
	Digits          float64 `json:"digits"`
	Wk52High        float64 `json:"52WkHigh"`
	Wk52Low         float64 `json:"52WkLow"`
	SecurityStatus  string  `json:"securityStatus"`
}

type OptionQ struct {
	Symbol                 string  `json:"symbol"`
	Description            string  `json:"description"`
	BidPrice               float64 `json:"bidPrice"`
	BidSize                float64 `json:"bidSize"`
	AskPrice               float64 `json:"askPrice"`
	AskSize                float64 `json:"askSize"`
	LastPrice              float64 `json:"lastPrice"`
	LastSize               float64 `json:"lastSize"`
	OpenPrice              float64 `json:"openPrice"`
	HighPrice              float64 `json:"highPrice"`
	LowPrice               float64 `json:"lowPrice"`
	ClosePrice             float64 `json:"closePrice"`
	NetChange              float64 `json:"netChange"`
	TotalVolume            float64 `json:"totalVolume"`
	QuoteTimeInLong        float64 `json:"quoteTimeInLong"`
	TradeTimeInLong        float64 `json:"tradeTimeInLong"`
	Mark                   float64 `json:"mark"`
	OpenInterest           float64 `json:"openInterest"`
	Volatility             float64 `json:"volatility"`
	MoneyIntrinsicValue    float64 `json:"moneyIntrinsicValue"`
	Multiplier             float64 `json:"multiplier"`
	StrikePrice            float64 `json:"strikePrice"`
	ContractType           string  `json:"contractType"`
	Underlying             string  `json:"underlying"`
	TimeValue              float64 `json:"timeValue"`
	Deliverables           string  `json:"deliverables"`
	Delta                  float64 `json:"delta"`
	Gamma                  float64 `json:"gamma"`
	Theta                  float64 `json:"theta"`
	Vega                   float64 `json:"vega"`
	Rho                    float64 `json:"rho"`
	SecurityStatus         string  `json:"securityStatus"`
	TheoreticalOptionValue float64 `json:"theoreticalOptionValue"`
	UnderlyingPrice        float64 `json:"underlyingPrice"`
	UvExpirationType       string  `json:"uvExpirationType"`
	Exchange               string  `json:"exchange"`
	ExchangeName           string  `json:"exchangeName"`
	SettlementType         string  `json:"settlementType"`
}

type Forex struct {
	Symbol             string  `json:"symbol"`
	BidPriceInDouble   float64 `json:"bidPriceInDouble"`
	AskPriceInDouble   float64 `json:"askPriceInDouble"`
	LastPriceInDouble  float64 `json:"lastPriceInDouble"`
	HighPriceInDouble  float64 `json:"highPriceInDouble"`
	LowPriceInDouble   float64 `json:"lowPriceInDouble"`
	ClosePriceInDouble float64 `json:"closePriceInDouble"`
	Exchange           string  `json:"exchange"`
	Description        string  `json:"description"`
	OpenPriceInDouble  float64 `json:"openPriceInDouble"`
	ChangeInDouble     float64 `json:"changeInDouble"`
	PercentChange      float64 `json:"percentChange"`
	ExchangeName       string  `json:"exchangeName"`
	Digits             float64 `json:"digits"`
	SecurityStatus     string  `json:"securityStatus"`
	Tick               float64 `json:"tick"`
	TickAmount         float64 `json:"tickAmount"`
	Product            string  `json:"product"`
	TradingHours       string  `json:"tradingHours"`
	IsTradable         bool    `json:"isTradable"`
	MarketMaker        string  `json:"marketMaker"`
	Wk52High           float64 `json:"52WkHigh"`
	Wk52Low            float64 `json:"52WkLow"`
	Mark               float64 `json:"mark"`
}

type ETF struct {
	Symbol                       string  `json:"symbol"`
	Description                  string  `json:"description"`
	BidPrice                     float64 `json:"bidPrice"`
	BidSize                      float64 `json:"bidSize"`
	BidID                        string  `json:"bidId"`
	AskPrice                     float64 `json:"askPrice"`
	AskSize                      float64 `json:"askSize"`
	AskID                        string  `json:"askId"`
	LastPrice                    float64 `json:"lastPrice"`
	LastSize                     float64 `json:"lastSize"`
	LastID                       string  `json:"lastId"`
	OpenPrice                    float64 `json:"openPrice"`
	HighPrice                    float64 `json:"highPrice"`
	LowPrice                     float64 `json:"lowPrice"`
	ClosePrice                   float64 `json:"closePrice"`
	NetChange                    float64 `json:"netChange"`
	TotalVolume                  float64 `json:"totalVolume"`
	QuoteTimeInLong              float64 `json:"quoteTimeInLong"`
	TradeTimeInLong              float64 `json:"tradeTimeInLong"`
	Mark                         float64 `json:"mark"`
	Exchange                     string  `json:"exchange"`
	ExchangeName                 string  `json:"exchangeName"`
	Marginable                   bool    `json:"marginable"`
	Shortable                    bool    `json:"shortable"`
	Volatility                   float64 `json:"volatility"`
	Digits                       float64 `json:"digits"`
	Wk52High                     float64 `json:"52WkHigh"`
	Wk52Low                      float64 `json:"52WkLow"`
	PeRatio                      float64 `json:"peRatio"`
	DivAmount                    float64 `json:"divAmount"`
	DivYield                     float64 `json:"divYield"`
	DivDate                      string  `json:"divDate"`
	SecurityStatus               string  `json:"securityStatus"`
	RegularMarketLastPrice       float64 `json:"regularMarketLastPrice"`
	RegularMarketLastSize        float64 `json:"regularMarketLastSize"`
	RegularMarketNetChange       float64 `json:"regularMarketNetChange"`
	RegularMarketTradeTimeInLong float64 `json:"regularMarketTradeTimeInLong"`
}

type EquityQ struct {
	Symbol                       string  `json:"symbol"`
	Description                  string  `json:"description"`
	BidPrice                     float64 `json:"bidPrice"`
	BidSize                      int64   `json:"bidSize"`
	BidID                        string  `json:"bidId"`
	AskPrice                     float64 `json:"askPrice"`
	AskSize                      int64   `json:"askSize"`
	AskID                        int64   `json:"askId"`
	LastPrice                    float64 `json:"lastPrice"`
	LastSize                     int64   `json:"lastSize"`
	LastID                       int64   `json:"lastId"`
	OpenPrice                    float64 `json:"openPrice"`
	HighPrice                    float64 `json:"highPrice"`
	LowPrice                     float64 `json:"lowPrice"`
	ClosePrice                   float64 `json:"closePrice"`
	NetChange                    float64 `json:"netChange"`
	TotalVolume                  float64 `json:"totalVolume"`
	QuoteTimeInLong              int64   `json:"quoteTimeInLong"`
	TradeTimeInLong              int64   `json:"tradeTimeInLong"`
	Mark                         float64 `json:"mark"`
	Exchange                     string  `json:"exchange"`
	ExchangeName                 string  `json:"exchangeName"`
	Marginable                   bool    `json:"marginable"`
	Shortable                    bool    `json:"shortable"`
	Volatility                   float64 `json:"volatility"`
	Digits                       float64 `json:"digits"`
	Wk52High                     float64 `json:"52WkHigh"`
	Wk52Low                      float64 `json:"52WkLow"`
	PeRatio                      float64 `json:"peRatio"`
	DivAmount                    float64 `json:"divAmount"`
	DivYield                     float64 `json:"divYield"`
	DivDate                      string  `json:"divDate"`
	SecurityStatus               string  `json:"securityStatus"`
	RegularMarketLastPrice       float64 `json:"regularMarketLastPrice"`
	RegularMarketLastSize        float64 `json:"regularMarketLastSize"`
	RegularMarketNetChange       float64 `json:"regularMarketNetChange"`
	RegularMarketTradeTimeInLong float64 `json:"regularMarketTradeTimeInLong"`
}
