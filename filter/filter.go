package filter

import (
	"os"
	"regexp"
	"sync"
)

func Filter(r *regexp.Regexp, filePath string, c chan Found) {
	file, _ := os.Open(filePath)
	found := createFound(r, file, filePath)
	found.Filter()
	if found.Match {
		c <- found
	}
	file.Close()
}

func processFileInChannel(r *regexp.Regexp, filesInChan <-chan string, foundChannel chan Found, wg *sync.WaitGroup) {
	defer wg.Done()
	for fpath := range filesInChan {
		Filter(r, fpath, foundChannel)
	}
}

func Process(query string, filesInChan <-chan string, buffer int) <-chan Found {
	foundChannel := make(chan Found, buffer)
	go func() {
		var wg sync.WaitGroup
		r, _ := regexp.Compile(query)

		for i := 0; i <= buffer; i++ {
			wg.Add(1)
			go processFileInChannel(r, filesInChan, foundChannel, &wg)
		}
		wg.Wait()
		close(foundChannel)
	}()
	return foundChannel
}
