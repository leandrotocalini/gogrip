package main

import (
	"log"

	ui "github.com/gizak/termui/v3"
)

const (
	NEXT_FILE_KEY     = "<Right>"
	PREVIOUS_FILE_KEY = "<Left>"
	MOVE_UP_KEY       = "<Up>"
	MOVE_DOWN_KEY     = "<Down>"
	TAB_KEY           = "<Tab>"
	ENTER_KEY         = "<Enter>"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	userInterface := CreateInterface()
	userInterface.run()
}
