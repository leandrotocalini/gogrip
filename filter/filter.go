package filter

import (
	"bufio"
	"github.com/leandrotocalini/gogrip/ds"
	"os"
	"regexp"
	"runtime"
	"sync"
)

var workers int

func init() {
	workers = runtime.NumCPU()
}

func FilterFile(r *regexp.Regexp, filePath string, scanner *bufio.Scanner) ds.Found {
	found := ds.Found{Match: false, FilePath: filePath, Regexp: r}
	for i := 0; scanner.Scan(); i++ {
		line := string(scanner.Text())
		found.Content = append(found.Content, line)
		if r.MatchString(line) {
			found.Match = true
			found.LineNumbers = append(found.LineNumbers, i)
		}
	}
	return found
}

func Filter(r *regexp.Regexp, filePath string, c chan ds.Found) {
	file, _ := os.Open(filePath)
	scanner := bufio.NewScanner(file)
	found := FilterFile(r, filePath, scanner)
	if found.Match {
		c <- found
	}
	file.Close()

}

func readFileChannel(r *regexp.Regexp, filesInChan <-chan string, foundChannel chan ds.Found, wg *sync.WaitGroup) {
	defer wg.Done()
	for fpath := range filesInChan {
		Filter(r, fpath, foundChannel)
	}
}

func FilterFileIn(query string, filesInChan <-chan string, foundChannel chan ds.Found) {
	var wg sync.WaitGroup
	r, _ := regexp.Compile(query)

	for i := 0; i <= workers; i++ {
		wg.Add(1)
		go readFileChannel(r, filesInChan, foundChannel, &wg)
	}
	wg.Wait()
	close(foundChannel)
}
