package main

import (
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/leandrotocalini/gogrip/filter"
)

type Screen struct {
	sideBarBox  *SidebarBox
	searchBox   *SearchBox
	contentBox  *ContentBox
	grid        *ui.Grid
	blocks      []filter.Block
	position    int
	total       int
	events      <-chan ui.Event
	sideBarChan chan int
}

func (u *Screen) search() {
	searchText := u.searchBox.getText()
	buffer := 2
	u.blocks = []filter.Block{}
	u.position = 0
	for block := range filter.SearchBlocks(".", buffer, searchText) {
		u.blocks = append(u.blocks, block)
		if len(u.blocks) == 1 {
			u.changeBlockFocus(0)
		}
	}
	u.sideBarBox.widget.Percent = 0
	u.total = len(u.blocks)
	u.position = 0
	u.changeBlockFocus(0)
}

func (u *Screen) changeBlockFocus(position int) {
	if position >= 0 && position < u.total {
		u.position = position
		u.sideBarBox.c <- []int{u.position, u.total}
		u.contentBox.contentChan <- u.blocks[u.position]
		ui.Render(u.grid)
	}
}

func (u *Screen) run() {
	ui.Render(u.grid)
	go u.contentBox.listen()
	go u.searchBox.listen()
	go u.sideBarBox.listen()
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
				u.searchBox.c <- e.ID
				u.search()
				ui.Render(u.grid)
			}
		case <-ticker:
			ui.Render(u.grid)
		}
	}
}

func CreateInterface() *Screen {
	sideBarBox := createSidebarBox()
	searchBox := createSearchBox()
	contentBox := createContentBox("FILENAME", "")
	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1,
			ui.NewCol(1.0/4,
				ui.NewRow(1.0/15, searchBox.widget),
				ui.NewRow(1.0/15, sideBarBox.widget),
			),
			ui.NewCol((1.0/4)*3, contentBox.widget),
		),
	)
	return &Screen{
		sideBarBox:  sideBarBox,
		searchBox:   searchBox,
		contentBox:  contentBox,
		grid:        grid,
		position:    0,
		total:       0,
		blocks:      []filter.Block{},
		events:      ui.PollEvents(),
		sideBarChan: make(chan int),
	}
}
