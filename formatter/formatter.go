package formatter

import (
	"bytes"
	"github.com/leandrotocalini/gogrip/blocks"
	"text/template"
)

const lineTmpl = `{{if .Found}}[31m{{else}}[37m{{end}}{{.FilePath}}	{{.LineNumber}}|		{{.Text}} `
const blockTmpl = `{{ range $i, $line := .Content }}{{formatLine $ $i $line}}
{{end}}`

func formatLine(b *blocks.Block, i int, line string) string {
	ln := i + b.FirstLine
	found, ok := b.Lines[ln]
	if !ok {
		found = false
	}
	type Line struct {
		Text       string
		Found      bool
		LineNumber int
		FilePath   string
	}
	l := &Line{Text: line, Found: found, LineNumber: ln + 1, FilePath: b.FilePath}
	out := &bytes.Buffer{}
	t, err := template.New("lineTmpl").Parse(lineTmpl)
	if err != nil {
		panic(err)
	}
	if err := t.ExecuteTemplate(out, "lineTmpl", &l); err != nil {
		panic(err)
	}
	return out.String()

}

func Format(block blocks.Block) string{
	out := &bytes.Buffer{}
	funcMap := template.FuncMap{
		"formatLine": formatLine,
	}
	t, err := template.New("blockTmpl").Funcs(funcMap).Parse(blockTmpl)
	if err != nil {
		panic(err)
	}
	if err := t.ExecuteTemplate(out, "blockTmpl", &block); err != nil {
		panic(err)
	}
	return out.String()
}

