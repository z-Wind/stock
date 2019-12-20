package twse

import (
	"reflect"
	"testing"
)

func TestTWSE_Quote(t *testing.T) {
	type fields struct {
		QPS int
	}
	type args struct {
		symbol string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		notWant []Trade
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Price", fields{1}, args{"1101"}, []Trade{}, false},
		{"Error", fields{1}, args{"123"}, []Trade{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &TWSE{
				qps: tt.fields.QPS,
			}
			got, err := a.Quote(tt.args.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("TWSE.LatestPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.DeepEqual(got, tt.notWant) && (err == nil) {
				t.Errorf("TWSE.LatestPrice() = %v, want %v", got, tt.notWant)
			}
		})
	}
}

func Test_renameSymbol(t *testing.T) {
	type args struct {
		symbol string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"tse", args{"0050"}, "tse_0050.tw"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renameSymbol(tt.args.symbol); got != tt.want {
				t.Errorf("renameSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTWSE_HistoryMonthly(t *testing.T) {
	type fields struct {
		QPS int
	}
	type args struct {
		symbol string
		year   int
		month  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		notWant []Trade
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Price", fields{1}, args{"1101", 2010, 8}, []Trade{}, false},
		{"Error", fields{1}, args{"123", 2010, 8}, []Trade{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &TWSE{
				qps: tt.fields.QPS,
			}
			got, err := a.HistoryMonthly(tt.args.symbol, tt.args.year, tt.args.month)
			if (err != nil) != tt.wantErr {
				t.Errorf("TWSE.HistoryMonthly() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.DeepEqual(got, tt.notWant) {
				t.Errorf("TWSE.HistoryMonthly() = %v, want %v", got, tt.notWant)
			}
		})
	}
}
