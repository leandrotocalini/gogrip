package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Found struct {
	LineNumbers []int
	Match       bool
	FilePath    string
	Content     []string
}

func findEndOfParraph(doc []string, lineNumber int) (int, int){
	ff := false
	fl := false
	first := lineNumber
	last := lineNumber
	for !ff || !fl{
		if ff == false {
			if doc[first] == "" {
				ff = true
			}else {
				if first > 0 {
					first--
				}
			}
		}
		if fl == false {
			if doc[last] == "" {
				fl = true
			}else {
				if last < len(doc) {
					last++
				}
			}
		}
	}
	return first, last
}

func PrintLine(doc []string, lineNumber int) {
	first, last := findEndOfParraph(doc, lineNumber)
	normal := color.New(color.FgWhite).SprintFunc()
	highlight := color.New(color.FgRed).SprintFunc()
	for i := first; i <= last; i++ {
		if i == lineNumber {
			fmt.Printf("%s - %s \n", normal(i), highlight(doc[i]))
		} else {
			fmt.Printf("%s - %s \n", normal(i), normal(doc[i]))
		}
	}
}

func (f *Found) Print() {
	fmt.Println("#########################")
	fmt.Println(f.FilePath)
	h := make(map[string][]int)
	for _, val := range f.LineNumbers {
		first, last := findEndOfParraph(f.Content, val)
		key := fmt.Sprint(first, "-", last)
		h[key] = append(h[key], val)
	}
	fmt.Println(h)

}

func SearchInFile(query string, filePath string, c chan Found, wg *sync.WaitGroup) {

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
	wg.Done()
}

func GetFiles(rootPath string) <-chan string {
	out := make(chan string)
	go func() {
		filepath.Walk(rootPath, func(path string, file os.FileInfo, err error) error {
			if !file.IsDir() {
				out <- path
			}
			return nil
		})
		close(out)
	}()
	return out
}

func SearchFilesInChannel(filesInChan <-chan string, query string) <-chan Found {
	out := make(chan Found)
	go func() {
		var wg sync.WaitGroup
		for fpath := range filesInChan {
			wg.Add(1)
			go SearchInFile(query, fpath, out, &wg)
		}
		wg.Wait()
		close(out)
	}()
	return out
}

func Search(query string, rootPath string) {
	filesInChan := GetFiles(rootPath)
	c := SearchFilesInChannel(filesInChan, query)
	for elem := range c {
		elem.Print()
	}
}

func main() {
	flag.Parse()
	query := flag.Arg(0)
	rootPath := flag.Arg(1)
	Search(query, rootPath)

}
