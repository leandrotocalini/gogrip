package main

import (
	"github.com/gizak/termui/v3/widgets"
	"github.com/leandrotocalini/gogrip/filter"
	"github.com/leandrotocalini/gogrip/formatter"
)

type ContentBox struct {
	widget      *widgets.Paragraph
	contentChan chan filter.Block
}

func (c *ContentBox) listen() {
	for block := range c.contentChan {
		c.widget.Text = formatter.Format(block)
		c.widget.Title = block.FilePath
	}
}

func createContentBox(title, text string) *ContentBox {
	content := widgets.NewParagraph()
	content.Text = text
	content.Title = title
	return &ContentBox{widget: content, contentChan: make(chan filter.Block)}
}
