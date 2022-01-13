package stocker

import (
	"context"
	"reflect"
	"testing"

	yfinance "github.com/z-Wind/yahoofinance"
)

func TestYahooFinance_PriceAdjHistory(t *testing.T) {
	yfinance, err := NewYahooFinance()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		symbol string
	}
	tests := []struct {
		name    string
		yf      *YahooFinance
		args    args
		want    []*DatePrice
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test", yfinance, args{"0050.TW"}, []*DatePrice{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.yf.PriceAdjHistory(context.TODO(), tt.args.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("YahooFinance.PriceAdjHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) <= 0 {
				t.Errorf("YahooFinance.PriceAdjHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYahooFinance_toAdjHistory(t *testing.T) {
	yf, err := NewYahooFinance()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		times   []int64
		history []float64
		event   yfinance.Events
	}
	tests := []struct {
		name    string
		yf      *YahooFinance
		args    args
		want    []float64
		wantErr bool
	}{
		// TODO: Add test cases.
		// standard test
		// {"Standard Test", yf, args{[]int64{216, 217, 218, 219, 220, 221, 222}, []float64{46.99, 48.3, 24.96, 24.91, 24.95, 24.53, 24.54}, yfinance.Events{
		// 	Dividends: map[string]yfinance.Dividend{"221": {Amount: 0.08, Date: 221}},
		// 	Splits: map[string]yfinance.Split{
		// 		"218": {
		// 			Date:        218,
		// 			Numerator:   2,
		// 			Denominator: 1,
		// 			SplitRatio:  "2:1",
		// 		},
		// 	},
		// }}, []float64{23.42, 24.07, 24.88, 24.83, 24.87, 24.53, 24.54}, false},
		{"Test", yf, args{[]int64{1217, 1220, 1221, 1222, 1223, 1227}, []float64{235.44, 232.13, 236.68, 239.02, 240.67, 242.96}, yfinance.Events{
			Dividends: map[string]yfinance.Dividend{"1227": {Amount: 0.859, Date: 1227}},
			Splits:    map[string]yfinance.Split{},
		}}, []float64{234.5996669298209, 231.301480990568, 235.83524111854405, 238.16688918436031, 239.81099999999998, 242.96}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.yf.toAdjHistory(tt.args.times, tt.args.history, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("YahooFinance.toAdjHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("YahooFinance.toAdjHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}
