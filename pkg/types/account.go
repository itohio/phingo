package types

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

type AccountsArr []*Account

func (arr AccountsArr) ById(id string) *Account {
	for _, a := range arr {
		if a == nil {
			continue
		}
		if a.Id == id {
			return a
		}
	}
	return nil
}

func (arr AccountsArr) ByName(name string) *Account {
	for _, a := range arr {
		if a == nil {
			continue
		}
		if a.Name == name {
			return a
		}
	}
	return nil
}

func (arr AccountsArr) ByDenom(denom string) *Account {
	for _, a := range arr {
		if a == nil {
			continue
		}
		if a.Denom == denom {
			return a
		}
	}
	return nil
}

func (a *Account) MakeId(salt string) string {
	buf := sha1.Sum([]byte(a.Denom + a.Name + salt))
	return base64.RawStdEncoding.EncodeToString(buf[:])
}

func (a *Account) MakeInvoiceCode(inv []*Invoice, totalCount int, now time.Time) string {
	format := "{MON} {Count}/{Day}"

	tokens := map[string]func() string{
		"{Month}": func() string {
			return now.Month().String()
		},
		"{Mon}": func() string {
			return now.Month().String()[:3]
		},
		"{MON}": func() string {
			return strings.ToUpper(now.Month().String()[:3])
		},
		"{Day}": func() string {
			return fmt.Sprint(now.Day())
		},
		"{Count}": func() string {
			return fmt.Sprint(len(inv) + 1)
		},
		"{Total Count}": func() string {
			return fmt.Sprint(totalCount + 1)
		},
	}
	for k, val := range tokens {
		format = strings.ReplaceAll(format, k, val())
	}
	return format
}
