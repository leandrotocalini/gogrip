package filter

import (
	"os"
	"regexp"
	"runtime"
	"sync"
)

var workers int

func init() {
	workers = runtime.NumCPU()
}

func Filter(r *regexp.Regexp, filePath string, c chan Found) {
	file, _ := os.Open(filePath)
	found := createFound(r, file, filePath)
	found.Filter()
	if found.Match {
		c <- found
	}
	file.Close()
}

func readFileChannel(r *regexp.Regexp, filesInChan <-chan string, foundChannel chan Found, wg *sync.WaitGroup) {
	defer wg.Done()
	for fpath := range filesInChan {
		Filter(r, fpath, foundChannel)
	}
}

func FileInChannel(query string, filesInChan <-chan string, foundChannel chan Found) {
	var wg sync.WaitGroup
	r, _ := regexp.Compile(query)

	for i := 0; i <= workers; i++ {
		wg.Add(1)
		go readFileChannel(r, filesInChan, foundChannel, &wg)
	}
	wg.Wait()
	close(foundChannel)
}
