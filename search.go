package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	MAX_LINES = 12
)

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
	_, err = file.Seek(0, 0)
	if err != nil {
		return false
	}
	contentType := http.DetectContentType(buffer)
	return contentType == "text/plain; charset=utf-8"
}

func _getFiles(rootPath string, filesInChan chan string) {
	err := filepath.Walk(rootPath, func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() && isEnabled(path) {
			filesInChan <- path
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	close(filesInChan)
}

func getFilesEnabledInPath(rootPath string, buffer int) <-chan string {
	filesInChan := make(chan string, buffer)
	go _getFiles(rootPath, filesInChan)
	return filesInChan
}

func searchInFile(query, filePath string, blockChannel chan *Block) {
	f, err := os.Open(filePath)
	if err == nil {
		scanner := bufio.NewScanner(f)
		content := []string{}
		lineNumbers := []int{}
		scanner.Split(bufio.ScanLines)
		counter := 0
		for scanner.Scan() {
			text := scanner.Text()
			if strings.Contains(text, query) {
				lineNumbers = append(lineNumbers, counter)
			}
			text = strings.Replace(text, "\t", "  ", 10)
			content = append(content, text)
			counter++
		}
		for _, lineNumber := range lineNumbers {
			blockChannel <- createBlock(lineNumber, filePath, content)
		}
		f.Close()
	}
}

func getFilesAndSearch(rootPath string, buffer int, query string, blockChannel chan *Block) {
	defer close(blockChannel)
	for filePath := range getFilesEnabledInPath(rootPath, buffer*2) {
		searchInFile(query, filePath, blockChannel)
	}
}

func Search(rootPath string, buffer int, query string) <-chan *Block {
	blockChannel := make(chan *Block, buffer)
	go getFilesAndSearch(rootPath, buffer, query, blockChannel)
	return blockChannel
}
