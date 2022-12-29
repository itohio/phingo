package types

import (
	"fmt"
	"math"
	"strings"

	"google.golang.org/protobuf/proto"
)

var (
	defaultDenoms = map[string]string{
		"EUR": "€",
		"USD": "$",
		"GBP": "£",
	}
)

func (c *Config) Translate(val string) string {
	if c.Locale == nil {
		return val
	}
	return c.Locale.Translate(val)
}

func (c *Locale) Translate(val string) string {
	if c == nil {
		return val
	}
	if c.Translations == nil {
		return val
	}
	if s, ok := c.Translations[strings.ToUpper(val)]; ok {
		return s
	}
	return val
}

func (c *Config) Format(price *Price) string {
	if c.Locale == nil {
		return price.Pretty()
	}
	return c.Locale.Format(price)
}

func (c *Locale) Format(price *Price) string {
	tmp := proto.Clone(price).(*Price)
	denom := strings.ToTitle(price.Denom)
	if s, ok := c.Translations[denom]; ok {
		tmp.Denom = s
	} else if s, ok := defaultDenoms[denom]; ok {
		tmp.Denom = s
	}
	return tmp.Pretty()
}

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
