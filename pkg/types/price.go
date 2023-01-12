package types

import (
	"fmt"
	"math"
	"strings"
)

func (p *Price) Pretty() string {
	denom := p.Denom
	if d, ok := defaultDenoms[strings.ToUpper(denom)]; ok {
		denom = d
	}
	if len(denom) == 1 {
		denom += " "
	}
	decimals := p.Decimals
	if decimals > 10 {
		decimals = 10
	}
	if decimals == 0 {
		decimals = 2
	}
	format := fmt.Sprintf("%%s%%0.%df", decimals)
	return fmt.Sprintf(format, denom, p.Amount)
}

func (p *Price) Words() string {
	// TODO/FIXME
	floor := math.Floor(float64(p.Amount))
	remainder := float64(p.Amount) - floor
	aboveZero := Num2Words(int(floor))
	belowFloor := math.Floor(remainder * math.Pow10(int(p.Decimals)))
	belowZero := Num2Words(int(belowFloor))
	return aboveZero + " " + p.Denom + " And " + belowZero + " Ct"
}

func (p *Price) Sub(b *Price) *Price {
	return &Price{
		Amount:   p.Amount - b.Amount,
		Denom:    p.Denom,
		Decimals: p.Decimals,
	}
}

func (p *Price) Add(b *Price) *Price {
	return &Price{
		Amount:   p.Amount + b.Amount,
		Denom:    p.Denom,
		Decimals: p.Decimals,
	}
}

func (p *Price) Mul(rate float32) *Price {
	return &Price{
		Amount:   p.Amount * rate,
		Denom:    p.Denom,
		Decimals: p.Decimals,
	}
}

// Convert converts one denomination to another given that
// the second is represented in the first denomination.
// E.g. first price is 10Eth while the second is 4000Eur (per Eth)
// Therefore 10Eth.Convert(4000Eur) will be equal to 40000Eur
func (p *Price) Convert(b *Price) *Price {
	return &Price{
		Amount:   p.Amount * b.Amount,
		Denom:    b.Denom,
		Decimals: b.Decimals,
	}
}
