package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createBlock(t *testing.T) {
	block := createBlock(2, "asd", []string{"asd", "22", "333"})
	block.Position = 1
	assert.Equal(t, 2, block.GetLine())
	assert.Equal(t, 1, block.GetPosition())
	assert.Equal(t, 4, len(block.GetContent()))
}

func Test_createBlockMoreContent(t *testing.T) {
	content := []string{}
	for i := 0; i < 100; i++ {
		content = append(content, fmt.Sprintf("content line: %d>>>", i))
	}
	block := createBlock(2, "asd", content)
	block.Position = 1
	assert.Equal(t, 2, block.GetLine())
	assert.Equal(t, 1, block.GetPosition())
	assert.Equal(t, 41, len(block.GetContent()))
}

func Test_createBlockMoreContentMidle(t *testing.T) {
	content := []string{}
	for i := 0; i < 100; i++ {
		content = append(content, fmt.Sprintf("content line: %d>>>", i))
	}
	block := createBlock(50, "asd", content)
	block.Position = 1
	assert.Equal(t, 50, block.GetLine())
	assert.Equal(t, 1, block.GetPosition())
	assert.Equal(t, 41, len(block.GetContent()))
}

func Test_createBlockFocusUp(t *testing.T) {
	content := []string{}
	for i := 0; i < 100; i++ {
		content = append(content, fmt.Sprintf("content line: %d>>>", i))
	}
	block := createBlock(10, "asd", content)
	block.Position = 1
	block.FocusUp()
	assert.Equal(t, 41, len(block.GetContent()))
	assert.Equal(t, 10, block.HighlightCursorLine())
	val, err := block.HighlightMatchedLine()
	assert.Equal(t, err, nil)
	assert.Equal(t, 11, val)
}

func Test_createBlockFocusDown(t *testing.T) {
	content := []string{}
	for i := 0; i < 100; i++ {
		content = append(content, fmt.Sprintf("content line: %d>>>", i))
	}
	block := createBlock(10, "asd", content)
	block.Position = 1
	block.FocusDown()
	block.FocusDown()
	block.FocusDown()
	assert.Equal(t, 41, len(block.GetContent()))
	assert.Equal(t, 14, block.HighlightCursorLine())
	val, err := block.HighlightMatchedLine()
	assert.Equal(t, err, nil)
	assert.Equal(t, 11, val)
}

func Test_createBlockFocusDownLine0(t *testing.T) {
	content := []string{}
	for i := 0; i < 100; i++ {
		content = append(content, fmt.Sprintf("content line: %d>>>", i))
	}
	block := createBlock(0, "asd", content)
	block.Position = 1
	for i := 0; i < 30; i++ {
		block.FocusDown()
	}

	assert.Equal(t, 41, len(block.GetContent()))
	assert.Equal(t, 21, block.HighlightCursorLine())
	val, err := block.HighlightMatchedLine()
	assert.NotEqual(t, err, nil)
	assert.Equal(t, -9, val)
}

func Test_createBlockFocusDownEndOfFile(t *testing.T) {
	content := []string{}
	for i := 0; i < 100; i++ {
		content = append(content, fmt.Sprintf("content line: %d>>>", i))
	}
	block := createBlock(0, "asd", content)
	block.Position = 1
	for i := 0; i < 111; i++ {
		block.FocusDown()
	}

	assert.Equal(t, 41, len(block.GetContent()))
	assert.Equal(t, 41, block.HighlightCursorLine())
	val, err := block.HighlightMatchedLine()
	assert.NotEqual(t, err, nil)
	assert.Equal(t, -59, val)
}
