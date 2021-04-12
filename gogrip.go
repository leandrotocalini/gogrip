package main

import (
	"flag"
	"os"
	"path/filepath"
	"sync"
)

type Found struct {
	LineNumbers []int
	Match       bool
	FilePath    string
	Content     []string
}

func GetFiles(rootPath string) <-chan string {
	out := make(chan string)
	go func() {
		filepath.Walk(rootPath, func(path string, file os.FileInfo, err error) error {
			if !file.IsDir() {
				out <- path
			}
			return nil
		})
		close(out)
	}()
	return out
}

func SearchFilesInChannel(filesInChan <-chan string, query string) <-chan Found {
	out := make(chan Found)
	go func() {
		var wg sync.WaitGroup
		for fpath := range filesInChan {
			wg.Add(1)
			go SearchInFile(query, fpath, out, &wg)
		}
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	flag.Parse()
	query := flag.Arg(0)
	rootPath := flag.Arg(1)
	filesInChan := GetFiles(rootPath)
	c := SearchFilesInChannel(filesInChan, query)
	for elem := range c {
		View(elem)
	}
}
