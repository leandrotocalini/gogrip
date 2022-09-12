package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/leandrotocalini/gogrip/filter"
)

type ContentBox struct {
	widget  *widgets.Table
	channel chan filter.BlockInterface
}

type ContentBoxInterface interface {
	listen()
	sendEvent(filter.BlockInterface)
}

func (s *ContentBox) sendEvent(message filter.BlockInterface) {
	s.channel <- message
}

func (c *ContentBox) listen() {
	for block := range c.channel {
		c.widget.Rows = block.GetContent()
		c.widget.Title = block.GetTitle()
		c.widget.ColumnResizer()
		c.widget.RowStyles[block.GetLine()] = ui.NewStyle(ui.ColorWhite, ui.ColorRed, ui.ModifierBold)
	}
}

func createContentBox(title, text string) *ContentBox {
	content := widgets.NewTable()
	content.Title = title
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
	return &ContentBox{widget: content, channel: make(chan filter.BlockInterface)}
}
