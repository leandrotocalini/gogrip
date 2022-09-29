package main

import "fmt"

type Block struct {
	MatchedLine int
	CursorLine  int
	firstLine   int
	lastLine    int
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
	HighlightCursorLine() int
	HighlightMatchedLine() (int, error)
}

func (block *Block) GetTitle() string {
	return block.FilePath
}

func (block *Block) GetMatchedWord() string {
	return block.Content[block.MatchedLine]
}

func (block *Block) GetPosition() int {
	return block.Position
}

func (block *Block) SetPosition(pos int) {
	block.Position = pos
}

func (block *Block) FocusUp() {
	if block.CursorLine > 0 {
		block.CursorLine -= 1
	}
}

func (block *Block) FocusDown() {
	if block.CursorLine < len(block.Content) {
		block.CursorLine += 1
	}
}

func (block *Block) setLines() {
	firstLine := block.CursorLine
	lastLine := firstLine + 1
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
	block.firstLine = firstLine
	block.lastLine = lastLine
}

func (block *Block) HighlightCursorLine() int {
	return block.CursorLine - block.firstLine + 1
}

func (block *Block) HighlightMatchedLine() (int, error) {
	l := block.MatchedLine - block.firstLine + 1
	if l >= 0 {
		return l, nil
	}
	return l, nil
}

func (block *Block) GetContent() [][]string {
	block.setLines()
	rows := [][]string{
		[]string{"", ""},
	}
	for i := block.firstLine; i < block.lastLine; i++ {
		lineCounter := fmt.Sprintf("%d", i)
		line := []string{lineCounter, block.Content[i]}

		rows = append(rows, line)
	}
	return rows
}

func (block *Block) GetLine() int {
	return block.MatchedLine
}

func createBlock(lineNumber int, filePath string, content []string) *Block {
	return &Block{
		MatchedLine: lineNumber,
		CursorLine:  lineNumber,
		FilePath:    filePath,
		Content:     content,
	}
}
