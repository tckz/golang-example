{{define "title"}}oreore{{end}}

{{define "content"}}
{{ .title }}
<br>
{{ .title | safeHtml }}
<br>

{{range $num := .numbers}}
	{{partial "_partial.tpl" .}}
{{end}}
{{end}}

