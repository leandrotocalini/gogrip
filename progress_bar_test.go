package main

import (
	"testing"

	ui "github.com/gizak/termui/v3"
	"github.com/stretchr/testify/assert"
)

func TestProgressBar_refreshLower(t *testing.T) {
	progressBar := createProgressBar()
	progressBar.position = 1
	progressBar.total = 10
	progressBar.refresh()
	assert.Equal(t, 10, progressBar.widget.Percent)
	assert.Equal(t, ui.ColorRed, progressBar.widget.BarColor)

}

func TestProgressBar_refreshMidle(t *testing.T) {
	progressBar := createProgressBar()
	progressBar.position = 5
	progressBar.total = 10
	progressBar.refresh()
	assert.Equal(t, 50, progressBar.widget.Percent)
	assert.Equal(t, ui.ColorYellow, progressBar.widget.BarColor)
}

func TestProgressBar_refreshHigh(t *testing.T) {
	progressBar := createProgressBar()
	progressBar.position = 85
	progressBar.total = 100
	progressBar.refresh()
	assert.Equal(t, 85, progressBar.widget.Percent)
	assert.Equal(t, ui.ColorGreen, progressBar.widget.BarColor)
}
