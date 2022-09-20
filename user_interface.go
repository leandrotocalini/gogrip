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
	widgets          []GoGripWidget
	listeners        []GoGripWidget
	grid             *ui.Grid
	blocks           []BlockInterface
	state            State
	active           bool
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
	u.renderScreen()
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
	u.renderScreen()
}

func (u *Screen) renderScreen() {
	if u.active {
		ui.Render(u.grid)
	}
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
			u.renderScreen()
			return
		}
	}
	u.widgets[0].activate()
	u.renderScreen()
}

func (u *Screen) nextBlock() {
	u.moveToBlock(u.state.position + 1)
}
func (u *Screen) previousBlock() {
	u.moveToBlock(u.state.position - 1)
}

func (u *Screen) keyEventHandler(key string) bool {
	switch key {
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
			u.searchBox.sendEvent(key)
			u.renderScreen()
		}
	}
	return true
}

func (u *Screen) run() {
	u.setScreen()
	u.renderScreen()
	for _, w := range u.listeners {
		go w.listen()
	}
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-u.events:
			if !u.keyEventHandler(e.ID) {
				return
			}
		case <-ticker:
			u.renderScreen()
		}
	}
}

func (u *Screen) setScreen() {
	termWidth, termHeight := ui.TerminalDimensions()
	u.grid.SetRect(0, 0, termWidth, termHeight)

	u.grid.Set(
		ui.NewCol(1.0/4,
			u.searchBox.getBoxItem(),
			u.searchHistoryBox.getBoxItem(),
		),
		ui.NewCol((1.0/4)*3,
			u.progressBar.getBoxItem(),
			u.contentBox.getBoxItem(),
		),
	)
	u.active = true
}

func CreateInterface() *Screen {
	progressBar := createProgressBar()
	searchBox := createSearchBox()
	searchBox.activate()
	searchHistoryBox := createSearchHistoryBox()
	contentBox := createContentBox()
	grid := ui.NewGrid()
	return &Screen{
		progressBar:      progressBar,
		searchBox:        searchBox,
		contentBox:       contentBox,
		searchHistoryBox: searchHistoryBox,
		widgets:          []GoGripWidget{searchBox, contentBox},
		listeners: []GoGripWidget{
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
		active: false,
		blocks: []BlockInterface{},
		events: ui.PollEvents(),
	}
}
