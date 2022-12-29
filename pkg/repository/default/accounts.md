{{range .Accounts}}
# Account {{.Name}}
ID: {{.Id}}
Denomination: {{.Denom}}
Decimal places: {{.Decimals}}

Contacts:
{{range Contacts .Contact}}- {{.}}
{{end}}

{{end}}
