package main

import (
	"flag"
	"fmt"
	"github.com/leandrotocalini/gogrip/filter"
	"github.com/leandrotocalini/gogrip/formatter"

	"runtime"
)

func main() {
	flag.Parse()
	query := flag.Arg(0)
	rootPath := flag.Arg(1)
	buffer := runtime.NumCPU()

	for block := range filter.FilterPath(rootPath, buffer, query) {
		fmt.Println(formatter.Format(block))
	}
}
