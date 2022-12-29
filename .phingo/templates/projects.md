{{range $p := .Projects}}
# Project {{.Name}}
| Parameter Name | Parameter Value |
| - | -: |
| ID | {{.Id}} |
| Client | {{if .Client}}{{.Client.Name}}{{end}} |
| Account | {{if .Account}}{{.Account.Name}}{{end}} |
| Start Date | {{.StartDate}} |
| **Deadline** | **{{if .EndDate}}{{.EndDate}}{{else}}-{{end}}** |
| Rate | {{Rate .}} |
{{range $key, $val := .Params}}| {{$key}} | {{$val}} |{{end}}

{{if .Account}}## Account {{.Account.Name}}
{{range Contacts .Account.Contact}}{{.}}
{{end}}

{{end}}{{if .Client}}## Client {{.Client.Name}}
{{.Client.Description}}

{{range Contacts .Client.Contact}}{{.}}
{{end}}

{{end}}


Log:
| Description         | Started           | Duration | Progress | Price      |
| --- | --- | --- | --- | --- |{{range .Log}}
| {{.Description}} | {{.Start}} | {{Duration .}} | {{.Progress}} | {{Price $p .}} |{{end}}
| **Total:** |   | {{TotalDuration .}} | {{if .Completed}}**{{end}}{{TotalProgress .}}{{if .Completed}}**{{end}} | {{Price $p nil}} |

Total in words: {{TotalPriceWords .}}

{{end}}