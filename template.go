package main

import (
	"bytes"
	"strings"
	"text/template"
)

const cmdTmpl = `
set -e
{{ range $key, $val := .Funcs }}
{{ $key }}() {
{{ $val | indent 2 }}
}
{{ end }}
{{ .Command }}
`

type Context struct {
	Funcs   map[string]string
	Command string
}

func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

func renderString(tpl string, ctx Context) (string, error) {
	funcMap := map[string]interface{}{
		"indent": indent,
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
