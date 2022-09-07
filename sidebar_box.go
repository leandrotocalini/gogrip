package main

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type SidebarBox struct {
	widget *widgets.Gauge
	c      chan []int
}

func (s *SidebarBox) listen() {
	for progress := range s.c {
		position := progress[0]
		total := progress[1]
		s.widget.Title = fmt.Sprintf("Matches: %d/%d", position+1, total)
		if position == 0 {
			s.widget.Percent = 0
		} else {
			s.widget.Percent = (position + 1) * 100 / total
		}
		if s.widget.Percent > 80 {
			s.widget.BarColor = ui.ColorGreen
		} else if s.widget.Percent > 40 {
			s.widget.BarColor = ui.ColorYellow
		} else {
			s.widget.BarColor = ui.ColorRed
		}
	}
}

func createSidebarBox() *SidebarBox {
	sideBar := widgets.NewGauge()
	sideBar.Title = "Matches"
	sideBar.Percent = 0
	sideBar.Label = " "

	return &SidebarBox{widget: sideBar, c: make(chan []int)}
}
