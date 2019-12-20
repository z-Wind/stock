package main

// type price struct {
// 	Symbol string  `json:"Symbol"`
// 	Price  float64 `json:"Price"`
// }

type savedOrderParas struct {
	Symbol      string  `json:"Symbol"`
	AssetType   string  `json:"AssetType"`
	Instruction string  `json:"Instruction"`
	Price       float64 `json:"Price"`
	Qunatity    float64 `json:"Qunatity"`
}

type datePrices []*datePrice
type datePrice struct {
	Date     string  `json:"Date"`
	Open     float64 `json:"Open"`
	High     float64 `json:"High"`
	Low      float64 `json:"Low"`
	Close    float64 `json:"Close"`
	CloseAdj float64 `json:"CloseAdj"`
	Volume   float64 `json:"Volume"`
}
