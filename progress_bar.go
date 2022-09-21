package main

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type ProgressBar struct {
	widget   *widgets.Gauge
	position int
	total    int
	channel  chan State
	active   bool
}

type ProgressBarInterface interface {
	GoGripWidget
	listen()
}

func (s *ProgressBar) listen() {
	for state := range s.channel {
		s.position = state.position
		s.total = state.total
		s.refresh()
	}
}

func (s *ProgressBar) getBoxItem() ui.GridItem {
	return ui.NewRow(1.0/15, s.widget)
}

func (s *ProgressBar) refresh() {
	s.widget.Title = fmt.Sprintf("%d/%d", s.position+1, s.total)
	if s.position+1 == s.total {
		s.widget.Percent = 100
	} else if s.position == 0 {
		s.widget.Percent = 0
	} else {
		s.widget.Percent = (s.position + 1) * 100 / s.total
	}

	if s.widget.Percent > 80 {
		s.widget.BarColor = ui.ColorGreen
	} else if s.widget.Percent > 40 {
		s.widget.BarColor = ui.ColorYellow
	} else {
		s.widget.BarColor = ui.ColorRed
	}
}

func (s *ProgressBar) expose(state State) {
	s.channel <- state
}

func (s *ProgressBar) isActive() bool {
	return s.active
}

func (s *ProgressBar) activate() {
	s.widget.BorderStyle.Fg = ui.ColorRed
	s.active = true
}

func (s *ProgressBar) deactivate() {
	s.widget.BorderStyle.Fg = ui.ColorWhite
	s.active = false
}

func createProgressBar() *ProgressBar {
	sideBar := widgets.NewGauge()
	sideBar.Title = "Matches"
	sideBar.Percent = 0
	sideBar.Label = " "
	return &ProgressBar{
		widget:  sideBar,
		channel: make(chan State),
		active:  false,
	}
}
