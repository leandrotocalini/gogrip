package main

import (
	"flag"
	"github.com/leandrotocalini/gogrip/fget"
	"github.com/leandrotocalini/gogrip/filter"
	"github.com/leandrotocalini/gogrip/viewer"
	"runtime"
)

func main() {
	flag.Parse()
	query := flag.Arg(0)
	rootPath := flag.Arg(1)
	buffer := runtime.NumCPU()

	for f := range filter.Process(query, fget.Get(rootPath, buffer), buffer) {
		viewer.View(f)
	}
}
