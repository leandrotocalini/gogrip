package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type ContentBox struct {
	widget  *widgets.Table
	channel chan State
	active  bool
}

type ContentBoxInterface interface {
	WidgetInterface
	listen()
}

func (c *ContentBox) listen() {
	for state := range c.channel {
		c.widget.Rows = state.currentBlock.GetContent()
		c.widget.Title = state.currentBlock.GetTitle()
		c.widget.ColumnResizer()
		c.widget.RowStyles[state.currentBlock.GetLine()] = ui.NewStyle(ui.ColorWhite, ui.ColorRed, ui.ModifierBold)
	}
}

func (s *ContentBox) isActive() bool {
	return s.active
}

func (s *ContentBox) activate() {
	s.widget.BorderStyle.Fg = ui.ColorRed
	s.active = true
}

func (s *ContentBox) deactivate() {
	s.widget.BorderStyle.Fg = ui.ColorWhite
	s.active = false
}
func (c *ContentBox) getBoxItem() ui.GridItem {
	return ui.NewRow((1.0/15)*13.5, c.widget)
}

func (c *ContentBox) update(state State) {
	c.channel <- state
}

func createContentBox() *ContentBox {
	content := widgets.NewTable()
	content.Title = "Filename"
	content.TextStyle = ui.NewStyle(ui.ColorWhite)
	content.Rows = [][]string{
		[]string{"", ""},
	}
	content.BorderStyle = ui.NewStyle(ui.ColorWhite)
	content.FillRow = false
	content.ColumnWidths = []int{5, 400}
	//content.TextAlignment = ui.AlignCenter
	content.RowSeparator = false
	//content.BorderStyle.Fg = ui.ColorRed
	return &ContentBox{
		widget:  content,
		channel: make(chan State),
		active:  false,
	}
}
