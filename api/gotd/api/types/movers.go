package types

type Movers []*Mover

// func (m *Movers) Index(i int) *Mover {
// 	return []*Mover(*m)[i]
// }

// func (m *Movers) Len() int {
// 	return len([]*Mover(*m))
// }

// func (m *Movers) Iterate() []*Mover {
// 	return []*Mover(*m)
// }

type Mover struct {
	Change      float64 `json:"change"`
	Description string  `json:"description"` //"string"
	Direction   string  `json:"direction"`   //"'up' or 'down'"
	Last        float64 `json:"last"`
	Symbol      string  `json:"symbol"` //"string"
	TotalVolume float64 `json:"totalVolume"`
}
