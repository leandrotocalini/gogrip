package filter

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
)

type FScanner interface {
	Scan() bool
	Text() string
}
type Block struct {
	FirstLine   int
	LastLine    int
	Content     []string
	MatchedLine int
	FilePath    string
}

const (
	MAX_LINES = 10
)

func filterFile(wg *sync.WaitGroup, query, filePath string, blockChannel chan Block) {
	wg.Add(1)
	defer wg.Done()
	f, err := os.Open(filePath)
	if err == nil {
		scanFile(bufio.NewScanner(f), filePath, query, blockChannel)
		f.Close()
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

func scanFile(scanner FScanner, filePath string, query string, blockChannel chan Block) {
	r, _ := regexp.Compile(query)
	i := 0
	content := make([]string, 0)
	lineNumbers := make([]int, 0)
	for scanner.Scan() {
		text := strings.Replace(scanner.Text(), "\t", "  ", 10)
		if r.MatchString(text) {
			lineNumbers = append(lineNumbers, i)
			text = fmt.Sprintf("[%s](fg:red,mod:bold)", text)
		}
		content = append(content, text)
		i++
	}
	for lineNumber := range lineNumbers {
		firstLine := lineNumber
		lastLine := lineNumber
		if firstLine-MAX_LINES >= 0 {
			firstLine -= MAX_LINES
		} else {
			firstLine = 0
		}
		if lastLine+MAX_LINES < len(content) {
			lastLine += MAX_LINES
		} else {
			lastLine = len(content)
		}
		blockChannel <- Block{
			FirstLine:   firstLine,
			LastLine:    lastLine,
			FilePath:    filePath,
			MatchedLine: lineNumber,
			Content:     content[firstLine:lastLine],
		}
	}
}
