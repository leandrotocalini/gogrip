package main

import (
	"time"

	ui "github.com/gizak/termui/v3"
)

type Screen struct {
	progressBar      ProgressBarInterface
	searchBox        SearchBoxInterface
	contentBox       ContentBoxInterface
	searchHistoryBox SearchHistoryBoxInterface
	widgets          []WidgetInterface
	grid             *ui.Grid
	blocks           []BlockInterface
	position         int
	total            int
	events           <-chan ui.Event
}

func (u *Screen) search() {
	searchText := u.searchBox.getText()
	u.searchHistoryBox.add(searchText)
	buffer := 2
	u.blocks = []BlockInterface{}
	u.position = 0
	for block := range Search(".", buffer, searchText) {
		u.blocks = append(u.blocks, block)
		if len(u.blocks) == 1 {
			u.total = 1
			u.changePosition(0)
		}
		block.SetPosition(len(u.blocks) - 1)
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
	u.progressBar.updatePosition(u.position, u.total)
	u.contentBox.sendEvent(u.blocks[u.position])
	u.reRender()
}

func (u *Screen) reRender() {
	ui.Render(u.grid)
}

func (u *Screen) focusOnNextWidget() {
	for idx, val := range u.widgets {
		if val.isActive() {
			val.deactivate()
			next := 0
			if idx < len(u.widgets)-1 {
				next = idx + 1
			}
			u.widgets[next].activate()
			u.reRender()
			return
		}
	}
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
				if u.contentBox.isActive() {
					u.changePosition(u.position - 1)
				}
			case "<Down>":
				if u.contentBox.isActive() {
					u.changePosition(u.position + 1)
				}
			case "<Enter>":
				if u.searchBox.isActive() {
					u.search()
					u.reRender()
				}
			case "<Tab>":
				u.focusOnNextWidget()
			default:
				if u.searchBox.isActive() {
					u.searchBox.sendEvent(e.ID)
					u.reRender()
				}
			}
		case <-ticker:
			u.reRender()
		}
	}
}

func setScreen(
	grid *ui.Grid,
	searchHistoryBox *SearchHistoryBox,
	searchBox *SearchBox,
	progressBar *ProgressBar,
	contentBox *ContentBox,
) {
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewCol(1.0/4,
			searchBox.getBoxItem(),
			searchHistoryBox.getBoxItem(),
		),
		ui.NewCol((1.0/4)*3,
			progressBar.getBoxItem(),
			contentBox.getBoxItem(),
		),
	)
}

func CreateInterface() *Screen {
	progressBar := createProgressBar()
	searchBox := createSearchBox()
	searchBox.activate()
	searchHistoryBox := createSearchHistoryBox()
	contentBox := createContentBox("FILENAME", "")
	grid := ui.NewGrid()
	setScreen(grid, searchHistoryBox, searchBox, progressBar, contentBox)
	return &Screen{
		progressBar:      progressBar,
		searchBox:        searchBox,
		contentBox:       contentBox,
		searchHistoryBox: searchHistoryBox,
		widgets:          []WidgetInterface{searchBox, contentBox},
		grid:             grid,
		position:         0,
		total:            0,
		blocks:           []BlockInterface{},
		events:           ui.PollEvents(),
	}
}
