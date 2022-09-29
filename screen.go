package main

import (
	"time"

	ui "github.com/gizak/termui/v3"
)

type Screen struct {
	progressBar      ProgressBarInterface
	searchBox        EventManager
	contentBox       ContentBoxInterface
	searchHistoryBox SearchHistoryBoxInterface
	stateExposer     StateExposer
	eventManagers    []EventManager
	exposers         []Listener
	grid             *ui.Grid
	state            State
	active           bool
	events           <-chan ui.Event
}

func (u *Screen) propagateState() {
	for _, w := range u.exposers {
		w.expose(u.state)
	}
	u.renderScreen()
}

func (u *Screen) renderScreen() {
	if u.active {
		ui.Render(u.grid)
	}
}

func (u *Screen) focusOnNextWidget() {
	for idx, val := range u.eventManagers {
		if val.isActive() {
			val.deactivate()
			next := 0
			if idx < len(u.eventManagers)-1 {
				next = idx + 1
			}
			u.eventManagers[next].activate()
			u.renderScreen()
			return
		}
	}
	u.eventManagers[0].activate()
	u.renderScreen()
}

func (u *Screen) sendEventToActiveEventManager(key string) {
	for _, val := range u.eventManagers {
		if val.isActive() {
			u.state = val.newEvent(u.state, key)
			u.propagateState()
		} else {
			val.newEvent(u.state, key)
		}

	}
}
func (u *Screen) keyEventHandler(key string) bool {
	switch key {
	case "<C-c>":
		return false
	case "<Tab>":
		u.focusOnNextWidget()
	default:
		u.sendEventToActiveEventManager(key)
	}
	return true
}

func (u *Screen) run() {
	u.setScreen()
	u.renderScreen()
	for _, w := range u.exposers {
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
		ui.NewCol(1.0/6,
			u.searchBox.getBoxItem(),
			u.searchHistoryBox.getBoxItem(),
			u.stateExposer.getBoxItem(),
		),
		ui.NewCol((1.0/6)*5,
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
	stateExposer := createStateExposer()
	grid := ui.NewGrid()
	return &Screen{
		progressBar:      progressBar,
		searchBox:        searchBox,
		contentBox:       contentBox,
		searchHistoryBox: searchHistoryBox,
		stateExposer:     *stateExposer,
		eventManagers:    []EventManager{searchBox, contentBox},
		exposers: []Listener{
			searchBox,
			searchHistoryBox,
			contentBox,
			progressBar,
			stateExposer,
		},
		grid: grid,
		state: State{
			position:     0,
			total:        0,
			searchString: "",
			blocks:       []BlockInterface{},
			cached:       false,
		},
		active: false,
		events: ui.PollEvents(),
	}
}
