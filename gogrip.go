package main

import (
	"log"
	"time"

	"runtime"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/leandrotocalini/gogrip/filter"
	"github.com/leandrotocalini/gogrip/formatter"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	fileList := widgets.NewList()
	fileList.Rows = []string{}
	fileList.Title = "Files"
	searchBox := widgets.NewParagraph()
	searchBox.Text = ""
	searchBox.Title = "Search"
	content := widgets.NewParagraph()
	content.Text = ""
	content.Title = "Content"

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1,
			ui.NewCol(1.0/4,
				ui.NewRow(1.0/15, searchBox),
				ui.NewRow((1.0/15)*14, fileList),
			),
			ui.NewCol((1.0/4)*3, content),
		),
	)

	ui.Render(grid)
	ticker := time.NewTicker(time.Second).C
	tickerCount := 0
	uiEvents := ui.PollEvents()
	blocks := []filter.Block{}
	blockPos := 0
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "<C-c>":
				return
			case "<C-f>":
				searchText := getSearchString(grid, searchBox, uiEvents)
				fileList.Title = searchText
				blocks = search(searchText, grid, fileList)
			case "<Up>":
				if blockPos > 1 {
					blockPos -= 1
					renderBlock(grid, content, blockPos, blocks)
				}
			case "<Down>":
				if blockPos < len(blocks) {
					blockPos += 1
					renderBlock(grid, content, blockPos, blocks)

				}
			}
		case <-ticker:
			tickerCount++
			ui.Render(grid)
		}
	}
}

func search(searchText string, grid *ui.Grid, fileList *widgets.List) []filter.Block {
	buffer := runtime.NumCPU()
	fileList.Rows = []string{}
	blocks := []filter.Block{}
	for block := range filter.SearchBlocks(".", buffer, searchText) {
		fileList.Rows = append(fileList.Rows, block.FilePath)
		blocks = append(blocks, block)
		ui.Render(grid)
	}
	return blocks
}

func renderBlock(grid *ui.Grid, p *widgets.Paragraph, blockPos int, blocks []filter.Block) {
	p.Text = formatter.Format(blocks[blockPos])
}

func getSearchString(grid *ui.Grid, p *widgets.Paragraph, uiEvents <-chan ui.Event) string {
	text := ""
	for i := range uiEvents {
		if i.ID == "<Enter>" {
			return text
		}
		text = text + i.ID
		p.Text = text
		ui.Render(grid)
	}
	return ""
}
