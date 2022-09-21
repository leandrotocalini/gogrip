package main

import (
	"strconv"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type StateExposer struct {
	widget  *widgets.Table
	active  bool
	channel chan State
}

func (s *StateExposer) getBoxItem() ui.GridItem {
	return ui.NewRow((1.0/15)*6, s.widget)
}

func (s *StateExposer) listen() {
	for state := range s.channel {
		if state.total > 0 {
			s.widget.Rows = [][]string{
				[]string{"Search string", state.searchString},
				[]string{"Current position", strconv.Itoa(state.position)},
				[]string{"Amount of blocks", strconv.Itoa(state.total)},
				[]string{"Current FilePath", state.currentBlock.GetTitle()},
			}
		}
	}
}

func (s *StateExposer) expose(state State) {
	s.channel <- state
}

func (s *StateExposer) isActive() bool {
	return s.active
}

func (s *StateExposer) activate() {
	s.widget.BorderStyle.Fg = ui.ColorRed
	s.active = true
}

func (s *StateExposer) deactivate() {
	s.widget.BorderStyle.Fg = ui.ColorWhite
	s.active = false
}

func createStateExposer() *StateExposer {
	sideBar := widgets.NewTable()
	sideBar.Title = "Current state"
	sideBar.RowSeparator = false
	sideBar.Rows = [][]string{
		[]string{"Key", "Value"},
	}
	return &StateExposer{
		widget:  sideBar,
		active:  false,
		channel: make(chan State),
	}
}
