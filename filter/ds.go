package filter

import (
	"bufio"
	"os"
	"regexp"
)

type Found struct {
	LineNumbers []int
	Match       bool
	FilePath    string
	Content     []string
	Regexp      *regexp.Regexp
	Scanner     *bufio.Scanner
}

func (f *Found) Filter() {
	for i := 0; f.Scanner.Scan(); i++ {
		line := string(f.Scanner.Text())
		f.Content = append(f.Content, line)
		if f.Regexp.MatchString(line) {
			f.Match = true
			f.LineNumbers = append(f.LineNumbers, i)
		}
	}
}

func createFound(r *regexp.Regexp, file *os.File, filePath string) Found {
	scanner := bufio.NewScanner(file)
	found := Found{Match: false, FilePath: filePath, Regexp: r, Scanner: scanner}
	return found
}
