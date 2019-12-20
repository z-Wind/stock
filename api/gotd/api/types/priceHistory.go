package types

type Candle struct {
	Close    float64 `json:"close"`
	Datetime int64   `json:"datetime"`
	High     float64 `json:"high"`
	Low      float64 `json:"low"`
	Open     float64 `json:"open"`
	Volume   float64 `json:"volume"`
}

type CandleList struct {
	Candles []Candle `json:"candles"`
	Empty   bool     `json:"empty"`
	Symbol  string   `json:"symbol"` //"string"
}
