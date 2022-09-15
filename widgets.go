package main

type WidgetInterface interface {
	isActive() bool
	activate()
	deactivate()
}
