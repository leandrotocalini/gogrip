package filter

import (
	"bufio"
	"os"
	"sync"
)




func filterFile(wg *sync.WaitGroup, query, filePath string, blockChannel chan Block) {
	wg.Add(1)
	defer wg.Done()
	f, err := os.Open(filePath)
	if err == nil {
		found := scanFile(bufio.NewScanner(f), filePath, query)
		f.Close()
		if found.Match {
			for _, val := range makeBlocks(found) {
				blockChannel <- *val
			}
		}
	}
}

func searchInFiles(rootPath string, buffer int, query string, blockChannel chan Block) {
	defer close(blockChannel)
	var wg sync.WaitGroup
	for filePath := range getFiles(rootPath, buffer*2) {
		go filterFile(&wg, query, filePath, blockChannel)
	}
	wg.Wait()
}

func SearchBlocks(rootPath string, buffer int, query string) <-chan Block {
	blockChannel := make(chan Block, buffer)
	go searchInFiles(rootPath, buffer, query, blockChannel)
	return blockChannel
}
