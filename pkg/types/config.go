package types

import (
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
