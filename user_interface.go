package main

import (
	"time"

	ui "github.com/gizak/termui/v3"
)

type State struct {
	currentBlock BlockInterface
	searchString string
	total        int
	position     int
}

type Screen struct {
	progressBar      ProgressBarInterface
	searchBox        SearchBoxInterface
	contentBox       ContentBoxInterface
	searchHistoryBox SearchHistoryBoxInterface
	widgets          []WidgetInterface
	listeners        []WidgetInterface
	grid             *ui.Grid
	blocks           []BlockInterface
	state            State
	events           <-chan ui.Event
}

func (u *Screen) search() {
	searchText := u.searchBox.getText()
	u.state.searchString = searchText
	buffer := 2
	u.blocks = []BlockInterface{}
	u.state.position = 0
	for block := range Search(".", buffer, searchText) {
		u.blocks = append(u.blocks, block)
		if len(u.blocks) == 1 {
			u.state.total = 1
			u.moveToBlock(0)
		}
		block.SetPosition(len(u.blocks) - 1)
	}
	u.state.total = len(u.blocks)
	u.moveToBlock(0)
	u.reRender()
}

func (u *Screen) moveToBlock(position int) {
	if position >= 0 && position < u.state.total {
		u.state.position = position
		u.state.currentBlock = u.blocks[u.state.position]
		u.propagateState()
	}
}

func (u *Screen) propagateState() {
	for _, w := range u.listeners {
		w.update(u.state)
	}
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
	u.widgets[0].activate()
	u.reRender()
}

func (u *Screen) nextBlock() {
	u.moveToBlock(u.state.position + 1)
}
func (u *Screen) previousBlock() {
	u.moveToBlock(u.state.position - 1)
}

func (u *Screen) eventHandler(e ui.Event) bool {
	switch e.ID {
	case "<C-c>":
		return false
	case "<Up>":
		if u.contentBox.isActive() {
			u.previousBlock()
		}
	case "<Down>":
		if u.contentBox.isActive() {
			u.nextBlock()
		}
	case "<Enter>":
		if u.searchBox.isActive() {
			u.search()
		}
	case "<Tab>":
		u.focusOnNextWidget()
	default:
		if u.searchBox.isActive() {
			u.searchBox.sendEvent(e.ID)
			u.reRender()
		}
	}
	return true
}

func (u *Screen) run() {
	ui.Render(u.grid)
	for _, w := range u.listeners {
		go w.listen()
	}
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-u.events:
			if !u.eventHandler(e) {
				return
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
	contentBox := createContentBox()
	grid := ui.NewGrid()
	setScreen(grid, searchHistoryBox, searchBox, progressBar, contentBox)
	return &Screen{
		progressBar:      progressBar,
		searchBox:        searchBox,
		contentBox:       contentBox,
		searchHistoryBox: searchHistoryBox,
		widgets:          []WidgetInterface{searchBox, contentBox},
		listeners: []WidgetInterface{
			searchBox,
			searchHistoryBox,
			contentBox,
			progressBar,
		},
		grid: grid,
		state: State{
			position:     0,
			total:        0,
			searchString: "",
		},
		blocks: []BlockInterface{},
		events: ui.PollEvents(),
	}
}
