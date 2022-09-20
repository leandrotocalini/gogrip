package main

import (
	"testing"

	ui "github.com/gizak/termui/v3"
	"github.com/stretchr/testify/assert"
)

func TestCreateSearchHistoryAndActive(t *testing.T) {
	searchH := createSearchHistoryBox()
	searchH.activate()
	assert.Equal(t, searchH.isActive(), true)
	assert.Equal(t, searchH.widget.BorderStyle.Fg, ui.ColorRed)
	searchH.deactivate()
	assert.Equal(t, searchH.isActive(), false)
	assert.Equal(t, searchH.widget.BorderStyle.Fg, ui.ColorWhite)
}

func TestSearchHistoryBox_add(t *testing.T) {
	searchH := createSearchHistoryBox()
	searchH.add("Batman")
	assert.Equal(t, len(searchH.widget.Rows), 2)
}

func TestSearchHistoryBox_addDuplicated(t *testing.T) {
	searchH := createSearchHistoryBox()
	searchH.add("Batman")
	searchH.add("Joker")
	searchH.add("Batman")
	assert.Equal(t, len(searchH.widget.Rows), 3)
}

func TestSearchHistoryBox_GetBoxItem(t *testing.T) {
	searchH := createSearchHistoryBox()
	item := searchH.getBoxItem()
	assert.IsType(t, ui.GridItem{}, item)
}
