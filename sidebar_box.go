package main

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type SidebarBox struct {
	widget   *widgets.Gauge
	position int
	total    int
	channel  chan []int
}

type SidebarBoxInterface interface {
	listen()
	setPersentage(int)
	updatePosition(int, int)
}

func (s *SidebarBox) listen() {
	for progress := range s.channel {
		s.position = progress[0]
		s.total = progress[1]
		s.refresh()
	}
}

func (s *SidebarBox) refresh() {
	s.widget.Title = fmt.Sprintf("Matches: %d/%d", s.position+1, s.total)
	if s.position == 0 {
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

func (s *SidebarBox) updatePosition(position, total int) {
	s.position = position
	s.total = total
	s.refresh()
}

func (s *SidebarBox) setPersentage(number int) {
	s.widget.Percent = number
}
func createSidebarBox() *SidebarBox {
	sideBar := widgets.NewGauge()
	sideBar.Title = "Matches"
	sideBar.Percent = 0
	sideBar.Label = " "
	return &SidebarBox{widget: sideBar, channel: make(chan []int)}
}
