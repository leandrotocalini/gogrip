package main

import (
	"fmt"
	"github.com/fatih/color"
)

func defineBlock(doc []string, lineNumber int) (int, int) {
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

func PrintLine(doc []string, first int, last int, lineNumbers map[int]bool) {
	normal := color.New(color.FgWhite).SprintFunc()
	highlight := color.New(color.FgRed).SprintFunc()
	for i := first; i <= last; i++ {
		_, ok := lineNumbers[i]
		if !ok {
			fmt.Printf("%s - %s \n", normal(i), normal(doc[i]))
		} else {
			fmt.Printf("%s - %s \n", normal(i), highlight(doc[i]))
		}
	}
}

type ViewStruct struct {
	First int
	Last  int
	Lines map[int]bool
}

func printTitle(f Found) {
	fmt.Println("\n\n\n")
	red := color.New(color.FgRed)
	boldRed := red.Add(color.Bold)
	boldRed.Println(f.FilePath)

}

func View(f Found) {
	printTitle(f)
	h := make(map[string]*ViewStruct)
	for _, val := range f.LineNumbers {
		first, last := defineBlock(f.Content, val)
		key := fmt.Sprintf("%d-%d", first, last)
		_, ok := h[key]
		if !ok {
			h[key] = &ViewStruct{First: first, Last: last, Lines: make(map[int]bool)}
		}
		h[key].Lines[val] = true
	}
	for _, v := range h {
		PrintLine(f.Content, v.First, v.Last, v.Lines)
	}
}
