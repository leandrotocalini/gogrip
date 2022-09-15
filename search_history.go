package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type SearchHistoryBox struct {
	widget *widgets.Table
	active bool
}

type SearchHistoryBoxInterface interface {
	WidgetInterface
	add(string)
}

func (s *SearchHistoryBox) getBoxItem() ui.GridItem {
	return ui.NewRow((1.0/15)*6, s.widget)
}

func (s *SearchHistoryBox) add(text string) {
	s.widget.Rows = append(s.widget.Rows, []string{text})
}

func (s *SearchHistoryBox) isActive() bool {
	return s.active
}

func (s *SearchHistoryBox) activate() {
	s.widget.BorderStyle.Fg = ui.ColorRed
	s.active = true
}

func (s *SearchHistoryBox) deactivate() {
	s.widget.BorderStyle.Fg = ui.ColorWhite
	s.active = false
}

func createSearchHistoryBox() *SearchHistoryBox {
	sideBar := widgets.NewTable()
	sideBar.Title = "Search History"
	sideBar.RowSeparator = false
	sideBar.Rows = [][]string{
		[]string{""},
	}
	return &SearchHistoryBox{
		widget: sideBar,
		active: false,
	}
}
