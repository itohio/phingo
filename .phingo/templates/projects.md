{{range $p := .Projects}}
# Project {{.Name}}
| Parameter Name | Parameter Value |
| - | - |
| ID | {{.Id}} |
| Client | {{if .Client}}{{.Client.Name}}{{end}} |
| Account | {{if .Account}}{{.Account.Name}}{{end}} |
| Start Date | {{.StartDate}} |
| **Deadline** | **{{.EndDate}}** |
| Rate | {{Rate .}} |
{{range $key, $val := .Params}}| {{$key}} | {{$val}} |{{end}}


Log:
| Description         | Started           | Duration | Progress | Price      |
| --- | --- | --- | --- | --- |{{range .Log}}
| {{.Description}} | {{.Start}} | {{.Duration}}| {{.Progress}} | {{Price $p .}} |{{end}}
| **Total:** | | {{TotalDuration .}} } | {{if .Completed}}**{{end}}{{TotalProgress .}}{{if .Completed}}**{{end}} | {{TotalPrice .}} |

{{end}}