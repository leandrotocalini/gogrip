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
	EventManager
}

func (c *ContentBox) listen() {
	for state := range c.channel {
		if state.total > 0 {
			c.widget.Rows = state.currentBlock.GetContent()
			c.widget.Title = state.currentBlock.GetTitle()
			for k := range c.widget.RowStyles {
				delete(c.widget.RowStyles, k)
			}
			c.widget.RowStyles[state.currentBlock.HighlightCursorLine()] = ui.NewStyle(ui.ColorYellow)
			if line, err := state.currentBlock.HighlightMatchedLine(); err == nil {
				c.widget.RowStyles[line] = ui.NewStyle(ui.ColorRed)

			}
		}
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

func (c *ContentBox) expose(state State) {
	c.channel <- state
}

func (c *ContentBox) newEvent(state State, message string) State {
	if c.isActive() {
		switch message {
		case "<Left>":
			state.previousBlock()
		case "<Right>":
			state.nextBlock()
		case "<Up>":
			state.currentBlock.FocusUp()
		case "<Down>":
			state.currentBlock.FocusDown()
		}
	}
	return state
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
