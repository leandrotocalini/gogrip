package viewer

import (
	"bytes"
	"fmt"
	"github.com/leandrotocalini/gogrip/filter"
	"os"
	"text/template"
)

const fileInfoTmpl = `


------------------------
Path: {{.FilePath}}
Matched lines: {{len .LineNumbers}}
------------------------
`

const lineTmpl = `|{{.LineNumber}}|		{{.Text}} `
const blockTmpl = `{{ range $i, $line := .Content }}{{formatLine $ $i $line}}
{{end}}`

type Line struct {
	Text       string
	Found      bool
	LineNumber int
}

func formatLine(b *Block, i int, line string) string {
	ln := i + b.FirstLine
	found, ok := b.Lines[ln]
	if !ok {
		found = false
	}

	l := &Line{Text: line, Found: found, LineNumber: ln}
	out := &bytes.Buffer{}

	t, err := template.New("lineTmpl").Parse(lineTmpl)
	if err != nil {
		panic(err)
	}
	if err := t.ExecuteTemplate(out, "lineTmpl", &l); err != nil {
		panic(err)
	}
	if l.Found {
		return "\033[31m" + out.String()
	}
	return "\033[37m" + out.String()

}

func (block *Block) String() string {
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

func View(f filter.Found) {
	t := template.Must(template.New("fileInfoTmpl").Parse(fileInfoTmpl))
	t.Execute(os.Stdout, f)
	for _, block := range getBlocks(f) {
		fmt.Println(block)
	}
}
