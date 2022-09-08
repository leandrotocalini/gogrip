package main

import (
	"strings"

	"github.com/gizak/termui/v3/widgets"
	"github.com/leandrotocalini/gogrip/filter"
)

type ContentBox struct {
	widget      *widgets.Paragraph
	contentChan chan filter.Block
}

func (c *ContentBox) listen() {
	for block := range c.contentChan {
		c.widget.Text = strings.Join(block.Content[:], "\n")
		c.widget.Title = block.FilePath
	}
}

func createContentBox(title, text string) *ContentBox {
	content := widgets.NewParagraph()
	content.Text = text
	content.Title = title
	//content.BorderStyle.Fg = ui.ColorRed
	return &ContentBox{widget: content, contentChan: make(chan filter.Block)}
}
