package main

import "testing"

func Test_quoteTD(t *testing.T) {
	type args struct {
		symbol string
	}
	tests := []struct {
		name    string
		args    args
		notWant float64
		wantErr bool
	}{
		// TODO: Add test cases.
		{"VTI", args{"VTI"}, 0, false},
		{"BND", args{"BND"}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := quoteTD(tt.args.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("quoteTD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == tt.notWant {
				t.Errorf("quoteTD() = %v, notWant %v", got, tt.notWant)
			}
		})
	}
}
