package main

import (
	"github.com/gizak/termui/v3/widgets"
)

type SearchBox struct {
	searchText string
	widget     *widgets.Paragraph
	c          chan string
}

func (s *SearchBox) listen() {
	for text := range s.c {
		if text == "<Backspace>" && len(s.searchText) > 0 {
			s.searchText = s.searchText[:len(s.searchText)-1]
		} else if text == "<Space>" {
			s.searchText = s.searchText + " "
		} else if len(text) == 1 {
			s.searchText = s.searchText + text
		}
		s.widget.Text = s.searchText
	}
}

func (s *SearchBox) getText() string {
	return s.searchText
}

func createSearchBox() *SearchBox {
	content := widgets.NewParagraph()
	content.Text = ""
	content.Title = "Search"
	return &SearchBox{widget: content, searchText: "", c: make(chan string, 10)}
}
