package main

import (
	"bytes"
	"fmt"
	"text/template"
)

const cmdTmpl = `
set -e

{{- with .Env }}
{{ range $env := . }}
{{ $env.Key }}="{{ $env.Value }}"
{{- end }}

export {{ range $env := . }}{{ $env.Key }} {{ end }}
{{- end }}

{{- with .Funcs }}
{{ range $key, $val := . }}
{{ $key }}() {
{{ $val }}
}
{{- end }}
{{- end }}

{{ .Command }}
`

type Env struct {
	Key, Value string
}

func (e *Env) String() string {
	return fmt.Sprintf("%s=%s", e.Key, e.Value)
}

type Context struct {
	Env     []*Env
	Funcs   map[string]string
	Command string
}

func renderString(tpl string, ctx Context) (string, error) {
	t := template.New("tpl")
	template.Must(t.Parse(tpl))
	buf := bytes.NewBuffer([]byte{})
	err := t.Execute(buf, ctx)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
