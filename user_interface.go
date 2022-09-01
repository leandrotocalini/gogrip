package main

import (
	"runtime"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/leandrotocalini/gogrip/filter"
	"github.com/leandrotocalini/gogrip/formatter"
)

type Screen struct {
	sideBar     *widgets.List
	searchText  string
	searchBox   *widgets.Paragraph
	content     *widgets.Paragraph
	grid        *ui.Grid
	blocks      []filter.Block
	position    int
	events      <-chan ui.Event
	contentChan chan int
	sideBarChan chan int
}

func (u *Screen) search(searchText string) {
	buffer := runtime.NumCPU()
	u.sideBar.Rows = []string{}
	u.blocks = []filter.Block{}
	u.position = 0
	for block := range filter.SearchBlocks(".", buffer, searchText) {
		u.sideBar.Rows = append(u.sideBar.Rows, block.FilePath)
		u.blocks = append(u.blocks, block)
	}
	u.changeBlock(0)
}

func (u *Screen) searchEventHandler(event ui.Event) {
	if event.ID == "<Backspace>" && len(u.searchText) > 0 {
		u.searchText = u.searchText[:len(u.searchText)-1]
		u.searchBox.Text = u.searchText
		u.search(u.searchText)
		ui.Render(u.grid)
	} else if event.ID == "<Space>" {
		u.searchText = u.searchText + " "
		u.searchBox.Text = u.searchText
		ui.Render(u.grid)
	} else if len(event.ID) == 1 {
		u.searchText = u.searchText + event.ID
		u.searchBox.Text = u.searchText
		u.search(u.searchText)
		ui.Render(u.grid)
		//u.renderBlock(0)
	}
}

func (u *Screen) listenContentChan() {
	for position := range u.contentChan {
		u.content.Text = formatter.Format(u.blocks[position])
		u.content.Title = u.blocks[position].FilePath
		ui.Render(u.grid)
	}
}

func (u *Screen) listenSideBarChan() {
	for position := range u.sideBarChan {
		u.sideBar.Title = u.blocks[position].FilePath
		ui.Render(u.grid)
	}
}

func (u *Screen) changeBlock(position int) {
	if position >= 0 && position < len(u.blocks) {
		u.position = position
		u.sideBarChan <- u.position
		u.contentChan <- u.position
	}
}

func (u *Screen) run() {
	ui.Render(u.grid)
	go u.listenContentChan()
	go u.listenSideBarChan()
	ticker := time.NewTicker(time.Second).C
	tickerCount := 0
	for {
		select {
		case e := <-u.events:
			switch e.ID {
			case "<C-c>":
				return
			case "<Up>":
				u.changeBlock(u.position - 1)
			case "<Down>":
				u.changeBlock(u.position + 1)
			default:
				u.searchEventHandler(e)
			}
		case <-ticker:
			tickerCount++
			ui.Render(u.grid)
		}
	}
}
func CreateInterface() *Screen {
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
	return &Screen{
		sideBar:     fileList,
		searchText:  "",
		searchBox:   searchBox,
		content:     content,
		grid:        grid,
		position:    0,
		blocks:      []filter.Block{},
		events:      ui.PollEvents(),
		sideBarChan: make(chan int),
		contentChan: make(chan int),
	}
}
