package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScreen_InterfaceExit(t *testing.T) {
	userInterface := CreateInterface()
	key := "<C-c>"
	result := userInterface.keyEventHandler(key)
	assert.False(t, result)
}

func TestScreen_MoveUp(t *testing.T) {
	userInterface := CreateInterface()
	key := "<Up>"
	result := userInterface.keyEventHandler(key)
	assert.True(t, result)
}

func TestScreen_MoveDown(t *testing.T) {
	userInterface := CreateInterface()
	userInterface.searchBox.deactivate()
	userInterface.contentBox.activate()
	key := "<Up>"
	result := userInterface.keyEventHandler(key)
	assert.True(t, result)
	key = "<Down>"
	result = userInterface.keyEventHandler(key)
	assert.True(t, result)
}

func TestScreen_FocusOnContent(t *testing.T) {
	userInterface := CreateInterface()
	userInterface.searchBox.deactivate()
	key := "<Tab>"
	for i := 0; i < 10; i++ {
		userInterface.keyEventHandler(key)
		if userInterface.contentBox.isActive() {
			assert.True(t, true)
			return
		}
	}
	assert.True(t, false)
}

func TestScreen_Search(t *testing.T) {
	userInterface := CreateInterface()
	for _, w := range userInterface.listeners {
		go w.listen()
	}
	for _, s := range "TestScreen_Search" {
		userInterface.keyEventHandler(string(s))
	}
	userInterface.keyEventHandler("<Enter>")
	time.Sleep(2 * time.Second)
	assert.Equal(t, userInterface.state.position, 0)
}
