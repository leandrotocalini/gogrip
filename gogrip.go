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

type UserInterface struct {
	fileList  *widgets.List
	searchBox *widgets.Paragraph
	content   *widgets.Paragraph
	grid      *ui.Grid
	blocks    []filter.Block
	position  int
	events    <-chan ui.Event
}

func createInterface() *UserInterface {
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
		fileList:  fileList,
		searchBox: searchBox,
		content:   content,
		grid:      grid,
		position:  0,
		blocks:    []filter.Block{},
		events:    ui.PollEvents(),
	}
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	userInterface := createInterface()

	ui.Render(userInterface.grid)
	ticker := time.NewTicker(time.Second).C
	tickerCount := 0
	for {
		select {
		case e := <-userInterface.events:
			switch e.ID {
			case "<C-c>":
				return
			case "<C-f>":
				searchText := userInterface.getSearchString()
				userInterface.search(searchText)
				userInterface.renderBlock(0)
			case "<Up>":
				if userInterface.position > 0 {
					userInterface.renderBlock(userInterface.position - 1)
				}
			case "<Down>":
				if userInterface.position < len(userInterface.blocks)-1 {
					userInterface.renderBlock(userInterface.position + 1)
				}
			}
		case <-ticker:
			tickerCount++
			ui.Render(userInterface.grid)
		}
	}
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
	u.content.Text = formatter.Format(u.blocks[position])
	u.content.Title = u.blocks[position].FilePath
	ui.Render(u.grid)
}

func (u *UserInterface) getSearchString() string {
	text := ""
	for i := range u.events {
		if i.ID == "<Enter>" {
			u.fileList.Title = text
			return text
		}
		text = text + i.ID
		u.searchBox.Text = text
		ui.Render(u.grid)
	}
	return ""
}
