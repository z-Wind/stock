package types

type Period struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type SessionHours struct {
	PreMarket     []Period `json:"preMarket"`
	RegularMarket []Period `json:"regularMarket"`
	PostMarket    []Period `json:"postMarket"`
}

type MarketHours map[string]map[string]*MarketHour

// func (m *MarketHours) Get(key string) map[string]*MarketHour {
// 	return map[string]map[string]*MarketHour(*m)[key]
// }

// func (m *MarketHours) Len() int {
// 	return len(map[string]map[string]*MarketHour(*m))
// }

// func (m *MarketHours) Iterate() map[string]map[string]*MarketHour {
// 	return map[string]map[string]*MarketHour(*m)
// }

type MarketHour struct {
	Category     string       `json:"category"`
	Date         string       `json:"date"`
	Exchange     string       `json:"exchange"`
	IsOpen       bool         `json:"isOpen"`
	MarketType   string       `json:"marketType"` //"'BOND' or 'EQUITY' or 'ETF' or 'FOREX' or 'FUTURE' or 'FUTURE_OPTION' or 'INDEX' or 'INDICATOR' or 'MUTUAL_FUND' or 'OPTION' or 'UNKNOWN'"
	Product      string       `json:"product"`
	ProductName  string       `json:"productName"`
	SessionHours SessionHours `json:"sessionHours"` //"object"
}
