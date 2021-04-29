package main

import (
	"flag"
	"fmt"
	"github.com/leandrotocalini/gogrip/blocks"
	"github.com/leandrotocalini/gogrip/fget"
	"github.com/leandrotocalini/gogrip/filter"
	"github.com/leandrotocalini/gogrip/formatter"
	"runtime"
)

func main() {
	flag.Parse()
	query := flag.Arg(0)
	rootPath := flag.Arg(1)
	buffer := runtime.NumCPU()

	for block := range blocks.Process(filter.Process(query, fget.Get(rootPath, buffer * 2), buffer), buffer) {
		fmt.Println(formatter.Format(block))
	}
}
