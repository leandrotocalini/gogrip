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
}


type FScanner interface {
	Scan() bool
	Text() string
}

func filterScanner(scanner FScanner, r *regexp.Regexp) ([]string, []int) {
	i := 0
	content := make([]string, 0)
	lineNumbers := make([]int, 0)
	for scanner.Scan() {
		text := scanner.Text()
		content = append(content, text)
		if r.MatchString(text) {
			lineNumbers = append(lineNumbers, i)
		}
		i++
	}

	return content, lineNumbers
}

func scanFile(scanner FScanner, filePath string, r *regexp.Regexp) Found {
	f := Found{Match: false, FilePath: filePath}
	content, lineNumbers := filterScanner(scanner, r)
	if len(lineNumbers) > 0 {
		f.Content = content
		f.LineNumbers = lineNumbers
		f.Match = true
	}
	return f
}

func getFilesAndFilter(rootPath string, buffer int, query string, blockChannel chan Block) {
	r, _ := regexp.Compile(query)
	defer close(blockChannel)
	var wg sync.WaitGroup
	for filePath := range getFiles(rootPath, buffer*2) {
		go func() {
			wg.Add(1)
			defer wg.Done()
			f, err := os.Open(filePath)
			if err == nil {
				found := scanFile(bufio.NewScanner(f), filePath, r)
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
