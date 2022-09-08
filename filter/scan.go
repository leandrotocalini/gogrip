package filter

import (
	"regexp"
	"strings"
)

type Found struct {
	LineNumbers []int
	Match       bool
	FilePath    string
	Content     []string
}

type FScanner interface {
	Scan() bool
	Text() string
}

func filterScanner(scanner FScanner, query string) ([]string, []int) {
	r, _ := regexp.Compile(query)
	i := 0
	content := make([]string, 0)
	lineNumbers := make([]int, 0)
	for scanner.Scan() {
		text := strings.Replace(scanner.Text(), "\t", "  ", 10)
		content = append(content, text)
		if r.MatchString(text) {
			lineNumbers = append(lineNumbers, i)
		}
		i++
	}

	return content, lineNumbers
}

func scanFile(scanner FScanner, filePath string, query string) Found {
	f := Found{Match: false, FilePath: filePath}
	content, lineNumbers := filterScanner(scanner, query)
	if len(lineNumbers) > 0 {
		f.Content = content
		f.LineNumbers = lineNumbers
		f.Match = true
	}
	return f
}
