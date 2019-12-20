package types

type Preferences struct {
	ExpressTrading                   bool    `json:"expressTrading"`
	DirectOptionsRouting             bool    `json:"directOptionsRouting"`
	DirectEquityRouting              bool    `json:"directEquityRouting"`
	DefaultEquityOrderLegInstruction string  `json:"defaultEquityOrderLegInstruction"` //"'BUY' or 'SELL' or 'BUY_TO_COVER' or 'SELL_SHORT' or 'NONE'"
	DefaultEquityOrderType           string  `json:"defaultEquityOrderType"`           //"'MARKET' or 'LIMIT' or 'STOP' or 'STOP_LIMIT' or 'TRAILING_STOP' or 'MARKET_ON_CLOSE' or 'NONE'"
	DefaultEquityOrderPriceLinkType  string  `json:"defaultEquityOrderPriceLinkType"`  //"'VALUE' or 'PERCENT' or 'NONE'"
	DefaultEquityOrderDuration       string  `json:"defaultEquityOrderDuration"`       //"'DAY' or 'GOOD_TILL_CANCEL' or 'NONE'"
	DefaultEquityOrderMarketSession  string  `json:"defaultEquityOrderMarketSession"`  //"'AM' or 'PM' or 'NORMAL' or 'SEAMLESS' or 'NONE'"
	DefaultEquityQuantity            float64 `json:"defaultEquityQuantity"`
	MutualFundTaxLotMethod           string  `json:"mutualFundTaxLotMethod"`    //"'FIFO' or 'LIFO' or 'HIGH_COST' or 'LOW_COST' or 'MINIMUM_TAX' or 'AVERAGE_COST' or 'NONE'"
	OptionTaxLotMethod               string  `json:"optionTaxLotMethod"`        //"'FIFO' or 'LIFO' or 'HIGH_COST' or 'LOW_COST' or 'MINIMUM_TAX' or 'AVERAGE_COST' or 'NONE'"
	EquityTaxLotMethod               string  `json:"equityTaxLotMethod"`        //"'FIFO' or 'LIFO' or 'HIGH_COST' or 'LOW_COST' or 'MINIMUM_TAX' or 'AVERAGE_COST' or 'NONE'"
	DefaultAdvancedToolLaunch        string  `json:"defaultAdvancedToolLaunch"` //"'TA' or 'N' or 'Y' or 'TOS' or 'NONE' or 'CC2'"
	AuthTokenTimeout                 string  `json:"authTokenTimeout"`          //"'FIFTY_FIVE_MINUTES' or 'TWO_HOURS' or 'FOUR_HOURS' or 'EIGHT_HOURS'"
}

type Key struct {
	Key string `json:"key"` //"string"
}

type SubscriptionKey struct {
	Keys []*Key `json:"keys"`
}

type StreamerInfo struct {
	StreamerBinaryURL string `json:"streamerBinaryUrl"` //"string"
	StreamerSocketURL string `json:"streamerSocketUrl"` //"string"
	Token             string `json:"token"`             //"string"
	TokenTimestamp    string `json:"tokenTimestamp"`    //"string"
	UserGroup         string `json:"userGroup"`         //"string"
	AccessLevel       string `json:"accessLevel"`       //"string"
	ACL               string `json:"acl"`               //"string"
	AppID             string `json:"appId"`             //"string"
}

type QuotesU struct {
	IsNyseDelayed   bool `json:"isNyseDelayed"`
	IsNasdaqDelayed bool `json:"isNasdaqDelayed"`
	IsOpraDelayed   bool `json:"isOpraDelayed"`
	IsAmexDelayed   bool `json:"isAmexDelayed"`
	IsCmeDelayed    bool `json:"isCmeDelayed"`
	IsIceDelayed    bool `json:"isIceDelayed"`
	IsForexDelayed  bool `json:"isForexDelayed"`
}

type Authorizations struct {
	Apex               bool   `json:"apex"`
	LevelTwoQuotes     bool   `json:"levelTwoQuotes"`
	StockTrading       bool   `json:"stockTrading"`
	MarginTrading      bool   `json:"marginTrading"`
	StreamingNews      bool   `json:"streamingNews"`
	OptionTradingLevel string `json:"optionTradingLevel"` //"'COVERED' or 'FULL' or 'LONG' or 'SPREAD' or 'NONE'"
	StreamerAccess     bool   `json:"streamerAccess"`
	AdvancedMargin     bool   `json:"advancedMargin"`
	ScottradeAccount   bool   `json:"scottradeAccount"`
}

type AccountU struct {
	AccountID         int64             `json:"accountId,string"`  //"string"
	Description       string            `json:"description"`       //"string"
	DisplayName       string            `json:"displayName"`       //"string"
	AccountCdDomainID string            `json:"accountCdDomainId"` //"string"
	Company           string            `json:"company"`           //"string"
	Segment           string            `json:"segment"`           //"string"
	SurrogateIds      map[string]string `json:"surrogateIds"`      //"object"
	Preferences       Preferences       `json:"preferences"`
	ACL               string            `json:"acl"` //"string"
	Authorizations    Authorizations    `json:"authorizations"`
}

type UserPrincipal struct {
	AuthToken                string           `json:"authToken"`               //"string"
	UserID                   string           `json:"userId"`                  //"string"
	UserCdDomainID           string           `json:"userCdDomainId"`          //"string"
	PrimaryAccountID         int64            `json:"primaryAccountId,string"` //"string"
	LastLoginTime            string           `json:"lastLoginTime"`           //"string"
	TokenExpirationTime      string           `json:"tokenExpirationTime"`     //"string"
	LoginTime                string           `json:"loginTime"`               //"string"
	AccessLevel              string           `json:"accessLevel"`             //"string"
	StalePassword            bool             `json:"stalePassword"`
	StreamerInfo             *StreamerInfo    `json:"streamerInfo"`
	ProfessionalStatus       string           `json:"professionalStatus"` //"'PROFESSIONAL' or 'NON_PROFESSIONAL' or 'UNKNOWN_STATUS'"
	Quotes                   *QuotesU         `json:"quotes"`
	StreamerSubscriptionKeys *SubscriptionKey `json:"streamerSubscriptionKeys"`
	Accounts                 []*AccountU      `json:"accounts"`
}
