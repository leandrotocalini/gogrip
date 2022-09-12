package main

import (
	"github.com/gizak/termui/v3/widgets"
)

type SearchBox struct {
	searchText string
	widget     *widgets.Paragraph
	channel    chan string
}

type SearchBoxInterface interface {
	listen()
	getText() string
	sendEvent(string)
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

func createSearchBox() *SearchBox {
	search := widgets.NewParagraph()
	search.Text = ""
	search.Title = "Search: "
	search.Border = false
	return &SearchBox{widget: search, searchText: "", channel: make(chan string, 10)}
}
