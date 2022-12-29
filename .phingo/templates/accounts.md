{{range .Accounts}}
# Account {{.Name}}
ID: {{.Id}}
Denomination: {{.Denom}}
Decimal places: {{.Decimals}}

Contacts:
{{range $what, $value := .Contact}}- **{{$what}}:** {{$value}}
{{end}}

{{end}}
