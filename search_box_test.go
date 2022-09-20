package main

import (
	"testing"

	ui "github.com/gizak/termui/v3"
	"github.com/stretchr/testify/assert"
)

func Test_createSearchBoxGetBoxItem(t *testing.T) {
	searchH := createSearchBox()
	item := searchH.getBoxItem()
	assert.IsType(t, ui.GridItem{}, item)
}

func Test_SearchBoxActive(t *testing.T) {
	searchH := createSearchBox()
	searchH.activate()
	assert.Equal(t, searchH.isActive(), true)
	assert.Equal(t, searchH.widget.BorderStyle.Fg, ui.ColorRed)
	searchH.deactivate()
	assert.Equal(t, searchH.isActive(), false)
	assert.Equal(t, searchH.widget.BorderStyle.Fg, ui.ColorWhite)
}

func Test_SearchBoxListenText(t *testing.T) {
	searchH := createSearchBox()
	searchH.listenText("b")
	searchH.listenText("a")
	searchH.listenText("t")
	assert.Equal(t, searchH.getText(), "bat")
}

func Test_SearchBoxListenBackSpace(t *testing.T) {
	searchH := createSearchBox()
	searchH.listenText("b")
	searchH.listenText("<Backspace>")
	searchH.listenText("t")
	assert.Equal(t, searchH.getText(), "t")
}

func Test_SearchBoxListenTextSpace(t *testing.T) {
	searchH := createSearchBox()
	searchH.listenText("b")
	searchH.listenText("<Space>")
	searchH.listenText("t")
	assert.Equal(t, searchH.getText(), "b t")
}
