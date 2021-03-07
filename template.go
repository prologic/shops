package main

import (
	"bytes"
	"fmt"
	"strings"
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
{{ $val | trim | indent 2 }}
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

func trim(v string) string {
	return strings.TrimSpace(v)
}

func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

func renderString(tpl string, ctx Context) (string, error) {
	funcMap := map[string]interface{}{
		"indent": indent,
		"trim":   trim,
	}

	t := template.New("tpl")
	t.Funcs(funcMap)
	template.Must(t.Parse(tpl))
	buf := bytes.NewBuffer([]byte{})
	err := t.Execute(buf, ctx)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
