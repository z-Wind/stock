package stocker

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

// Stocker get stock infomation
type Stocker interface {
	// 得到股票價格
	Quote(ctx context.Context, symbol string) (float64, error)
	// 得到股票歷史價格
	PriceHistory(ctx context.Context, symbol string) ([]*DatePrice, error)
	// 得到股票歷史 Adj 價格
	PriceAdjHistory(ctx context.Context, symbol string) ([]*DatePrice, error)
}

// DatePrice price with date
type DatePrice struct {
	Date     Time    `json:"Date"`
	Open     float64 `json:"Open"`
	High     float64 `json:"High"`
	Low      float64 `json:"Low"`
	Close    float64 `json:"Close"`
	CloseAdj float64 `json:"CloseAdj"`
	Volume   float64 `json:"Volume"`
}

// Time redefine time.Time for JSON
type Time time.Time

const (
	timeFormart = "2006-01-02"
)

// UnmarshalJSON 輸入 json string 格式，兩邊有雙引號
func (t *Time) UnmarshalJSON(data []byte) error {
	result, err := time.Parse(`"`+timeFormart+`"`, string(data))
	if err != nil {
		return errors.Wrapf(err, "time.Parse")
	}

	*t = Time(result)
	return nil
}

// MarshalJSON 輸出 json string 格式，兩邊加上雙引號
func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormart)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(timeFormart)
}
