package main

import (
	"fmt"
	"os"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/leandrotocalini/mimaps"
)

type SearchBox struct {
	searchText string
	widget     *widgets.Paragraph
	channel    chan string
	active     bool
	cache      mimaps.InMemoryCache[string, State]
}

type SearchBoxInterface interface {
	EventManager
	getText() string
}

func (s *SearchBox) search(state State) State {
	searchText := s.getText()
	if cachedState, err := s.cache.Get(searchText); err == nil {
		cachedState.cached = true
		return cachedState
	}
	state.searchString = searchText
	state.cached = false
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
	if err := s.cache.Put(searchText, state); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
	if !s.isActive() {
		return state
	}

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
	} else {
		s.widget.Text = message

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
		cache:      mimaps.NewInMemoryCache[string, State](600),
	}
}
