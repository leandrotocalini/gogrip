package main

import (
	"fmt"
	"runtime"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/leandrotocalini/gogrip/filter"
	"github.com/leandrotocalini/gogrip/formatter"
)

type Screen struct {
	sideBar     *widgets.Gauge
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
	u.blocks = []filter.Block{}
	u.position = 0
	for block := range filter.SearchBlocks(".", buffer, searchText) {

		u.blocks = append(u.blocks, block)
	}
	u.sideBar.Percent = 0
	u.changeBlockFocus(0)
}

func (u *Screen) changeSearch(event ui.Event) {
	if event.ID == "<Backspace>" && len(u.searchText) > 0 {
		u.searchText = u.searchText[:len(u.searchText)-1]
		u.searchBox.Text = u.searchText
	} else if event.ID == "<Space>" {
		u.searchText = u.searchText + " "
		u.searchBox.Text = u.searchText
	} else if len(event.ID) == 1 {
		u.searchText = u.searchText + event.ID
		u.searchBox.Text = u.searchText
	}
	u.search(u.searchText)
	ui.Render(u.grid)
}

func (u *Screen) contentListener() {
	for position := range u.contentChan {
		u.content.Text = formatter.Format(u.blocks[position])
		u.content.Title = u.blocks[position].FilePath
		ui.Render(u.grid)
	}
}

func (u *Screen) sideBarListener() {
	for position := range u.sideBarChan {
		u.sideBar.Title = fmt.Sprintf("founds (%d/%d)", position+1, len(u.blocks))
		u.sideBar.Percent = (position + 1) * 100 / len(u.blocks)
		if u.sideBar.Percent > 80 {
			u.sideBar.BarColor = ui.ColorGreen
		} else if u.sideBar.Percent > 40 {
			u.sideBar.BarColor = ui.ColorYellow
		} else {
			u.sideBar.BarColor = ui.ColorRed
		}
		ui.Render(u.grid)
	}
}

func (u *Screen) changeBlockFocus(position int) {
	if position >= 0 && position < len(u.blocks) {
		u.position = position
		u.sideBarChan <- u.position
		u.contentChan <- u.position
	}
}

func (u *Screen) run() {
	ui.Render(u.grid)
	go u.contentListener()
	go u.sideBarListener()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-u.events:
			switch e.ID {
			case "<C-c>":
				return
			case "<Up>":
				u.changeBlockFocus(u.position - 1)
			case "<Down>":
				u.changeBlockFocus(u.position + 1)
			default:
				u.changeSearch(e)
			}
		case <-ticker:
			ui.Render(u.grid)
		}
	}
}

func CreateInterface() *Screen {
	sideBar := widgets.NewGauge()
	sideBar.Title = "Files"
	sideBar.Percent = 10
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
				ui.NewRow(1.0/15, sideBar),
			),
			ui.NewCol((1.0/4)*3, content),
		),
	)
	return &Screen{
		sideBar:     sideBar,
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
