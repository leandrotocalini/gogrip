package main

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

func SearchInFile(query string, filePath string, c chan Found, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	found := Found{Match: false, FilePath: filePath}
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
