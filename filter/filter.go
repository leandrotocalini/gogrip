package filter

import (
	"bufio"
	"os"
	"sync"
)

type Found struct {
	LineNumbers []int
	Match       bool
	FilePath    string
	Content     []string
}


type FScanner interface {
	Scan() bool
	Text() string
}


func getFilesAndFilter(rootPath string, buffer int, query string, blockChannel chan Block) {
	defer close(blockChannel)
	var wg sync.WaitGroup
	for filePath := range getFiles(rootPath, buffer*2) {
		go func() {
			wg.Add(1)
			defer wg.Done()
			f, err := os.Open(filePath)
			if err == nil {
				found := scanFile(bufio.NewScanner(f), filePath, query)
				f.Close()
				if found.Match {
					foundToBlocks(found, blockChannel)
				}
			}

		}()
	}
	wg.Wait()
}

func FilterPath(rootPath string, buffer int, query string) <-chan Block {
	blockChannel := make(chan Block, buffer)
	go getFilesAndFilter(rootPath, buffer, query, blockChannel)
	return blockChannel
}
