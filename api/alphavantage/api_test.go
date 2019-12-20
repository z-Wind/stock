package alphavantage

import (
	"math/rand"
	"testing"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	keyLen      = 3
)

var (
	n = 0
)

func randkey(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	if n > 0 {
		time.Sleep(time.Second * 15)
	}
	n++
	return string(b)
}

func TestAlphaVantage_Intraday(t *testing.T) {
	type fields struct {
		Key string
	}
	type args struct {
		symbol     string
		interval   string
		outputsize string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TimeSeries
		wantErr bool
	}{
		// TODO: Add test cases.
		// normal
		{"1min", fields{randkey(keyLen)}, args{"VTI", "1min", OutputsizeCompact}, &TimeSeries{}, false},
		{"5min", fields{randkey(keyLen)}, args{"VTI", "5min", OutputsizeCompact}, &TimeSeries{}, false},
		{"30min", fields{randkey(keyLen)}, args{"VTI", "30min", OutputsizeCompact}, &TimeSeries{}, false},
		{"60min", fields{randkey(keyLen)}, args{"VTI", "60min", OutputsizeCompact}, &TimeSeries{}, false},

		// error
		{"Information", fields{"demo"}, args{"VTI", "1min", OutputsizeCompact}, &TimeSeries{}, true},
		{"Error Message", fields{randkey(keyLen)}, args{"VTI", "2min", OutputsizeCompact}, &TimeSeries{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AlphaVantage{
				key: tt.fields.Key,
			}
			got, err := a.Intraday(tt.args.symbol, tt.args.interval, tt.args.outputsize)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlphaVantage.Intraday() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil &&
				len(got.TimeSeries1min) == 0 &&
				len(got.TimeSeries5min) == 0 &&
				len(got.TimeSeries30min) == 0 &&
				len(got.TimeSeries60min) == 0 {
				t.Errorf("AlphaVantage.Intraday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlphaVantage_Daily(t *testing.T) {
	type fields struct {
		Key string
	}
	type args struct {
		symbol     string
		outputsize string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TimeSeries
		wantErr bool
	}{
		// TODO: Add test cases.
		// normal
		{"Daily", fields{randkey(keyLen)}, args{"VTI", OutputsizeCompact}, &TimeSeries{}, false},

		// error
		{"Information", fields{"demo"}, args{"VTI", OutputsizeCompact}, &TimeSeries{}, true},
		{"Error Message", fields{randkey(keyLen)}, args{"VTIaa", OutputsizeCompact}, &TimeSeries{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AlphaVantage{
				key: tt.fields.Key,
			}
			got, err := a.Daily(tt.args.symbol, tt.args.outputsize)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlphaVantage.Daily() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil &&
				len(got.TimeSeriesDaily) == 0 {
				t.Errorf("AlphaVantage.Daily() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlphaVantage_DailyAdj(t *testing.T) {
	type fields struct {
		Key string
	}
	type args struct {
		symbol     string
		outputsize string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TimeSeries
		wantErr bool
	}{
		// TODO: Add test cases.
		{"DailyAdj", fields{randkey(keyLen)}, args{"VTI", OutputsizeCompact}, &TimeSeries{}, false},

		// error
		{"Information", fields{"demo"}, args{"VTI", OutputsizeCompact}, &TimeSeries{}, true},
		{"Error Message", fields{randkey(keyLen)}, args{"VTIaa", OutputsizeCompact}, &TimeSeries{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AlphaVantage{
				key: tt.fields.Key,
			}
			got, err := a.DailyAdj(tt.args.symbol, tt.args.outputsize)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlphaVantage.DailyAdj() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil &&
				len(got.TimeSeriesDaily) == 0 {
				t.Errorf("AlphaVantage.DailyAdj() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlphaVantage_Weekly(t *testing.T) {
	type fields struct {
		Key string
	}
	type args struct {
		symbol     string
		outputsize string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TimeSeries
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Weekly", fields{randkey(keyLen)}, args{"VTI", OutputsizeCompact}, &TimeSeries{}, false},

		// error
		{"Information", fields{"demo"}, args{"VTI", OutputsizeCompact}, &TimeSeries{}, true},
		{"Error Message", fields{randkey(keyLen)}, args{"VTIaa", OutputsizeCompact}, &TimeSeries{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AlphaVantage{
				key: tt.fields.Key,
			}
			got, err := a.Weekly(tt.args.symbol, tt.args.outputsize)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlphaVantage.Weekly() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil &&
				len(got.TimeSeriesWeekly) == 0 {
				t.Errorf("AlphaVantage.Weekly() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlphaVantage_WeeklyAdj(t *testing.T) {
	type fields struct {
		Key string
	}
	type args struct {
		symbol     string
		outputsize string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TimeSeries
		wantErr bool
	}{
		// TODO: Add test cases.
		{"WeeklyAdj", fields{randkey(keyLen)}, args{"VTI", OutputsizeCompact}, &TimeSeries{}, false},

		// error
		{"Information", fields{"demo"}, args{"VTI", OutputsizeCompact}, &TimeSeries{}, true},
		{"Error Message", fields{randkey(keyLen)}, args{"VTIaa", OutputsizeCompact}, &TimeSeries{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AlphaVantage{
				key: tt.fields.Key,
			}
			got, err := a.WeeklyAdj(tt.args.symbol, tt.args.outputsize)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlphaVantage.WeeklyAdj() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil &&
				len(got.TimeSeriesWeeklyAdj) == 0 {
				t.Errorf("AlphaVantage.WeeklyAdj() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlphaVantage_Monthly(t *testing.T) {
	type fields struct {
		Key string
	}
	type args struct {
		symbol     string
		outputsize string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TimeSeries
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Monthly", fields{randkey(keyLen)}, args{"VTI", OutputsizeCompact}, &TimeSeries{}, false},

		// error
		{"Information", fields{"demo"}, args{"VTI", OutputsizeCompact}, &TimeSeries{}, true},
		{"Error Message", fields{randkey(keyLen)}, args{"VTIaa", OutputsizeCompact}, &TimeSeries{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AlphaVantage{
				key: tt.fields.Key,
			}
			got, err := a.Monthly(tt.args.symbol, tt.args.outputsize)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlphaVantage.Monthly() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil &&
				len(got.TimeSeriesMonthly) == 0 {
				t.Errorf("AlphaVantage.Monthly() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlphaVantage_MonthlyAdj(t *testing.T) {
	type fields struct {
		Key string
	}
	type args struct {
		symbol     string
		outputsize string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TimeSeries
		wantErr bool
	}{
		// TODO: Add test cases.
		{"MonthlyAdj", fields{randkey(keyLen)}, args{"VTI", OutputsizeCompact}, &TimeSeries{}, false},

		// error
		{"Information", fields{"demo"}, args{"VTI", OutputsizeCompact}, &TimeSeries{}, true},
		{"Error Message", fields{randkey(keyLen)}, args{"VTIaa", OutputsizeCompact}, &TimeSeries{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AlphaVantage{
				key: tt.fields.Key,
			}
			got, err := a.MonthlyAdj(tt.args.symbol, tt.args.outputsize)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlphaVantage.MonthlyAdj() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil &&
				len(got.TimeSeriesMonthlyAdj) == 0 {
				t.Errorf("AlphaVantage.MonthlyAdj() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlphaVantage_LatestPrice(t *testing.T) {
	type fields struct {
		Key string
	}
	type args struct {
		symbol string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"LastPrice", fields{randkey(keyLen)}, args{"VTI"}, "VTI", false},

		// error
		{"Information", fields{"demo"}, args{"VTI"}, "", true},
		{"empty", fields{randkey(keyLen)}, args{"VTIaa"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AlphaVantage{
				key: tt.fields.Key,
			}
			got, err := a.LatestPrice(tt.args.symbol)
			t.Log(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlphaVantage.LatestPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil &&
				got.Symbol != tt.want {
				t.Errorf("AlphaVantage.LatestPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
