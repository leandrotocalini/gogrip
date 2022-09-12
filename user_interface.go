package main

import (
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/leandrotocalini/gogrip/filter"
)

type Screen struct {
	sideBarBox SidebarBoxInterface
	searchBox  SearchBoxInterface
	contentBox ContentBoxInterface
	grid       *ui.Grid
	blocks     []filter.BlockInterface
	position   int
	total      int
	events     <-chan ui.Event
}

func (u *Screen) search() {
	searchText := u.searchBox.getText()
	buffer := 2
	u.blocks = []filter.BlockInterface{}
	u.position = 0
	for block := range filter.Search(".", buffer, searchText) {
		u.blocks = append(u.blocks, block)
		if len(u.blocks) == 1 {
			u.total = 1
			u.changePosition(0)
		}
	}
	u.total = len(u.blocks)
	u.changePosition(0)
}

func (u *Screen) changePosition(position int) {
	if position >= 0 && position < u.total {
		u.position = position
		u.focusOnBlock()
	}
}

func (u *Screen) focusOnBlock() {
	u.sideBarBox.updatePosition(u.position, u.total)
	u.contentBox.sendEvent(u.blocks[u.position])
	u.reRender()
}

func (u *Screen) reRender() {
	ui.Render(u.grid)
}

func (u *Screen) run() {
	ui.Render(u.grid)
	go u.contentBox.listen()
	go u.searchBox.listen()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-u.events:
			switch e.ID {
			case "<C-c>":
				return
			case "<Up>":
				u.changePosition(u.position - 1)
			case "<Down>":
				u.changePosition(u.position + 1)
			case "<Enter>":
				u.search()
				u.reRender()
			default:
				u.searchBox.sendEvent(e.ID)
				u.reRender()
			}
		case <-ticker:
			u.reRender()
		}
	}
}

func setScreen(grid *ui.Grid, searchBox *SearchBox, sideBarBox *SidebarBox, contentBox *ContentBox) {
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
}

func CreateInterface() *Screen {
	sideBarBox := createSidebarBox()
	searchBox := createSearchBox()
	contentBox := createContentBox("FILENAME", "")
	grid := ui.NewGrid()
	setScreen(grid, searchBox, sideBarBox, contentBox)
	return &Screen{
		sideBarBox: sideBarBox,
		searchBox:  searchBox,
		contentBox: contentBox,
		grid:       grid,
		position:   0,
		total:      0,
		blocks:     []filter.BlockInterface{},
		events:     ui.PollEvents(),
	}
}
