{{range .Invoices}}# Invoice {{.Code}}

{{range Account .Account}}{{.}}
{{end}}
---

{{range Client .Client}}{{.}}
{{end}}
---

| No. | Description | Amount | Unit | Rate | Total |
| -| - | - | - | - | - |
{{range $i, $val := .Items}}| {{$i}} | {{$val.Name}} | {{.Amount}} | {{.Unit}} | {{.Rate}} | - |
{{end}}
{{end}}