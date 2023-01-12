package engine

import (
	"text/template"
	"time"

	"github.com/itohio/phingo/pkg/types"
)

func addDefaultFuncs(m template.FuncMap) template.FuncMap {
	ret := make(template.FuncMap, len(m)+10)
	for k, v := range m {
		ret[k] = v
	}

	defaults := template.FuncMap{
		"sub": func(y, x int) int {
			return x - y
		},
		"add": func(y, x int) int {
			return x + y
		},
		"get": func(m map[string]string, key, fallback string) string {
			if v, ok := m[key]; ok {
				return v
			}
			return fallback
		},
		"duration": func(d time.Duration) string {
			return d.String()
		},
		"date": func(t time.Time) string {
			return t.Format("2006-11-12")
		},
		"time": func(t time.Time) string {
			return t.Format("17:45")
		},
		"since": func(t time.Time) time.Duration {
			return time.Since(t)
		},
		"now": func() time.Time {
			return time.Now()
		},
		"Pretty": func(price *types.Price) string {
			if price == nil {
				return "-"
			}
			return price.Pretty()
		},
		"Words": func(price *types.Price) string {
			if price == nil {
				return "-"
			}
			return price.Words()
		},
	}

	for k, v := range defaults {
		ret[k] = v
	}
	return ret
}
