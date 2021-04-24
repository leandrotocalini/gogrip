package viewer

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/leandrotocalini/gogrip/filter"
)

func PrintBlock(doc []string, first int, last int, lineNumbers map[int]bool) {
	normal := color.New(color.FgWhite).SprintFunc()
	red := color.New(color.FgMagenta)
	boldRed := red.Add(color.Bold)
	highlight := boldRed.SprintFunc()
	for i := first; i <= last; i++ {
		_, ok := lineNumbers[i]
		fmt.Printf("%s\t", normal(i+1, "|"))
		if !ok {
			fmt.Printf("%s\n", normal(doc[i]))
		} else {
			fmt.Printf("%s\n", highlight(doc[i]))
		}
	}
}

func PrintInfo(f filter.Found) {
	fmt.Println("\n\n\n")
	red := color.New(color.FgWhite)
	boldRed := red.Add(color.Underline)
	boldRed.Printf("%s:\n\n", f.FilePath)

}

func View(f filter.Found) {
	PrintInfo(f)
	blocks := getBlocks(f)
	for _, v := range blocks {
		PrintBlock(f.Content, v.FirstLine, v.LastLine, v.Lines)
		fmt.Println("\n")

	}
}
