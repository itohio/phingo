package types

import (
	"testing"
)

func Test_convert99(t *testing.T) {
	tests := []struct {
		name string
		args int
	}{
		{"One", 1},
		{"Two", 2},
		{"Zero", 0},
		{"Twelve", 12},
		{"Forty-Two", 42},
		{"Seventy-Nine", 79},
		{"Ninety-Nine", 99},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convert99(tt.args); got != tt.name {
				t.Errorf("convert99() = %v, want %v", got, tt.name)
			}
		})
	}
}

func Test_convert999(t *testing.T) {
	tests := []struct {
		name string
		args int
	}{
		{"One", 1},
		{"Two", 2},
		{"Zero", 0},
		{"Twelve", 12},
		{"Forty-Two", 42},
		{"Seventy-Nine", 79},
		{"Ninety-Nine", 99},
		{"One Hundred", 100},
		{"One Hundred Thirteen", 113},
		{"Five Hundred Seven", 507},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convert999(tt.args); got != tt.name {
				t.Errorf("convert999() = %v, want %v", got, tt.name)
			}
		})
	}
}

func TestNum2Words(t *testing.T) {
	tests := []struct {
		name string
		args int
	}{
		{"One", 1},
		{"Two", 2},
		{"Zero", 0},
		{"Twelve", 12},
		{"Forty-Two", 42},
		{"Seventy-Nine", 79},
		{"Ninety-Nine", 99},
		{"One Hundred", 100},
		{"One Hundred Thirteen", 113},
		{"Five Hundred Seven", 507},
		{"Nine Hundred Ninety-Nine", 999},
		{"One Thousand Twenty-Four", 1_024},
		{"Sixteen Thousand Three Hundred Eighty-Six", 16_386},
		{"One Hundred Fifty Thousand Seven Hundred Eighty", 150_780},
		{"One Hundred Fifty-One Thousand Seven Hundred Eighty-Three", 151_783},
		{"Four Hundred Forty-Six Million Five Hundred Seventy-Seven Thousand Seven Hundred Eighty-Three", 446_577_783},
		{"Five Hundred Fifteen Billion Four Hundred Forty-Six Million Five Hundred Seventy-Seven Thousand Seven Hundred Eighty-Three", 515_446_577_783},
		{"More than a trillion", 1_515_446_577_783},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Num2Words(tt.args); got != tt.name {
				t.Errorf("Num2Words() = %v, want %v", got, tt.name)
			}
		})
	}
}
