---
pageSize: "a4"
permissions: 
    - Print
    - Copy
signed: true
trace: true
---
{{$length := len .Invoices}}{{range $i, $v := .Invoices}}# INVOICE Serial Number {{$v.Code}} {style="text-align: center" {{if not (eq $length 1)}}{{if not (eq $i 0)}}page-break="true"{{end}} }
### {{$v.IssueDate}} {style="text-align: center"}

:::{#{{$i}} style="width: 100%; display: flex; justify-content: space-between"}
:::{#{{$i}}1}
{{range Account $v.Account}}{{.}}<br>
{{end}}
:::
:::{#{{$i}}2}
{{range Client $v.Client}}{{.}}

{{end}}
:::
:::

### Services {style="text-align: center"}
| No. | Description | Amount | Unit | Rate | Total |
| -| - | - | - | - | - |
{{range $j, $val := .Items}}| {{$j}} | {{$val.Name}} | {{.Amount}} | {{.Unit}} | {{.Rate}} | - |
{{end}}

:::{#{{$i}}total style="width: 100%; display: flex; justify-content: right"}
:::{#{{$i}}subtotal style="width: 300px"}
Subtotal: <br>

Tax: <br>

Discount: <br>

Total: <br>

---

:::
:::

Total in words: 

{{end}}{{end}}