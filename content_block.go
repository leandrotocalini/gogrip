package main

import "fmt"

type Block struct {
	Line        int
	FocusedLine int
	Content     []string
	FilePath    string
	Position    int
}

type BlockInterface interface {
	GetContent() [][]string
	GetTitle() string
	GetLine() int
	SetPosition(int)
	GetPosition() int
	GetMatchedWord() string
	FocusUp()
	FocusDown()
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
	return block.GetContentFromLine(block.FocusedLine)
}

func (block *Block) FocusUp() {
	if block.FocusedLine > 0 {
		block.FocusedLine -= 1
	}
}

func (block *Block) FocusDown() {
	if block.FocusedLine < len(block.Content) {
		block.FocusedLine += 1
	}
}

func (block *Block) GetContentFromLine(line int) [][]string {
	firstLine := line
	lastLine := line + 1
	contentLen := len(block.Content)
	if firstLine-MAX_LINES >= 0 {
		firstLine -= MAX_LINES
	} else {
		firstLine = 0
	}
	if firstLine == 0 {
		lastLine = MAX_LINES * 2
		if lastLine > contentLen {
			lastLine = contentLen
		}
	} else {
		if lastLine+MAX_LINES > contentLen {
			lastLine = contentLen
			if lastLine-(MAX_LINES*2) >= 0 {
				firstLine = lastLine - (MAX_LINES * 2)
			}
		} else {
			lastLine += MAX_LINES - 1
		}
	}

	rows := [][]string{
		[]string{"", ""},
	}
	for i := firstLine; i < lastLine; i++ {
		lineCounter := fmt.Sprintf("%d", i)
		if i == block.Line {
			lineCounter = fmt.Sprintf("%d>>>", i)

		}
		line := []string{lineCounter, block.Content[i]}

		rows = append(rows, line)
	}
	return rows
}

func (block *Block) GetLine() int {
	return block.Line
}

func createBlock(lineNumber int, filePath string, content []string) *Block {
	return &Block{
		Line:        lineNumber,
		FocusedLine: lineNumber,
		FilePath:    filePath,
		Content:     content,
	}
}
