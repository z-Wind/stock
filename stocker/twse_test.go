package stocker

import (
	"context"
	"reflect"
	"testing"
)

func TestTWSE_Quote(t *testing.T) {
	twse, err := NewTWSE()
	if err != nil {
		t.Fatal(err)
	}
	twse.symbolsTWSE_Path = "../TWSE.csv"
	twse.symbolsTPEx_Path = "../TPEx.csv"

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
			got, err := tt.tr.Quote(context.TODO(), tt.args.symbol)
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

func TestTWSE_getSymbolList(t *testing.T) {
	twse, err := NewTWSE()
	if err != nil {
		t.Fatal(err)
	}
	twse.symbolsTWSE_Path = "../TWSE.csv"
	twse.symbolsTPEx_Path = "../TPEx.csv"

	type args struct {
		url string
	}
	tests := []struct {
		name    string
		t       *TWSE
		args    args
		want    map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"TWSE", twse, args{TWSE_url}, nil, false},
		{"OTC", twse, args{TPEx_url}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.getSymbolList(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("TWSE.getSymbolList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tt.want) > 0 {
				t.Errorf("TWSE.getSymbolList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTWSE_isInTWSE(t *testing.T) {
	twse, err := NewTWSE()
	if err != nil {
		t.Fatal(err)
	}
	twse.symbolsTWSE_Path = "../TWSE.csv"
	twse.symbolsTPEx_Path = "../TPEx.csv"

	type args struct {
		symbol string
	}
	tests := []struct {
		name string
		t    *TWSE
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"test", twse, args{"0050.TW"}, true},
		{"test", twse, args{"006208.TW"}, true},
		{"test", twse, args{"006208"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.isInTWSE(tt.args.symbol); got != tt.want {
				t.Errorf("TWSE.isInTWSE() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTWSE_isInTPEx(t *testing.T) {
	twse, err := NewTWSE()
	if err != nil {
		t.Fatal(err)
	}
	twse.symbolsTWSE_Path = "../TWSE.csv"
	twse.symbolsTPEx_Path = "../TPEx.csv"

	type args struct {
		symbol string
	}
	tests := []struct {
		name string
		t    *TWSE
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"test", twse, args{"712435.TWO"}, true},
		{"test", twse, args{"73338P.TWO"}, true},
		{"test", twse, args{"712435"}, true},
		{"test", twse, args{"73338P"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.isInTPEx(tt.args.symbol); got != tt.want {
				t.Errorf("TWSE.isInTPEx() = %v, want %v", got, tt.want)
			}
		})
	}
}
