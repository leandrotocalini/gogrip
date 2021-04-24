package formatter

import (
	"fmt"
	"github.com/leandrotocalini/gogrip/filter"
)

type Block struct {
	FirstLine int
	LastLine  int
	Lines     map[int]bool
}

func getBlock(doc []string, lineNumber int) (int, int) {
	ff := false
	fl := false
	first := lineNumber
	last := lineNumber
	for !ff || !fl {
		if ff == false {
			if doc[first] == "" {
				ff = true
			} else {
				if first > 0 {
					first--
				} else {
					ff = true
				}
			}
		}
		if fl == false {
			if doc[last] == "" {
				fl = true
			} else {
				if last < len(doc)-1 {
					last++
				} else {
					fl = true
				}
			}
		}
	}
	return first, last
}

func getBlocks(f filter.Found) map[string]*Block {
	blocks := make(map[string]*Block)
	for _, val := range f.LineNumbers {
		first, last := getBlock(f.Content, val)
		key := fmt.Sprintf("%d-%d", first, last)
		_, ok := blocks[key]
		if !ok {
			blocks[key] = &Block{FirstLine: first, LastLine: last, Lines: make(map[int]bool)}
		}
		blocks[key].Lines[val] = true
	}
	return blocks
}
