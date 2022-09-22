package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateExposer_updateWidgetEmpty(t *testing.T) {

	stateExposer := createStateExposer()
	state := State{
		position:     0,
		total:        0,
		searchString: "",
		blocks:       []BlockInterface{},
	}
	stateExposer.updateWidget(state)
	assert.Equal(t, 1, len(stateExposer.widget.Rows))
}

func TestStateExposer_updateWidgetWithRows(t *testing.T) {
	stateExposer := createStateExposer()
	block1 := createBlock(2, "asd", []string{"asd"})
	block2 := createBlock(2, "asd3", []string{"asd3"})
	state := State{
		position:     0,
		total:        2,
		searchString: "",
		blocks:       []BlockInterface{block1, block2},
		currentBlock: block1,
	}
	stateExposer.updateWidget(state)
	assert.Equal(t, 6, len(stateExposer.widget.Rows))
}
