package output

import (
	"bytes"
	"text/template"
)

var resultTmpl = `
// Mark Template
{{define "mark" -}}
	{{if .Success -}}
		✓
	{{- else -}}
		✗
	{{- end}}
{{- end -}}

// Add File Template
{{define "file" -}}
	{{with .FileName -}}
		[{{.}}]
	{{- end}}
{{- end -}}

// Add Tries Template
{{define "tries" -}}
	{{if gt .Tries 1 -}}
		{{.Tries -}}
	{{- end}}
{{- end -}}

// BaseResult
{{define "baseResult" -}}
	{{template "mark" .}} {{template "file" .}} [{{ .Node }}]
{{- end -}}`

var commanderTmpl = `
// Duration
{{define "duration" -}}
	Duration: {{printf "%.3fs" .Duration.Seconds}}
{{- end -}}

// Summary
{{define "summary" -}}
	Count: {{len .TestResults}}, Failed: {{ .Failed }}
{{- end -}}

// Result
{{define "result" -}}
	{{template "baseResult" .}}{{ .Title }}{{template "tries" .}}
{{- end -}}

// Failure
{{- define "failure" -}}
	{{- template "baseResult" .}} '{{ .Title }}', on property {{ .FailedProperty }}
{{- end -}}

// Error
{{define "error" -}}
		{{template "baseResult" .}} '{{ .Title }}' could not be executed with error message:
{{- end}}`

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

func (t cliTemplate) failures(testResult TestResult) string {
	tpl := t.getTemplatedString("failure", testResult)
	return tpl.String()
}

func (t cliTemplate) errors(testResult TestResult) string {
	tpl := t.getTemplatedString("error", testResult)
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
