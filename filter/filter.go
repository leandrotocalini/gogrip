package filter

import (
	"bufio"
	"os"
	"io"
	"regexp"
	"sync"
)

type Found struct {
	LineNumbers []int
	Match       bool
	FilePath    string
	Content     []string
}

type Line struct{
	Text	string
	Ln		int
}


func filterScanChannel(r *regexp.Regexp, filePath string, c chan Line) Found {
	f := Found{Match: false, FilePath: filePath}
	for line := range c{
		f.Content = append(f.Content, line.Text)
		if r.MatchString(line.Text) {
			f.Match = true
			f.LineNumbers = append(f.LineNumbers, line.Ln)
		}
	}
	return f
}

func scanFile(handle io.Reader, c chan Line) {
	scanner := bufio.NewScanner(handle)
	i := 0
	for scanner.Scan() {
		text := string(scanner.Text())
		c <- Line{Text: text, Ln: i}
		i++	
	}
}

func openScanFile(filePath string, c chan Line){
	f, _ := os.Open(filePath)
	defer f.Close()
	defer close(c)
	scanFile(f, c)
}


func filterPath(r *regexp.Regexp, filePath string) Found{
	c := make(chan Line, 10)
	go openScanFile(filePath, c)
	return filterScanChannel(r, filePath, c)
}

func filterWorker(r *regexp.Regexp, filesInChan <-chan string, foundChan chan Found, wg *sync.WaitGroup) {
	defer wg.Done()
	for fpath := range filesInChan {
		f := filterPath(r, fpath)
		if f.Match {
			foundChan <- f
		}
	}
}

func filterPool(filesInChan <-chan string, foundChan chan Found, query string, buffer int){
	defer close(foundChan)
	var wg sync.WaitGroup
	r, _ := regexp.Compile(query)
	for i := 0; i <= buffer; i++ {
		wg.Add(1)
		go filterWorker(r, filesInChan, foundChan, &wg)
	}
	wg.Wait()
}

func Process(query string, filesInChan <-chan string, buffer int) <-chan Found {
	foundChan := make(chan Found, buffer)
	go filterPool(filesInChan, foundChan, query, buffer)
	return foundChan
}
