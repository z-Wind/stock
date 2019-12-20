package api

import (
	"github.com/z-Wind/stock/api/gotd/api/types"
	"time"
)

// NewOrder default Order
func NewOrder(symbol, assetType, instruction string, price, qunatity float64) *types.Order {
	return &types.Order{
		Session:                  "NORMAL",
		Duration:                 "GOOD_TILL_CANCEL",
		OrderType:                "LIMIT",
		CancelTime:               time.Now().AddDate(0, 4, 0).UTC().Format("2006-01-02"),
		ComplexOrderStrategyType: "NONE",
		Price:                    price,
		OrderLegCollections: []*types.OrderLegCollection{
			&types.OrderLegCollection{
				OrderLegType: assetType,
				InstrumentA: &types.InstrumentA{
					AssetType: assetType,
					Data: &types.Equity{
						Symbol: symbol,
					},
				},
				Instruction: instruction,
				Quantity:    qunatity,
			},
		},
		OrderStrategyType: "SINGLE",
	}
}

// NewPreferences default Preferences
func NewPreferences(timeout string) *types.Preferences {
	return &types.Preferences{
		ExpressTrading:                   false,
		DefaultEquityOrderLegInstruction: "NONE",
		DefaultEquityOrderType:           "LIMIT",
		DefaultEquityOrderPriceLinkType:  "NONE",
		DefaultEquityOrderDuration:       "GOOD_TILL_CANCEL",
		DefaultEquityOrderMarketSession:  "NORMAL",
		DefaultEquityQuantity:            0,
		MutualFundTaxLotMethod:           "FIFO",
		OptionTaxLotMethod:               "FIFO",
		EquityTaxLotMethod:               "FIFO",
		DefaultAdvancedToolLaunch:        "NONE",
		AuthTokenTimeout:                 timeout, //'FIFTY_FIVE_MINUTES', 'TWO_HOURS', 'FOUR_HOURS', 'EIGHT_HOURS'."
	}
}

// NewWatchlistBasic default WatchlistBasic
func NewWatchlistBasic(name string, symbols []string) *types.WatchlistBasic {
	watchlistItems := make([]*types.WatchlistItemBasic, len(symbols))
	for i, symbol := range symbols {
		watchlistItems[i] = &types.WatchlistItemBasic{
			Instrument: &types.InstrumentW{
				Symbol:    symbol,
				AssetType: "EQUITY",
			},
		}
	}

	return &types.WatchlistBasic{
		Name:           name,
		WatchlistItems: watchlistItems,
	}
}
