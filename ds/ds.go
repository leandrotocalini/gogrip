package ds

import "regexp"

type Found struct {
	LineNumbers []int
	Match       bool
	FilePath    string
	Content     []string
	Regexp      *regexp.Regexp
}
