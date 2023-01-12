{{range $inv := .Invoices}}# Invoice {{.Code}}

{{range Account $inv.Account}}{{.}}
{{end}}
---

{{range Client $inv.Client}}{{.}}
{{end}}
---

| No. | Description | Amount | Unit | Rate | Total |
| - | - | - | - | - | - |
{{range $i, $val := $inv.Items}}{{if not $val.Extra}}| {{add $i 1}} | {{$val.Name}} | {{.Amount}} | {{.Unit}} | {{.Rate}} | {{Pretty (ItemPrice $inv .)}} |
{{end}}{{end}}
{{$summary := Summary $inv}}
Subtotal: {{Pretty $summary.Subtotal}}

Discount: {{Pretty $summary.Discount}}

Tax: {{Pretty $summary.Tax}}

Total: {{Pretty $summary.Total}}

Total: {{Words $summary.Total}}

---

Taxes:
{{range $summary.Taxes}}- {{.Name}} = {{Pretty .Amount}}
{{end}}

Discounts:
{{range $summary.Discounts}}- {{.Name}} = {{Pretty .Amount}}
{{end}}

{{end}}