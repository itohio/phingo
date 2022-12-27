{{range .Accounts}}
# Account {{.Name}}
ID: {{.Id}}
Denomination: {{.Denom}}

Contacts:
{{range $what, $value := .Contact}}
- **{{$what}}:** {{$value}}{{end}}

{{end}}
