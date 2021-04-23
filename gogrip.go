package main

import (
	"flag"
	"github.com/leandrotocalini/gogrip/ds"
	"github.com/leandrotocalini/gogrip/filter"
	"github.com/leandrotocalini/gogrip/formatter"
	"os"
	"path/filepath"
)

func GetFiles(rootPath string) <-chan string {
	filesInChan := make(chan string)
	go func() {
		filepath.Walk(rootPath, func(path string, file os.FileInfo, err error) error {
			if !file.IsDir() {
				filesInChan <- path
			}
			return nil
		})
		close(filesInChan)
	}()
	return filesInChan
}

func main() {
	flag.Parse()
	query := flag.Arg(0)
	rootPath := flag.Arg(1)
	filesInChan := GetFiles(rootPath)
	foundChannel := make(chan ds.Found)
	go filter.FilterFileIn(query, filesInChan, foundChannel)
	for elem := range foundChannel {
		formatter.View(elem)
	}
}
