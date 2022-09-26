package main

import ui "github.com/gizak/termui/v3"

type GoGripWidget interface {
	getBoxItem() ui.GridItem
}

type Listener interface {
	listen()
	expose(State)
}

type GoGripListenerWidget interface {
	GoGripWidget
	Listener
}
type EventManager interface {
	GoGripListenerWidget
	isActive() bool
	activate()
	deactivate()
	newEvent(State, string) State
}
