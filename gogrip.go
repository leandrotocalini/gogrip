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
	foundChannel := make(chan Found)
	go func() {
		var wg sync.WaitGroup
		for fpath := range filesInChan {
			wg.Add(1)
			go SearchInFile(query, fpath, foundChannel, &wg)
		}
		wg.Wait()
		close(foundChannel)
	}()
	for elem := range foundChannel {
		View(elem)
	}
}
