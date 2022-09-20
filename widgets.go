package main

import ui "github.com/gizak/termui/v3"

type GoGripWidget interface {
	isActive() bool
	activate()
	deactivate()
	update(State)
	listen()
	getBoxItem() ui.GridItem
}
