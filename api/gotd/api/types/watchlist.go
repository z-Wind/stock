package types

type InstrumentW struct {
	Symbol      string `json:"symbol"`                //"string"
	AssetType   string `json:"assetType"`             //"'EQUITY' or 'OPTION' or 'MUTUAL_FUND' or 'FIXED_INCOME' or 'INDEX'"
	Description string `json:"description,omitempty"` //"string"
}

type WatchlistItemBasic struct {
	Quantity      float64      `json:"quantity,omitempty"`
	AveragePrice  float64      `json:"averagePrice,omitempty"`
	Commission    float64      `json:"commission,omitempty"`
	PurchasedDate string       `json:"purchasedDate,omitempty"` //"DateParam\"",
	Instrument    *InstrumentW `json:"instrument"`
}

type WatchlistBasic struct {
	Name           string                `json:"name"` //"string"
	WatchlistItems []*WatchlistItemBasic `json:"watchlistItems"`
}

type WatchlistItem struct {
	*WatchlistItemBasic

	SequenceID int64  `json:"sequenceId"`
	Status     string `json:"status"` //"'UNCHANGED' or 'CREATED' or 'UPDATED' or 'DELETED'"
}

type Watchlists []*Watchlist

// func (w *Watchlists) Index(i int) *Watchlist {
// 	return []*Watchlist(*w)[i]
// }

// func (w *Watchlists) Len() int {
// 	return len([]*Watchlist(*w))
// }

// func (w *Watchlists) Iterate() []*Watchlist {
// 	return []*Watchlist(*w)
// }

type Watchlist struct {
	*WatchlistBasic

	WatchlistID    string           `json:"watchlistId"`      //"string"
	AccountID      int64            `json:"accountId,string"` //"string"
	Status         string           `json:"status"`           //"'UNCHANGED' or 'CREATED' or 'UPDATED' or 'DELETED'"
	WatchlistItems []*WatchlistItem `json:"watchlistItems"`
}
