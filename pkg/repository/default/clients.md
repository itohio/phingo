{{range .Clients}}
# Client {{.Name}}
ID: {{.Id}}
Name: {{.Name}}
Description: {{.Description}}
{{if .Notes}}
Notes:
{{range Notes .Notes}}- {{.}}
{{end}}{{end}}

Contacts:
{{range Contacts .Contact}}- {{.}}
{{end}}

{{end}}