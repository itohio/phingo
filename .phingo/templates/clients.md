{{range .Clients}}
# Client {{.Name}}
ID: {{.Id}}
Name: {{.Name}}
Description: {{.Description}}
{{if .Notes}}Notes:
{{range $what, $value := .Notes}}
- **{{$what}}:** {{$value}}
{{end}}{{end}}

Contacts:
{{range $what, $value := .Contact}}- **{{$what}}:** {{$value}}
{{end}}

{{end}}