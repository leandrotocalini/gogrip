package filter

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Block struct {
	Line     int
	Content  []string
	FilePath string
}

type BlockInterface interface {
	GetContent() [][]string
	GetTitle() string
	GetLine() int
}

const (
	MAX_LINES = 7
)

func (block Block) GetTitle() string {
	return block.FilePath
}

func (block Block) GetContent() [][]string {
	firstLine := block.Line
	lastLine := block.Line + 1
	contentLen := len(block.Content)
	if firstLine-MAX_LINES >= 0 {
		firstLine -= MAX_LINES
	} else {
		firstLine = 0
	}

	if lastLine+MAX_LINES > contentLen {
		lastLine = contentLen
	} else {
		lastLine += MAX_LINES - 1
	}
	rows := [][]string{
		[]string{"", ""},
	}
	for i := firstLine; i < lastLine; i++ {
		line := []string{fmt.Sprintf("%d", i), block.Content[i]}
		rows = append(rows, line)
	}
	return rows
}

func (block Block) GetLine() int {
	return block.Line
}
func isEnabled(path string) bool {
	if strings.Contains(path, ".git") {
		return false
	}
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	buffer := make([]byte, 16)
	_, err = file.Read(buffer)
	if err != nil {
		return false
	}
	file.Seek(0, 0)
	contentType := http.DetectContentType(buffer)
	return contentType == "text/plain; charset=utf-8"
}

func _getFiles(rootPath string, filesInChan chan string) {
	filepath.Walk(rootPath, func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() && isEnabled(path) {
			filesInChan <- path
		}
		return nil
	})
	close(filesInChan)
}

func getFilesEnabledInPath(rootPath string, buffer int) <-chan string {
	filesInChan := make(chan string, buffer)
	go _getFiles(rootPath, filesInChan)
	return filesInChan
}
func searchInFile(wg *sync.WaitGroup, query, filePath string, blockChannel chan Block) {
	wg.Add(1)
	defer wg.Done()
	f, err := os.Open(filePath)
	if err == nil {
		searchBlocksInFile(bufio.NewScanner(f), filePath, query, blockChannel)
		f.Close()
	}
}

func getFilesAndSearch(rootPath string, buffer int, query string, blockChannel chan Block) {
	defer close(blockChannel)
	var wg sync.WaitGroup
	for filePath := range getFilesEnabledInPath(rootPath, buffer*2) {
		go searchInFile(&wg, query, filePath, blockChannel)
	}
	wg.Wait()
}

func Search(rootPath string, buffer int, query string) <-chan Block {
	blockChannel := make(chan Block, buffer)
	go getFilesAndSearch(rootPath, buffer, query, blockChannel)
	return blockChannel
}

func searchBlocksInFile(scanner *bufio.Scanner, filePath string, query string, blockChannel chan Block) {
	content := []string{}
	lineNumbers := []int{}
	scanner.Split(bufio.ScanLines)
	counter := 0
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.Replace(text, "\t", "  ", 10)
		if strings.Contains(text, query) {
			lineNumbers = append(lineNumbers, counter)
		}

		content = append(content, text)
		counter++
	}
	for _, lineNumber := range lineNumbers {
		blockChannel <- Block{
			Line:     lineNumber,
			FilePath: filePath,
			Content:  content,
		}
	}
}
