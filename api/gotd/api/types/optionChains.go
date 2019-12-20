package types

type Underlying struct {
	Ask               float64 `json:"ask"`
	AskSize           int64   `json:"askSize"`
	Bid               float64 `json:"bid"`
	BidSize           int64   `json:"bidSize"`
	Change            float64 `json:"change"`
	Close             float64 `json:"close"`
	Delayed           bool    `json:"delayed"`
	Description       string  `json:"description"`  //"string"
	ExchangeName      string  `json:"exchangeName"` //"string"
	FiftyTwoWeekHigh  float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow   float64 `json:"fiftyTwoWeekLow"`
	HighPrice         float64 `json:"highPrice"`
	Last              float64 `json:"last"`
	LowPrice          float64 `json:"lowPrice"`
	Mark              float64 `json:"mark"`
	MarkChange        float64 `json:"markChange"`
	MarkPercentChange float64 `json:"markPercentChange"`
	OpenPrice         float64 `json:"openPrice"`
	PercentChange     float64 `json:"percentChange"`
	QuoteTime         int64   `json:"quoteTime"`
	Symbol            string  `json:"symbol"` //"string"
	TotalVolume       int64   `json:"totalVolume"`
	TradeTime         int64   `json:"tradeTime"`
}

type Leg struct {
	Symbol      string  `json:"symbol"`
	PutCallInd  string  `json:"putCallInd"`
	Description string  `json:"description"`
	Bid         float64 `json:"bid"`
	Ask         float64 `json:"ask"`
	Range       string  `json:"range"`
	StrikePrice float64 `json:"strikePrice"`
	TotalVolume float64 `json:"totalVolume"`
}

type OptionStrategy struct {
	PrimaryLeg     *Leg    `json:"primaryLeg"`
	SecondaryLeg   *Leg    `json:"secondaryLeg"`
	StrategyStrike string  `json:"strategyStrike"`
	StrategyBid    float64 `json:"strategyBid"`
	StrategyAsk    float64 `json:"strategyAsk"`
}

type MonthlyStrategy struct {
	Month              string            `json:"month"`
	Year               int64             `json:"year"`
	Day                int64             `json:"day"`
	DaysToExp          int64             `json:"daysToExp"`
	SecondaryMonth     *string           `json:"secondaryMonth"`
	SecondaryYear      *int64            `json:"secondaryYear"`
	SecondaryDay       *int64            `json:"secondaryDay"`
	SecondaryDaysToExp *int64            `json:"secondaryDaysToExp"`
	Type               string            `json:"type"`
	SecondaryType      string            `json:"secondaryType"`
	SecondaryLeap      bool              `json:"secondaryLeap"`
	OptionStrategyList []*OptionStrategy `json:"optionStrategyList"`
}

type OptionChain struct {
	Symbol              string                          `json:"symbol"` //"string"
	Status              string                          `json:"status"` //"string"
	Underlying          *Underlying                     `json:"underlying"`
	Strategy            string                          `json:"strategy"` //"'SINGLE' or 'ANALYTICAL' or 'COVERED' or 'VERTICAL' or 'CALENDAR' or 'STRANGLE' or 'STRADDLE' or 'BUTTERFLY' or 'CONDOR' or 'DIAGONAL' or 'COLLAR' or 'ROLL'"
	Interval            *float64                        `json:"interval"`
	Intervals           []float64                       `json:"intervals"`
	IsDelayed           bool                            `json:"isDelayed"`
	IsIndex             bool                            `json:"isIndex"`
	DaysToExpiration    float64                         `json:"daysToExpiration"`
	InterestRate        float64                         `json:"interestRate"`
	UnderlyingPrice     float64                         `json:"underlyingPrice"`
	Volatility          float64                         `json:"volatility"`
	MonthlyStrategyList []*MonthlyStrategy              `json:"monthlyStrategyList"`
	CallExpDateMap      map[string]map[string][]*Option `json:"callExpDateMap"` //"object"
	PutExpDateMap       map[string]map[string][]*Option `json:"putExpDateMap"`  //"object"
}

type OptionDeliverables struct {
	Symbol           string `json:"symbol"`           //"string"
	AssetType        string `json:"assetType"`        //"string"
	DeliverableUnits string `json:"deliverableUnits"` //"string"
	CurrencyType     string `json:"currencyType"`     //"string"
}

type Option struct {
	PutCall                string                `json:"putCall"`      //"'PUT' or 'CALL'"
	Symbol                 string                `json:"symbol"`       //"string"
	Description            string                `json:"description"`  //"string"
	ExchangeName           string                `json:"exchangeName"` //"string"
	Bid                    float64               `json:"bid"`
	Ask                    float64               `json:"ask"`
	Last                   float64               `json:"last"`
	Mark                   float64               `json:"mark"`
	BidSize                float64               `json:"bidSize"`
	AskSize                float64               `json:"askSize"`
	LastSize               float64               `json:"lastSize"`
	HighPrice              float64               `json:"highPrice"`
	LowPrice               float64               `json:"lowPrice"`
	OpenPrice              float64               `json:"openPrice"`
	ClosePrice             float64               `json:"closePrice"`
	TotalVolume            float64               `json:"totalVolume"`
	TradeDate              *int64                `json:"tradeDate"`
	QuoteTimeInLong        float64               `json:"quoteTimeInLong"`
	TradeTimeInLong        float64               `json:"tradeTimeInLong"`
	NetChange              float64               `json:"netChange"`
	Volatility             float64               `json:"volatility"`
	Delta                  float64               `json:"delta"`
	Gamma                  float64               `json:"gamma"`
	Theta                  float64               `json:"theta"`
	Vega                   float64               `json:"vega"`
	Rho                    float64               `json:"rho"`
	TimeValue              float64               `json:"timeValue"`
	OpenInterest           float64               `json:"openInterest"`
	IsInTheMoney           bool                  `json:"isInTheMoney"`
	TheoreticalOptionValue float64               `json:"theoreticalOptionValue"`
	TheoreticalVolatility  float64               `json:"theoreticalVolatility"`
	IsMini                 bool                  `json:"isMini"`
	IsNonStandard          bool                  `json:"isNonStandard"`
	OptionDeliverablesList []*OptionDeliverables `json:"optionDeliverablesList"`
	StrikePrice            float64               `json:"strikePrice"`
	ExpirationDate         int64                 `json:"expirationDate"`
	DaysToExpiration       int64                 `json:"daysToExpiration"`
	ExpirationType         string                `json:"expirationType"` //"string"
	LastTradingDay         int64                 `json:"lastTradingDay"`
	Multiplier             float64               `json:"multiplier"`
	SettlementType         string                `json:"settlementType"`  //"string"
	DeliverableNote        string                `json:"deliverableNote"` //"string"
	IsIndexOption          *bool                 `json:"isIndexOption"`
	PercentChange          float64               `json:"percentChange"`
	MarkChange             float64               `json:"markChange"`
	MarkPercentChange      float64               `json:"markPercentChange"`
	NonStandard            bool                  `json:"nonStandard"`
	Mini                   bool                  `json:"mini"`
	InTheMoney             bool                  `json:"inTheMoney"`
}
