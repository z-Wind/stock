package api

import (
	"time"
)

const (
	AccountFieldPositions = "positions"
	AccountFieldOrders    = "orders"

	InstrumentsProjectionSymbolSearch = "symbol-search"
	InstrumentsProjectionSymbolRegex  = "symbol-regex"
	InstrumentsProjectionDescSearch   = "desc-search"
	InstrumentsProjectionDescRegex    = "desc-regex"
	InstrumentsProjectionFundamental  = "fundamental"

	HoursMarketsEQUITY = "EQUITY"
	HoursMarketsOPTION = "OPTION"
	HoursMarketsFUTURE = "FUTURE"
	HoursMarketsBOND   = "BOND"
	HoursMarketsFOREX  = "FOREX"

	MoversIndexCOMPX    = "$COMPX"
	MoversIndexDJI      = "$DJI"
	MoversIndexSPXX     = "$SPX.X"
	MoversDirectionUp   = "up"
	MoversDirectionDown = "down"
	MoversChangePercent = "percent"
	MoversChangeValue   = "value"

	OptionChainContractTypeCALL   = "CALL"
	OptionChainContractTypePUT    = "PUT"
	OptionChainContractTypeALL    = "ALL"
	OptionChainStrategySINGLE     = "SINGLE"
	OptionChainStrategyANALYTICAL = "ANALYTICAL"
	OptionChainStrategyCOVERED    = "COVERED"
	OptionChainStrategyVERTICAL   = "VERTICAL"
	OptionChainStrategyCALENDAR   = "CALENDAR"
	OptionChainStrategySTRANGLE   = "STRANGLE"
	OptionChainStrategySTRADDLE   = "STRADDLE"
	OptionChainStrategyBUTTERFLY  = "BUTTERFLY"
	OptionChainStrategyCONDOR     = "CONDOR"
	OptionChainStrategyDIAGONAL   = "DIAGONAL"
	OptionChainStrategyCOLLAR     = "COLLAR"
	OptionChainStrategyROLL       = "ROLL"
	OptionChainRangeITM           = "ITM"
	OptionChainRangeNTM           = "NTM"
	OptionChainRangeOTM           = "OTM"
	OptionChainRangeSAK           = "SAK"
	OptionChainRangeSBK           = "SBK"
	OptionChainRangeSNK           = "SNK"
	OptionChainRangeALL           = "ALL"
	OptionChainExpMonthJAN        = "JAN"
	OptionChainExpMonthFEB        = "FEB"
	OptionChainExpMonthMAR        = "MAR"
	OptionChainExpMonthAPR        = "APR"
	OptionChainExpMonthMAY        = "MAY"
	OptionChainExpMonthJUN        = "JUN"
	OptionChainExpMonthJULY       = "JULY"
	OptionChainExpMonthAUG        = "AUG"
	OptionChainExpMonthSEP        = "SEP"
	OptionChainExpMonthOCT        = "OCT"
	OptionChainExpMonthNOV        = "NOV"
	OptionChainExpMonthDEC        = "DEC"
	OptionChainExpMonthALL        = "ALL"
	OptionChainTypeS              = "S"
	OptionChainTypeNS             = "NS"
	OptionChainTypeALL            = "ALL"

	PriceHistoryPeriodTypeDay = "day"
	PriceHistoryPeriodTypeMonth = "month"
	PriceHistoryPeriodTypeYear = "year"
	PriceHistoryPeriodTypeYtd = "ytd"
	PriceHistoryFrequencyTypeMinute = "minute"
	PriceHistoryFrequencyTypeDaily = "daily"
	PriceHistoryFrequencyTypeWeekly = "weekly"
	PriceHistoryFrequencyTypeMonthly = "monthly"

	TransactionsKindALL             = "ALL"
	TransactionsKindTRADE           = "TRADE"
	TransactionsKindBUYONLY         = "BUY_ONLY"
	TransactionsKindSELLONLY        = "SELL_ONLY"
	TransactionsKindCASHINORCASHOUT = "CASH_IN_OR_CASH_OUT"
	TransactionsKindCHECKING        = "CHECKING"
	TransactionsKindDIVIDEND        = "DIVIDEND"
	TransactionsKindINTEREST        = "INTEREST"
	TransactionsKindOTHER           = "OTHER"
	TransactionsKindADVISORFEES     = "ADVISOR_FEES"

	UserPrincipalsFieldsStreamerSubscriptionKeys = "streamerSubscriptionKeys"
	UserPrincipalsFieldsStreamerConnectionInfo   = "streamerConnectionInfo"
	UserPrincipalsFieldsPreferences              = "preferences"
	UserPrincipalsFieldsSurrogateIds             = "surrogateIds"

	OrdersStatusAWAITINGPARENTORDER      = "AWAITING_PARENT_ORDER"
	OrdersStatusAWAITINGCONDITION        = "AWAITING_CONDITION"
	OrdersStatusAWAITINGMANUALREVIEW     = "AWAITING_MANUAL_REVIEW"
	OrdersStatusACCEPTED                 = "ACCEPTED"
	OrdersStatusAWAITINGUROUT            = "AWAITING_UR_OUT"
	OrdersStatusPENDINGACTIVATION        = "PENDING_ACTIVATION"
	OrdersStatusQUEUED                   = "QUEUED"
	OrdersStatusWORKING                  = "WORKING"
	OrdersStatusREJECTED                 = "REJECTED"
	OrdersStatusPENDINGCANCEL            = "PENDING_CANCEL"
	OrdersStatusCANCELED                 = "CANCELED"
	OrdersStatusPENDINGREPLACE           = "PENDING_REPLACE"
	OrdersStatusREPLACED                 = "REPLACED"
	OrdersStatusFILLED                   = "FILLED"
	OrdersStatusEXPIRED                  = "EXPIRED"
	OrderInstructionBuy                  = "BUY"
	OrderInstructionSELL                 = "SELL"
	OrderAssetTypeEQUITY                 = "EQUITY"
	OrderAssetTypeOPTION                 = "OPTION"
	OrderAssetTypeINDEX                  = "INDEX"
	OrderAssetTypeMUTUALFUND             = "MUTUAL_FUND"
	OrderAssetTypeCASHEQUIVALENT         = "CASH_EQUIVALENT"
	OrderAssetTypeFIXEDINCOME            = "FIXED_INCOME"
	OrderAssetTypeCURRENCYOrderAssetType = "CURRENCYOrderAssetType"

	PreferencesTimeoutFIFTYFIVEMINUTES = "FIFTY_FIVE_MINUTES"
	PreferencesTimeoutTWOHOURS         = "TWO_HOURS"
	PreferencesTimeoutFOURHOURS        = "FOUR_HOURS"
	PreferencesTimeoutEIGHTHOURS       = "EIGHT_HOURS"
)

// TokenConfig 命名無法符合 go，因同時為 api 的參數
type TokenConfig struct {
	Grant_type    string // authorization_code, refresh_token
	Refresh_token string
	Access_type   string // Set to offline to receive a refresh token
	Code          string
	Client_id     string
	Redirect_uri  string
}

func (config *TokenConfig) SetAuthorization(code, redirectURI string) *TokenConfig {
	config.Grant_type = "authorization_code"
	config.Access_type = "offline"
	config.Code = code
	config.Redirect_uri = redirectURI

	return config
}

func (config *TokenConfig) SetRefresh(refreshToken string) *TokenConfig {
	config.Grant_type = "refresh_token"
	config.Refresh_token = refreshToken

	return config
}

type PriceHistoryConfig struct {
	Apikey                string
	PeriodType            string //day, month, year, or ytd (year to date). Default is day.
	Period                int64  //day: 1, 2, 3, 4, 5, 10* month: 1*, 2, 3, 6 year: 1*, 2, 3, 5, 10, 15, 20	ytd: 1*
	FrequencyType         string //day: minute* month: daily, weekly*	year: daily, weekly, monthly* ytd: daily, weekly*
	Frequency             int64  //minute: 1*, 5, 10, 15, 30 daily: 1* weekly: 1* monthly: 1*
	EndDate               int64  //End date as milliseconds since epoch If startDate and endDate are provided, period should not be provided. Default is previous trading day.
	StartDate             int64  //Start date as milliseconds since epoch If startDate and endDate are provided, period should not be provided. Default is previous trading day.
	NeedExtendedHoursData *bool
}

func (config *PriceHistoryConfig) SetPeriod(periodType string, period int64) *PriceHistoryConfig {
	config.PeriodType = periodType
	config.Period = period

	return config
}

func (config *PriceHistoryConfig) SetFrequency(frequencyType string, frequency int64) *PriceHistoryConfig {
	config.FrequencyType = frequencyType
	config.Frequency = frequency

	return config
}

func (config *PriceHistoryConfig) SetNeedExtendedHoursData(needExtendedHoursData bool) *PriceHistoryConfig {
	config.NeedExtendedHoursData = &needExtendedHoursData

	return config
}

func (config *PriceHistoryConfig) SetDate(startDate, endDate time.Time) *PriceHistoryConfig {
	config.StartDate = startDate.UnixNano() / int64(time.Millisecond)
	config.EndDate = endDate.UnixNano() / int64(time.Millisecond)

	return config
}

type OptionChainConfig struct {
	Apikey           string
	Symbol           string
	ContractType     string
	StrikeCount      int64
	IncludeQuotes    *bool
	Strategy         string
	Interval         int64
	Strike           float64
	Range            string
	FromDate         string
	ToDate           string
	Volatility       float64
	UnderlyingPrice  float64
	InterestRate     float64
	DaysToExpiration int64
	ExpMonth         string
	OptionType       string
}

func (config *OptionChainConfig) SetSingle() *OptionChainConfig {
	config.Strategy = OptionChainStrategySINGLE

	return config
}

func (config *OptionChainConfig) SetVertical() *OptionChainConfig {
	config.Strategy = OptionChainStrategyVERTICAL

	return config
}

func (config *OptionChainConfig) SetAnalytical(volatility, underlyingPrice, interestRate float64, daysToExpiration int64) *OptionChainConfig {
	config.Strategy = OptionChainStrategyANALYTICAL
	config.Volatility = volatility
	config.UnderlyingPrice = underlyingPrice
	config.InterestRate = interestRate
	config.DaysToExpiration = daysToExpiration

	return config
}

func (config *OptionChainConfig) SetDate(fromDate, toDate time.Time) *OptionChainConfig {
	config.FromDate = fromDate.UTC().Format("2006-01-02")
	config.ToDate = toDate.UTC().Format("2006-01-02")

	return config
}

type TransactionsConfig struct {
	Type      string
	Symbol    string
	StartDate string
	EndDate   string
}

func (config *TransactionsConfig) SetDate(StartDate, EndDate time.Time) *TransactionsConfig {
	config.StartDate = StartDate.UTC().Format("2006-01-02")
	config.EndDate = EndDate.UTC().Format("2006-01-02")

	return config
}

type OrdersConfig struct {
	MaxResults      int64
	FromEnteredTime string
	ToEnteredTime   string
}

func (config *OrdersConfig) SetEnteredTime(fromEnteredTime, toEnteredTime time.Time) *OrdersConfig {
	config.FromEnteredTime = fromEnteredTime.UTC().Format("2006-01-02")
	config.ToEnteredTime = toEnteredTime.UTC().Format("2006-01-02")

	return config
}
