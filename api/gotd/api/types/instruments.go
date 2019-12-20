package types

type Instruments map[string]*Instrument

// func (i *Instruments) Get(key string) *Instrument {
// 	return map[string]*Instrument(*i)[key]
// }

// func (i *Instruments) Len() int {
// 	return len(map[string]*Instrument(*i))
// }

// func (i *Instruments) Iterate() map[string]*Instrument {
// 	return map[string]*Instrument(*i)
// }

type Instrument struct {
	Cusip       string `json:"cusip"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	Exchange    string `json:"exchange"`
	AssetType   string `json:"assetType"` //"'EQUITY' or 'ETF' or 'FOREX' or 'FUTURE' or 'FUTURE_OPTION' or 'INDEX' or 'INDICATOR' or 'MUTUAL_FUND' or 'OPTION' or 'UNKNOWN' or 'BOND'"

	Fundamental *FundamentalData `json:"fundamental,omitempty"`
	BondPrice   float64          `json:"bondPrice,omitempty"` // for "'BOND'"
}

type FundamentalData struct {
	Symbol              string  `json:"symbol"`
	High52              float64 `json:"high52"`
	Low52               float64 `json:"low52"`
	DividendAmount      float64 `json:"dividendAmount"`
	DividendYield       float64 `json:"dividendYield"`
	DividendDate        string  `json:"dividendDate"`
	PeRatio             float64 `json:"peRatio"`
	PegRatio            float64 `json:"pegRatio"`
	PbRatio             float64 `json:"pbRatio"`
	PrRatio             float64 `json:"prRatio"`
	PcfRatio            float64 `json:"pcfRatio"`
	GrossMarginTTM      float64 `json:"grossMarginTTM"`
	GrossMarginMRQ      float64 `json:"grossMarginMRQ"`
	NetProfitMarginTTM  float64 `json:"netProfitMarginTTM"`
	NetProfitMarginMRQ  float64 `json:"netProfitMarginMRQ"`
	OperatingMarginTTM  float64 `json:"operatingMarginTTM"`
	OperatingMarginMRQ  float64 `json:"operatingMarginMRQ"`
	ReturnOnEquity      float64 `json:"returnOnEquity"`
	ReturnOnAssets      float64 `json:"returnOnAssets"`
	ReturnOnInvestment  float64 `json:"returnOnInvestment"`
	QuickRatio          float64 `json:"quickRatio"`
	CurrentRatio        float64 `json:"currentRatio"`
	InterestCoverage    float64 `json:"interestCoverage"`
	TotalDebtToCapital  float64 `json:"totalDebtToCapital"`
	LtDebtToEquity      float64 `json:"ltDebtToEquity"`
	TotalDebtToEquity   float64 `json:"totalDebtToEquity"`
	EpsTTM              float64 `json:"epsTTM"`
	EpsChangePercentTTM float64 `json:"epsChangePercentTTM"`
	EpsChangeYear       float64 `json:"epsChangeYear"`
	EpsChange           float64 `json:"epsChange"`
	RevChangeYear       float64 `json:"revChangeYear"`
	RevChangeTTM        float64 `json:"revChangeTTM"`
	RevChangeIn         float64 `json:"revChangeIn"`
	SharesOutstanding   float64 `json:"sharesOutstanding"`
	MarketCapFloat      float64 `json:"marketCapFloat"`
	MarketCap           float64 `json:"marketCap"`
	BookValuePerShare   float64 `json:"bookValuePerShare"`
	ShortIntToFloat     float64 `json:"shortIntToFloat"`
	ShortIntDayToCover  float64 `json:"shortIntDayToCover"`
	DivGrowthRate3Year  float64 `json:"divGrowthRate3Year"`
	DividendPayAmount   float64 `json:"dividendPayAmount"`
	DividendPayDate     string  `json:"dividendPayDate"`
	Beta                float64 `json:"beta"`
	Vol1DayAvg          float64 `json:"vol1DayAvg"`
	Vol10DayAvg         float64 `json:"vol10DayAvg"`
	Vol3MonthAvg        float64 `json:"vol3MonthAvg"`
}
