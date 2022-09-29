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
