{{range .Invoices}}# Invoice .Code

- {{.Client.Name}}
- {{.Account.Name}}
{{end}}