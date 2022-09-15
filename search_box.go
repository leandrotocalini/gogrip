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
	WidgetInterface
	listen()
	getText() string
	sendEvent(string)
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

func (s *SearchBox) listen() {
	for text := range s.channel {
		if text == "<Backspace>" && len(s.searchText) > 0 {
			s.searchText = s.searchText[:len(s.searchText)-1]
			s.widget.Text = s.searchText

		} else if text == "<Space>" {
			s.searchText = s.searchText + " "
			s.widget.Text = s.searchText
		} else if len(text) == 1 {
			s.searchText = s.searchText + text
			s.widget.Text = s.searchText
		}
	}
}

func (s *SearchBox) getText() string {
	return s.searchText
}

func (s *SearchBox) sendEvent(message string) {
	s.channel <- message
}

func (s *SearchBox) getBoxItem() ui.GridItem {
	return ui.NewRow(1.0/15, s.widget)
}

func (s *SearchBox) update(state State) {
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
