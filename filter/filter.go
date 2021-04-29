package filter

import (
	"bufio"
	"os"
	"regexp"
	"sync"
)

type Found struct {
	LineNumbers []int
	Match       bool
	FilePath    string
	Content     []string
	Regexp      *regexp.Regexp
	Scanner     *bufio.Scanner
}

func (f *Found) Filter() {
	for i := 0; f.Scanner.Scan(); i++ {
		line := string(f.Scanner.Text())
		f.Content = append(f.Content, line)
		if f.Regexp.MatchString(line) {
			f.Match = true
			f.LineNumbers = append(f.LineNumbers, i)
		}
	}
}

func createFound(r *regexp.Regexp, file *os.File, filePath string) Found {
	scanner := bufio.NewScanner(file)
	found := Found{Match: false, FilePath: filePath, Regexp: r, Scanner: scanner}
	return found
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

func Process(query string, filesInChan <-chan string, buffer int) <-chan Found {
	foundChannel := make(chan Found, buffer)
	go func() {
		var wg sync.WaitGroup
		r, _ := regexp.Compile(query)
		for i := 0; i <= buffer; i++ {
			wg.Add(1)
			go func(r *regexp.Regexp, filesInChan <-chan string, foundChannel chan Found, wg *sync.WaitGroup) {
				defer wg.Done()
				for fpath := range filesInChan {
					Filter(r, fpath, foundChannel)
				}
			}(r, filesInChan, foundChannel, &wg)
		}
		wg.Wait()
		close(foundChannel)
	}()
	return foundChannel
}
