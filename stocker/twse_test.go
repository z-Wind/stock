package stocker

import (
	"reflect"
	"testing"
)

func TestTWSE_Quote(t *testing.T) {
	twse, err := NewTWSE()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		symbol string
	}
	tests := []struct {
		name    string
		tr      *TWSE
		args    args
		want    reflect.Kind
		wantErr bool
	}{
		// TODO: Add test cases.
		{"TWSE", twse, args{"0050.TW"}, reflect.Float64, false},
		{"TWSE", twse, args{"006208.tw"}, reflect.Float64, false},
		{"OTC", twse, args{"6237.TWO"}, reflect.Float64, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.Quote(tt.args.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("TWSE.Quote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got).Kind() != tt.want {
				t.Errorf("TWSE.Quote() = %v, want %v", got, tt.want)
			}
		})
	}
}
