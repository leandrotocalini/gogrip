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
	assert.Equal(t, 0, len(stateExposer.widget.Rows))
}

func TestStateExposer_updateWidgetWithRows(t *testing.T) {
	stateExposer := createStateExposer()
	block1 := createBlock(2, "asd", []string{"asd", "22", "333"})
	block2 := createBlock(2, "asd", []string{"asd", "22", "333"})
	state := State{
		position:     0,
		total:        2,
		searchString: "",
		blocks:       []BlockInterface{block1, block2},
		currentBlock: block1,
	}
	stateExposer.updateWidget(state)
	assert.Equal(t, 8, len(stateExposer.widget.Rows))
}

func TestStateExposer_updateWidgetWithRowsInTheMiddle(t *testing.T) {
	stateExposer := createStateExposer()
	block1 := createBlock(2, "asd", []string{"asd", "22", "333"})
	block2 := createBlock(2, "asd", []string{"asd", "22", "333"})
	block3 := createBlock(2, "asd", []string{"asd", "22", "333"})
	state := State{
		position:     1,
		total:        3,
		searchString: "",
		blocks:       []BlockInterface{block1, block2, block3},
		currentBlock: block2,
	}
	stateExposer.updateWidget(state)
	assert.Equal(t, 8, len(stateExposer.widget.Rows))
}

func TestStateExposer_updateWidgetWithRowsInTheEnd(t *testing.T) {
	stateExposer := createStateExposer()
	block1 := createBlock(2, "asd", []string{"asd", "22", "333"})
	block2 := createBlock(2, "asd", []string{"asd", "22", "333"})
	block3 := createBlock(2, "asd", []string{"asd", "22", "333"})
	state := State{
		position:     2,
		total:        3,
		searchString: "",
		blocks:       []BlockInterface{block1, block2, block3},
		currentBlock: block2,
	}
	stateExposer.updateWidget(state)
	assert.Equal(t, 8, len(stateExposer.widget.Rows))
}
