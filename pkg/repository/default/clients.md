{{range .Clients}}
# Client {{.Name}}
ID: {{.Id}}
Account ID: {{.Account.Id}}
Account Name: {{.Account.Name}}
Account Denom: {{.Account.Denom}}

Contacts:
{{range $what, $value := .Contact}}
- **{{$what}}:** {{$value}}
{{end}}

{{end}}