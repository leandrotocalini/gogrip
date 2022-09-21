package main

import ui "github.com/gizak/termui/v3"

type GoGripWidget interface {
	listen()
	getBoxItem() ui.GridItem
}

type ExposerWidget interface {
	GoGripWidget
	expose(State)
}

type EventManagerWidget interface {
	GoGripWidget
	isActive() bool
	activate()
	deactivate()
	newEvent(State, string) State
}
