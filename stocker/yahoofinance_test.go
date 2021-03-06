package stocker

import (
	"context"
	"testing"
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
