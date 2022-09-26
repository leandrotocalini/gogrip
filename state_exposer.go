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

func (s *StateExposer) updateWidget(state State) {
	rows := [][]string{}
	if state.total > 0 {
		rows = append(rows, []string{"Search string", state.searchString})
		rows = append(rows, []string{"Current file", strconv.Itoa(state.position + 1)})
		rows = append(rows, []string{"Line number", strconv.Itoa(state.currentBlock.GetLine())})
		rows = append(rows, []string{"Current Line", state.currentBlock.GetMatchedWord()})

		rows = append(rows, []string{"Amount of blocks", strconv.Itoa(state.total)})
		if state.position > 0 {
			rows = append(rows, []string{"Previous file", state.blocks[state.position-1].GetTitle()})
		} else {
			rows = append(rows, []string{"Previous file", "--"})

		}
		rows = append(rows, []string{"Current file", state.currentBlock.GetTitle()})
		if state.position < state.total-1 {
			rows = append(rows, []string{"Next file", state.blocks[state.position+1].GetTitle()})
		} else {
			rows = append(rows, []string{"Next file", "--"})

		}
	}

	s.widget.Rows = rows

}

func (s *StateExposer) listen() {
	for state := range s.channel {
		s.updateWidget(state)
	}
}

func (s *StateExposer) expose(state State) {
	s.channel <- state
}

func createStateExposer() *StateExposer {
	stateExposer := widgets.NewTable()
	stateExposer.Title = "Current state"
	stateExposer.RowSeparator = true
	stateExposer.BorderStyle = ui.NewStyle(ui.ColorGreen)

	stateExposer.Rows = [][]string{
		[]string{"", ""},
	}
	return &StateExposer{
		widget:  stateExposer,
		active:  false,
		channel: make(chan State),
	}
}
