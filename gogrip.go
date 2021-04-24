package main

import (
	"flag"
	"github.com/leandrotocalini/gogrip/fget"
	"github.com/leandrotocalini/gogrip/filter"
	"github.com/leandrotocalini/gogrip/viewer"
)

func main() {
	flag.Parse()
	query := flag.Arg(0)
	rootPath := flag.Arg(1)
	foundChannel := make(chan filter.Found)
	go filter.FileInChannel(query, fget.Get(rootPath), foundChannel)
	for f := range foundChannel {
		viewer.View(f)
	}
}
