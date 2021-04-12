package main

import (
	"fmt"
	"github.com/fatih/color"
)

type Block struct {
	FirstLine int
	LastLine  int
	Lines     map[int]bool
}

func DefineBlock(doc []string, lineNumber int) (int, int) {
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

func PrintBlock(doc []string, first int, last int, lineNumbers map[int]bool) {
	normal := color.New(color.FgWhite).SprintFunc()
	red := color.New(color.FgMagenta)
	boldRed := red.Add(color.Bold)
	highlight := boldRed.SprintFunc()
	for i := first; i <= last; i++ {
		_, ok := lineNumbers[i]
		fmt.Printf("%s\t", normal(i, "|"))
		if !ok {
			fmt.Printf("%s\n", normal(doc[i]))
		} else {
			fmt.Printf("%s\n", highlight(doc[i]))
		}
	}
}

func PrintInfo(f Found) {
	fmt.Println("\n\n\n")
	red := color.New(color.FgWhite)
	boldRed := red.Add(color.Underline)
	boldRed.Printf("%s:\n\n", f.FilePath)

}

func View(f Found) {
	PrintInfo(f)
	blocks := make(map[string]*Block)
	for _, val := range f.LineNumbers {
		first, last := DefineBlock(f.Content, val)
		key := fmt.Sprintf("%d-%d", first, last)
		_, ok := blocks[key]
		if !ok {
			blocks[key] = &Block{FirstLine: first, LastLine: last, Lines: make(map[int]bool)}
		}
		blocks[key].Lines[val] = true
	}
	for _, v := range blocks {
		PrintBlock(f.Content, v.FirstLine, v.LastLine, v.Lines)
		fmt.Println("\n")

	}
}
