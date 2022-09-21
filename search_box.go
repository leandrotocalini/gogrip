package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type SearchBox struct {
	searchText string
	widget     *widgets.Paragraph
	channel    chan string
	active     bool
}

type SearchBoxInterface interface {
	EventManagerWidget
	getText() string
}

func (s *SearchBox) search(state State) State {
	searchText := s.getText()
	state.searchString = searchText
	buffer := 2
	state.blocks = []BlockInterface{}
	state.position = 0
	for block := range Search(".", buffer, searchText) {
		state.blocks = append(state.blocks, block)
		if len(state.blocks) == 1 {
			state.total = 1
			state.moveToBlock(0)
		}
		block.SetPosition(len(state.blocks) - 1)
	}
	state.total = len(state.blocks)
	state.moveToBlock(0)
	return state
}

func (s *SearchBox) isActive() bool {
	return s.active
}

func (s *SearchBox) activate() {
	s.widget.BorderStyle.Fg = ui.ColorRed
	s.active = true
}

func (s *SearchBox) deactivate() {
	s.widget.BorderStyle.Fg = ui.ColorWhite
	s.active = false
}

// notest
func (s *SearchBox) listen() {
	//
}

func (s *SearchBox) getText() string {
	return s.searchText
}

func (s *SearchBox) newEvent(state State, message string) State {
	if message == "<Backspace>" && len(s.searchText) > 0 {
		s.searchText = s.searchText[:len(s.searchText)-1]
		s.widget.Text = s.searchText

	} else if message == "<Space>" {
		s.searchText = s.searchText + " "
		s.widget.Text = s.searchText
	} else if message == "<Enter>" {
		state.searchString = s.searchText
		state = s.search(state)
	} else if len(message) == 1 {
		s.searchText = s.searchText + message
		s.widget.Text = s.searchText
	}
	return state
}

func (s *SearchBox) getBoxItem() ui.GridItem {
	return ui.NewRow(1.0/15, s.widget)
}

func (s *SearchBox) expose(state State) {
	// pass
}

func createSearchBox() *SearchBox {
	search := widgets.NewParagraph()
	search.Text = ""
	search.Title = "Search: "
	return &SearchBox{
		widget:     search,
		searchText: "",
		channel:    make(chan string, 10),
		active:     false,
	}
}
