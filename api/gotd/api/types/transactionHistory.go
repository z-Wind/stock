package types

type InstrumentT struct {
	Symbol               string  `json:"symbol"`               //"string"
	UnderlyingSymbol     string  `json:"underlyingSymbol"`     //"string"
	OptionExpirationDate string  `json:"optionExpirationDate"` //"string"
	OptionStrikePrice    float64 `json:"optionStrikePrice"`
	PutCall              string  `json:"putCall"`          //"string"
	Cusip                string  `json:"cusip"`            //"string"
	Description          string  `json:"description"`      //"string"
	AssetType            string  `json:"assetType"`        //"string"
	BondMaturityDate     string  `json:"bondMaturityDate"` //"string"
	BondInterestRate     float64 `json:"bondInterestRate"`
}

type TransactionItem struct {
	AccountID            int64        `json:"accountId"`
	Amount               float64      `json:"amount"`
	Price                float64      `json:"price"`
	Cost                 float64      `json:"cost"`
	ParentOrderKey       float64      `json:"parentOrderKey"`
	ParentChildIndicator string       `json:"parentChildIndicator"` //"string"
	Instruction          string       `json:"instruction"`          //"string"
	PositionEffect       string       `json:"positionEffect"`       //"string"
	Instrument           *InstrumentT `json:"instrument"`
}

type Fees struct {
	RFee          float64 `json:"rFee"`
	AdditionalFee float64 `json:"additionalFee"`
	CdscFee       float64 `json:"cdscFee"`
	RegFee        float64 `json:"regFee"`
	OtherCharges  float64 `json:"otherCharges"`
	Commission    float64 `json:"commission"`
	OptRegFee     float64 `json:"optRegFee"`
	SecFee        float64 `json:"secFee"`
}

type Transactions []*Transaction

// func (t *Transactions) Index(i int) *Transaction {
// 	return []*Transaction(*t)[i]
// }

// func (t *Transactions) Len() int {
// 	return len([]*Transaction(*t))
// }

// func (t *Transactions) Iterate() []*Transaction {
// 	return []*Transaction(*t)
// }

type Transaction struct {
	Type                          string           `json:"type"`                    //"'TRADE' or 'RECEIVE_AND_DELIVER' or 'DIVIDEND_OR_INTEREST' or 'ACH_RECEIPT' or 'ACH_DISBURSEMENT' or 'CASH_RECEIPT' or 'CASH_DISBURSEMENT' or 'ELECTRONIC_FUND' or 'WIRE_OUT' or 'WIRE_IN' or 'JOURNAL' or 'MEMORANDUM' or 'MARGIN_CALL' or 'MONEY_MARKET' or 'SMA_ADJUSTMENT'"
	ClearingReferenceNumber       string           `json:"clearingReferenceNumber"` //"string"
	SubAccount                    string           `json:"subAccount"`              //"string"
	SettlementDate                string           `json:"settlementDate"`          //"string"
	OrderID                       int64            `json:"orderId"`                 //"string"
	Sma                           float64          `json:"sma"`
	RequirementReallocationAmount float64          `json:"requirementReallocationAmount"`
	DayTradeBuyingPowerEffect     float64          `json:"dayTradeBuyingPowerEffect"`
	NetAmount                     float64          `json:"netAmount"`
	TransactionDate               string           `json:"transactionDate"`    //"string"
	OrderDate                     string           `json:"orderDate"`          //"string"
	TransactionSubType            string           `json:"transactionSubType"` //"string"
	TransactionID                 int64            `json:"transactionId"`
	CashBalanceEffectFlag         bool             `json:"cashBalanceEffectFlag"`
	Description                   string           `json:"description"` //"string"
	AchStatus                     string           `json:"achStatus"`   //"'Approved' or 'Rejected' or 'Cancel' or 'Error'"
	AccruedInterest               float64          `json:"accruedInterest"`
	Fees                          *Fees            `json:"fees"` //"object"
	TransactionItem               *TransactionItem `json:"transactionItem"`
}
