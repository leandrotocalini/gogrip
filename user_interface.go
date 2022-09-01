package main

import (
	"runtime"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/leandrotocalini/gogrip/filter"
	"github.com/leandrotocalini/gogrip/formatter"
)

type UserInterface struct {
	fileList   *widgets.List
	searchText string
	searchBox  *widgets.Paragraph
	content    *widgets.Paragraph
	grid       *ui.Grid
	blocks     []filter.Block
	position   int
	events     <-chan ui.Event
}

func (u *UserInterface) search(searchText string) {
	buffer := runtime.NumCPU()
	u.fileList.Rows = []string{}
	u.blocks = []filter.Block{}
	u.position = 0
	for block := range filter.SearchBlocks(".", buffer, searchText) {
		u.fileList.Rows = append(u.fileList.Rows, block.FilePath)
		u.blocks = append(u.blocks, block)
	}
}

func (u *UserInterface) renderBlock(position int) {
	u.position = position
	if len(u.blocks) > position && position >= 0 {
		u.content.Text = formatter.Format(u.blocks[position])
		u.content.Title = u.blocks[position].FilePath
		ui.Render(u.grid)
	}
}

func (u *UserInterface) searchEventHandler(event ui.Event) {
	if event.ID == "<Enter>" {
		u.fileList.Title = u.searchText
		u.content.Text = ""
		u.search(u.searchText)
		u.renderBlock(0)
		return
	} else if event.ID == "<Backspace>" && len(u.searchText) > 0 {
		u.searchText = u.searchText[:len(u.searchText)-1]
		u.searchBox.Text = u.searchText
		ui.Render(u.grid)
	} else if event.ID == "<Space>" {
		u.searchText = u.searchText + " "
		u.searchBox.Text = u.searchText
		ui.Render(u.grid)
	} else if len(event.ID) == 1 {
		u.searchText = u.searchText + event.ID
		u.searchBox.Text = u.searchText
		ui.Render(u.grid)
	}
}

func (u *UserInterface) run() {
	ui.Render(u.grid)
	ticker := time.NewTicker(time.Second).C
	tickerCount := 0
	for {
		select {
		case e := <-u.events:
			switch e.ID {
			case "<C-c>":
				return
			case "<Up>":
				if u.position > 0 {
					u.renderBlock(u.position - 1)
				}
			case "<Down>":
				if u.position < len(u.blocks)-1 {
					u.renderBlock(u.position + 1)
				}
			default:
				u.searchEventHandler(e)

			}
		case <-ticker:
			tickerCount++
			ui.Render(u.grid)
		}
	}
}
func CreateInterface() *UserInterface {
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
	return &UserInterface{
		fileList:   fileList,
		searchText: "",
		searchBox:  searchBox,
		content:    content,
		grid:       grid,
		position:   0,
		blocks:     []filter.Block{},
		events:     ui.PollEvents(),
	}
}
