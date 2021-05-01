package blocks

import (
	"github.com/leandrotocalini/gogrip/filter"
	"sync"
)

type Block struct {
	FirstLine int
	LastLine  int
	Lines     map[int]bool
	Content   []string
	FilePath  string
}

func findBlockOneWay(doc []string, lineNumber, order int) int {
	found := lineNumber
	for {
		if doc[found] == "" {
			return found
		} else {
			if found > 0 && found < len(doc) -1 {
				found += order
			} else {
				return found
			}
		}
	}
	return found
}

func findBlock(doc []string, lineNumber int) (int, int) {
	return findBlockOneWay(doc, lineNumber, -1), findBlockOneWay(doc, lineNumber, 1)
}

func foundToBlocks(f filter.Found, bchan chan Block) {
	type Key struct {
		FirstLine int
		LastLine  int
	}
	blocks := make(map[Key]*Block)
	for _, lineNumber := range f.LineNumbers {
		first, last := findBlock(f.Content, lineNumber)
		key := Key{
			FirstLine: first,
			LastLine:  last,
		}
		_, ok := blocks[key]
		if !ok {
			blocks[key] = &Block{
				FirstLine: first,
				LastLine:  last,
				Lines:     make(map[int]bool),
				FilePath:  f.FilePath,
				Content:   f.Content[first:last]}
		}
		blocks[key].Lines[lineNumber] = true
	}
	for _, val := range blocks {
		bchan <- *val
	}
}

func Process(foundChan <-chan filter.Found, buffer int) <-chan Block {
	bchan := make(chan Block, buffer)
	go func() {
		var wg sync.WaitGroup
		for i := 0; i <= buffer; i++ {
			wg.Add(1)
			go func(foundChan <-chan filter.Found, bchan chan Block, wg *sync.WaitGroup) {
				defer wg.Done()
				for f := range foundChan {
					foundToBlocks(f, bchan)
				}
			}(foundChan, bchan, &wg)
		}
		wg.Wait()
		close(bchan)
	}()
	return bchan
}
