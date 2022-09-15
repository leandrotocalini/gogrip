package main

import "fmt"

type Block struct {
	Line     int
	Content  []string
	FilePath string
	Position int
}

type BlockInterface interface {
	GetContent() [][]string
	GetTitle() string
	GetLine() int
	SetPosition(int)
	GetPosition() int
	GetMatchedWord() string
}

func (block *Block) GetTitle() string {
	return block.FilePath
}

func (block *Block) GetMatchedWord() string {
	return block.Content[block.Line]
}

func (block *Block) GetPosition() int {
	return block.Position
}

func (block *Block) SetPosition(pos int) {
	block.Position = pos
}

func (block *Block) GetContent() [][]string {
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

func (block *Block) GetLine() int {
	return block.Line
}

func createBlock(lineNumber int, filePath string, content []string) *Block {
	return &Block{
		Line:     lineNumber,
		FilePath: filePath,
		Content:  content,
	}
}
