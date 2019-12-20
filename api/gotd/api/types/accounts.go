package types

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type Accounts []*Account

// func (a *Accounts) Index(i int) *Account {
// 	return []*Account(*a)[i]
// }

// func (a *Accounts) Len() int {
// 	return len([]*Account(*a))
// }

// func (a *Accounts) Iterate() []*Account {
// 	return []*Account(*a)
// }

type Account struct {
	SecuritiesAccount `json:"securitiesAccount"`
}

type InstrumentA struct {
	AssetType string `json:"assetType"`

	// EQUITY
	// OPTION
	// MUTUAL_FUND
	// CASH_EQUIVALENT
	// FIXED_INCOME
	// CURRENCY 未定義
	// INDEX 未定義
	Data interface{}
}

type _InstrumentA InstrumentA

func (i *InstrumentA) UnmarshalJSON(bs []byte) (err error) {
	instrument := _InstrumentA{}

	err = json.Unmarshal(bs, &instrument)
	if err != nil {
		return errors.Wrap(err, "json.Unmarshal")
	}

	switch instrument.AssetType {
	case "EQUITY":
		instrument.Data = &Equity{}
	case "OPTION":
		instrument.Data = &OptionA{}
	case "MUTUAL_FUND":
		instrument.Data = &MutualFund{}
	case "CASH_EQUIVALENT":
		instrument.Data = &CashEquivalent{}
	case "FIXED_INCOME":
		instrument.Data = &FixedIncome{}
	//case "CURRENCY":
	//	instrument.Data = &{}
	//case "INDEX":
	//  instrument.Data = &Index{}
	default:
		return fmt.Errorf("Not support type %s", instrument.AssetType)
	}
	err = json.Unmarshal(bs, instrument.Data)
	*i = InstrumentA(instrument)

	return err
}

func (i *InstrumentA) MarshalJSON() ([]byte, error) {
	switch data := i.Data.(type) {
	case *Equity:
		return json.Marshal(&struct {
			AssetType string `json:"assetType"`
			*Equity
		}{
			AssetType: i.AssetType,
			Equity:    data,
		})
	case *OptionA:
		return json.Marshal(&struct {
			AssetType string `json:"assetType"`
			*OptionA
		}{
			AssetType: i.AssetType,
			OptionA:   data,
		})
	case *MutualFund:
		return json.Marshal(&struct {
			AssetType string `json:"assetType"`
			*MutualFund
		}{
			AssetType:  i.AssetType,
			MutualFund: data,
		})
	case *CashEquivalent:
		return json.Marshal(&struct {
			AssetType string `json:"assetType"`
			*CashEquivalent
		}{
			AssetType:      i.AssetType,
			CashEquivalent: data,
		})
	case *FixedIncome:
		return json.Marshal(&struct {
			AssetType string `json:"assetType"`
			*FixedIncome
		}{
			AssetType:   i.AssetType,
			FixedIncome: data,
		})
	//case "CURRENCY":
	//	instrument.Data = &{}
	//case "INDEX":
	//  instrument.Data = &Index{}
	default:
		return nil, fmt.Errorf("unexpected type %T: %v", data, data)
	}

}

type Position struct {
	ShortQuantity                  float64      `json:"shortQuantity"`
	AveragePrice                   float64      `json:"averagePrice"`
	CurrentDayProfitLoss           float64      `json:"currentDayProfitLoss"`
	CurrentDayProfitLossPercentage float64      `json:"currentDayProfitLossPercentage"`
	LongQuantity                   float64      `json:"longQuantity"`
	SettledLongQuantity            float64      `json:"settledLongQuantity"`
	SettledShortQuantity           float64      `json:"settledShortQuantity"`
	AgedQuantity                   float64      `json:"agedQuantity"`
	InstrumenA                     *InstrumentA `json:"instrument"`
	MarketValue                    float64      `json:"marketValue"`
}

type OrderLegCollection struct {
	OrderLegType   string       `json:"orderLegType"` //"'EQUITY' or 'OPTION' or 'INDEX' or 'MUTUAL_FUND' or 'CASH_EQUIVALENT' or 'FIXED_INCOME' or 'CURRENCY'",
	LegID          int64        `json:"legId,omitempty"`
	InstrumentA    *InstrumentA `json:"instrument"`               //"\"The type <Instrument> has the following subclasses [Option, MutualFund, CashEquivalent, Equity, FixedIncome] descriptions are listed below\"",
	Instruction    string       `json:"instruction"`              //"'BUY' or 'SELL' or 'BUY_TO_COVER' or 'SELL_SHORT' or 'BUY_TO_OPEN' or 'BUY_TO_CLOSE' or 'SELL_TO_OPEN' or 'SELL_TO_CLOSE' or 'EXCHANGE'",
	PositionEffect string       `json:"positionEffect,omitempty"` //"'OPENING' or 'CLOSING' or 'AUTOMATIC'",
	Quantity       float64      `json:"quantity"`
	QuantityType   string       `json:"quantityType,omitempty"` //"'ALL_SHARES' or 'DOLLARS' or 'SHARES'"
}

type OrderStrategie struct {
	Session                  string                `json:"session"`   //"'NORMAL' or 'AM' or 'PM' or 'SEAMLESS'",
	Duration                 string                `json:"duration"`  //"'DAY' or 'GOOD_TILL_CANCEL' or 'FILL_OR_KILL'",
	OrderType                string                `json:"orderType"` //"'MARKET' or 'LIMIT' or 'STOP' or 'STOP_LIMIT' or 'TRAILING_STOP' or 'MARKET_ON_CLOSE' or 'EXERCISE' or 'TRAILING_STOP_LIMIT' or 'NET_DEBIT' or 'NET_CREDIT' or 'NET_ZERO'",
	CancelTime               string                `json:"cancelTime"`
	ComplexOrderStrategyType string                `json:"complexOrderStrategyType"` //"'NONE' or 'COVERED' or 'VERTICAL' or 'BACK_RATIO' or 'CALENDAR' or 'DIAGONAL' or 'STRADDLE' or 'STRANGLE' or 'COLLAR_SYNTHETIC' or 'BUTTERFLY' or 'CONDOR' or 'IRON_CONDOR' or 'VERTICAL_ROLL' or 'COLLAR_WITH_STOCK' or 'DOUBLE_DIAGONAL' or 'UNBALANCED_BUTTERFLY' or 'UNBALANCED_CONDOR' or 'UNBALANCED_IRON_CONDOR' or 'UNBALANCED_VERTICAL_ROLL' or 'CUSTOM'",
	Quantity                 float64               `json:"quantity"`
	FilledQuantity           float64               `json:"filledQuantity"`
	RemainingQuantity        float64               `json:"remainingQuantity"`
	RequestedDestination     string                `json:"requestedDestination"` //"'INET' or 'ECN_ARCA' or 'CBOE' or 'AMEX' or 'PHLX' or 'ISE' or 'BOX' or 'NYSE' or 'NASDAQ' or 'BATS' or 'C2' or 'AUTO'",
	DestinationLinkName      string                `json:"destinationLinkName"`
	ReleaseTime              string                `json:"releaseTime"`
	StopPrice                float64               `json:"stopPrice"`
	StopPriceLinkBasis       string                `json:"stopPriceLinkBasis"` //"'MANUAL' or 'BASE' or 'TRIGGER' or 'LAST' or 'BID' or 'ASK' or 'ASK_BID' or 'MARK' or 'AVERAGE'",
	StopPriceLinkType        string                `json:"stopPriceLinkType"`  //"'VALUE' or 'PERCENT' or 'TICK'",
	StopPriceOffset          float64               `json:"stopPriceOffset"`
	StopType                 string                `json:"stopType"`       //"'STANDARD' or 'BID' or 'ASK' or 'LAST' or 'MARK'",
	PriceLinkBasis           string                `json:"priceLinkBasis"` //"'MANUAL' or 'BASE' or 'TRIGGER' or 'LAST' or 'BID' or 'ASK' or 'ASK_BID' or 'MARK' or 'AVERAGE'",
	PriceLinkType            string                `json:"priceLinkType"`  //"'VALUE' or 'PERCENT' or 'TICK'",
	Price                    float64               `json:"price"`
	TaxLotMethod             string                `json:"taxLotMethod"` //"'FIFO' or 'LIFO' or 'HIGH_COST' or 'LOW_COST' or 'AVERAGE_COST' or 'SPECIFIC_LOT'",
	OrderLegCollections      []*OrderLegCollection `json:"orderLegCollection"`
	ActivationPrice          float64               `json:"activationPrice"`
	SpecialInstruction       string                `json:"specialInstruction"` //"'ALL_OR_NONE' or 'DO_NOT_REDUCE' or 'ALL_OR_NONE_DO_NOT_REDUCE'",
	OrderStrategyType        string                `json:"orderStrategyType"`  //"'SINGLE' or 'OCO' or 'TRIGGER'",
	OrderID                  int64                 `json:"orderId"`
	Cancelable               bool                  `json:"cancelable"`
	Editable                 bool                  `json:"editable"`
	Status                   string                `json:"status"` //"'AWAITING_PARENT_ORDER' or 'AWAITING_CONDITION' or 'AWAITING_MANUAL_REVIEW' or 'ACCEPTED' or 'AWAITING_UR_OUT' or 'PENDING_ACTIVATION' or 'QUEUED' or 'WORKING' or 'REJECTED' or 'PENDING_CANCEL' or 'CANCELED' or 'PENDING_REPLACE' or 'REPLACED' or 'FILLED' or 'EXPIRED'",
	EnteredTime              string                `json:"enteredTime"`
	CloseTime                string                `json:"closeTime"`
	Tag                      string                `json:"tag"`
	AccountID                int64                 `json:"accountId"`
	OrderActivityCollection  []*Execution          `json:"orderActivityCollection"`  //: ["\"The type <OrderActivity> has the following subclasses [Execution] descriptions are listed below\""],
	ReplacingOrderCollection []json.RawMessage     `json:"replacingOrderCollection"` //: [ {} ],
	ChildOrderStrategies     []json.RawMessage     `json:"childOrderStrategies"`     //: [ {}  ],
	StatusDescription        string                `json:"statusDescription"`
}

type InitialBalances struct {
	AccruedInterest            float64 `json:"accruedInterest"`
	CashAvailableForTrading    float64 `json:"cashAvailableForTrading"`
	CashAvailableForWithdrawal float64 `json:"cashAvailableForWithdrawal"`
	CashBalance                float64 `json:"cashBalance"`
	BondValue                  float64 `json:"bondValue"`
	CashReceipts               float64 `json:"cashReceipts"`
	LiquidationValue           float64 `json:"liquidationValue"`
	LongOptionMarketValue      float64 `json:"longOptionMarketValue"`
	LongStockValue             float64 `json:"longStockValue"`
	MoneyMarketFund            float64 `json:"moneyMarketFund"`
	MutualFundValue            float64 `json:"mutualFundValue"`
	ShortOptionMarketValue     float64 `json:"shortOptionMarketValue"`
	ShortStockValue            float64 `json:"shortStockValue"`
	IsInCall                   bool    `json:"isInCall"`
	UnsettledCash              float64 `json:"unsettledCash"`
	CashDebitCallValue         float64 `json:"cashDebitCallValue"`
	PendingDeposits            float64 `json:"pendingDeposits"`
	AccountValue               float64 `json:"accountValue"`
}

type CurrentBalances struct {
	AccruedInterest              float64 `json:"accruedInterest"`
	CashBalance                  float64 `json:"cashBalance"`
	CashReceipts                 float64 `json:"cashReceipts"`
	LongOptionMarketValue        float64 `json:"longOptionMarketValue"`
	LiquidationValue             float64 `json:"liquidationValue"`
	LongMarketValue              float64 `json:"longMarketValue"`
	MoneyMarketFund              float64 `json:"moneyMarketFund"`
	Savings                      float64 `json:"savings"`
	ShortMarketValue             float64 `json:"shortMarketValue"`
	PendingDeposits              float64 `json:"pendingDeposits"`
	CashAvailableForTrading      float64 `json:"cashAvailableForTrading"`
	CashAvailableForWithdrawal   float64 `json:"cashAvailableForWithdrawal"`
	CashCall                     float64 `json:"cashCall"`
	LongNonMarginableMarketValue float64 `json:"longNonMarginableMarketValue"`
	TotalCash                    float64 `json:"totalCash"`
	ShortOptionMarketValue       float64 `json:"shortOptionMarketValue"`
	MutualFundValue              float64 `json:"mutualFundValue"`
	BondValue                    float64 `json:"bondValue"`
	CashDebitCallValue           float64 `json:"cashDebitCallValue"`
	UnsettledCash                float64 `json:"unsettledCash"`
}
type ProjectedBalances struct {
	AccruedInterest              float64 `json:"accruedInterest"`
	CashBalance                  float64 `json:"cashBalance"`
	CashReceipts                 float64 `json:"cashReceipts"`
	LongOptionMarketValue        float64 `json:"longOptionMarketValue"`
	LiquidationValue             float64 `json:"liquidationValue"`
	LongMarketValue              float64 `json:"longMarketValue"`
	MoneyMarketFund              float64 `json:"moneyMarketFund"`
	Savings                      float64 `json:"savings"`
	ShortMarketValue             float64 `json:"shortMarketValue"`
	PendingDeposits              float64 `json:"pendingDeposits"`
	CashAvailableForTrading      float64 `json:"cashAvailableForTrading"`
	CashAvailableForWithdrawal   float64 `json:"cashAvailableForWithdrawal"`
	CashCall                     float64 `json:"cashCall"`
	LongNonMarginableMarketValue float64 `json:"longNonMarginableMarketValue"`
	TotalCash                    float64 `json:"totalCash"`
	ShortOptionMarketValue       float64 `json:"shortOptionMarketValue"`
	MutualFundValue              float64 `json:"mutualFundValue"`
	BondValue                    float64 `json:"bondValue"`
	CashDebitCallValue           float64 `json:"cashDebitCallValue"`
	UnsettledCash                float64 `json:"unsettledCash"`
}

type SecuritiesAccount struct {
	Type                    string             `json:"type"`
	AccountID               int64              `json:"accountId,string"`
	RoundTrips              float64            `json:"roundTrips"`
	IsDayTrader             bool               `json:"isDayTrader"`
	IsClosingOnlyRestricted bool               `json:"isClosingOnlyRestricted"`
	Positions               []*Position        `json:"positions"`
	OrderStrategies         []*OrderStrategie  `json:"orderStrategies"`
	InitialBalances         *InitialBalances   `json:"initialBalances"`
	CurrentBalances         *CurrentBalances   `json:"currentBalances"`
	ProjectedBalances       *ProjectedBalances `json:"projectedBalances"`
}

type OptionDeliverable struct {
	Symbol           string  `json:"symbol"`
	DeliverableUnits float64 `json:"deliverableUnits"`
	CurrencyType     string  `json:"currencyType"` //"'USD' or 'CAD' or 'EUR' or 'JPY'",
	AssetType        string  `json:"assetType"`    //"'EQUITY' or 'OPTION' or 'INDEX' or 'MUTUAL_FUND' or 'CASH_EQUIVALENT' or 'FIXED_INCOME' or 'CURRENCY'"
}

type OptionA struct {
	Cusip              string               `json:"cusip,omitempty"`
	Symbol             string               `json:"symbol"`
	Description        string               `json:"description,omitempty"`
	Type               string               `json:"type"`    //"'VANILLA' or 'BINARY' or 'BARRIER'",
	PutCall            string               `json:"putCall"` //"'PUT' or 'CALL'",
	UnderlyingSymbol   string               `json:"underlyingSymbol"`
	OptionMultiplier   float64              `json:"optionMultiplier"`
	OptionDeliverables []*OptionDeliverable `json:"optionDeliverables"`
}

type MutualFund struct {
	Cusip       string `json:"cusip,omitempty"`
	Symbol      string `json:"symbol"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"` //"'NOT_APPLICABLE' or 'OPEN_END_NON_TAXABLE' or 'OPEN_END_TAXABLE' or 'NO_LOAD_NON_TAXABLE' or 'NO_LOAD_TAXABLE'"
}

type CashEquivalent struct {
	Cusip       string `json:"cusip,omitempty"`
	Symbol      string `json:"symbol"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"` //"'SAVINGS' or 'MONEY_MARKET_FUND'"
}

type Equity struct {
	Cusip       string `json:"cusip,omitempty"`
	Symbol      string `json:"symbol"`
	Description string `json:"description,omitempty"`
}

type FixedIncome struct {
	Cusip        string  `json:"cusip"`
	Symbol       string  `json:"symbol"`
	Description  string  `json:"description"`
	MaturityDate string  `json:"maturityDate"`
	VariableRate float64 `json:"variableRate"`
	Factor       float64 `json:"factor"`
}

//The class <OrderActivity> has the
//following subclasses:
//-Execution
//JSON for each are listed below:

type ExecutionLeg struct {
	LegID             int64   `json:"legId"`
	Quantity          float64 `json:"quantity"`
	MismarkedQuantity float64 `json:"mismarkedQuantity"`
	Price             float64 `json:"price"`
	Time              string  `json:"time"`
}

type Execution struct {
	ActivityType           string          `json:"activityType"`  //"'EXECUTION' or 'ORDER_ACTION'",
	ExecutionType          string          `json:"executionType"` //"'FILL'",
	Quantity               float64         `json:"quantity"`
	OrderRemainingQuantity float64         `json:"orderRemainingQuantity"`
	ExecutionLegs          []*ExecutionLeg `json:"executionLegs"`
}
