package filter

import (
	"bufio"
	"github.com/leandrotocalini/gogrip/ds"
	"os"
	"strings"
	"sync"
)

func Filter(query string, filePath string, c chan ds.Found) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	found := ds.Found{Match: false, FilePath: filePath}
	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		line := string(scanner.Text())
		found.Content = append(found.Content, line)
		if strings.Contains(line, query) {
			found.Match = true
			found.LineNumbers = append(found.LineNumbers, i)
		}
	}
	file.Close()
	if found.Match {
		c <- found
	}
}

func readFileChannel(query string, filesInChan <-chan string, foundChannel chan ds.Found, wg *sync.WaitGroup) {
	defer wg.Done()
	for fpath := range filesInChan {
		Filter(query, fpath, foundChannel)
	}
}

func FilterFileIn(query string, filesInChan <-chan string, foundChannel chan ds.Found) {
	var wg sync.WaitGroup
	for i := 0; i <= 5; i++ {
		wg.Add(1)
		go readFileChannel(query, filesInChan, foundChannel, &wg)
	}
	wg.Wait()
	close(foundChannel)
}