package output

import (
	"bytes"
	"text/template"
)

var resultTmpl = `
{{define "mark"}}{{if .Success}}✓{{else}}✗{{end}}{{end}}
{{define "file"}}{{with .FileName -}}[{{.}}] {{end}}{{end}}
{{define "tries"}}{{if gt .Tries 1}}{{.Tries}}{{end}}{{end}}
{{define "baseResult"}}{{template "mark" .}} {{template "file" .}}[{{ .Node }}] {{ .Title }}{{end}}`

var commanderTmpl = `
{{define "duration"}}Duration: {{printf "%.3fs" .Duration.Seconds}}{{end}}
{{define "summary"}}Count: {{len .TestResults}}, Failed: {{ .Failed }}{{end}}
{{define "result"}}{{template "baseResult" .}}{{template "tries" .}}{{end}}`

// cliTemplate template object for all of commanders cli tenplate
// each methood is designed to be a wrapper on each commander template
type cliTemplate struct {
	template *template.Template
}

func newCliTemplate() cliTemplate {
	t := template.Must(template.New("").Parse(resultTmpl))
	t = template.Must(t.New("").Parse(commanderTmpl))

	return cliTemplate{
		template: t,
	}
}

func (t cliTemplate) duration(result Result) string {
	tpl := t.getTemplatedString("duration", result)
	return tpl.String()
}

func (t cliTemplate) summary(result Result) string {
	tpl := t.getTemplatedString("summary", result)
	return tpl.String()
}

func (t cliTemplate) testResult(testResult TestResult) string {
	tpl := t.getTemplatedString("result", testResult)
	return tpl.String()
}

func (t cliTemplate) getTemplatedString(name string, data interface{}) bytes.Buffer {
	var tpl bytes.Buffer
	err := t.template.ExecuteTemplate(&tpl, name, data)
	if err != nil {
		panic(err)
	}

	return tpl
}
