package filter

import (
	"regexp"
	"testing"
	"github.com/leandrotocalini/gogrip/fget"
)

func Test_filterFileChannel(t *testing.T) {

	c := make(chan Line, 10)
	go func(c chan Line) {
		defer close(c)
		c <- Line{
			Text: "Hello World",
			Ln:   0}
	}(c)
	r, _ := regexp.Compile("Hello")
	t.Run("match test", func(t *testing.T) {
		got := filterScanChannel(r, "test.go", c)
		if !got.Match {
			t.Errorf("No Match!")
		}
		if len(got.LineNumbers) != len(got.Content) {
			t.Errorf("Wrong amount of lines")

		}
	})
}

func Test_filterFileChannelMoreLines(t *testing.T) {
	c := make(chan Line, 10)
	go func(c chan Line) {
		defer close(c)
		c <- Line{
			Text: "Hello World",
			Ln:   0}
		c <- Line{
			Text: "Bye World",
			Ln:   1}
	}(c)
	r, _ := regexp.Compile("Hello")
	t.Run("match with more line test", func(t *testing.T) {
		got := filterScanChannel(r, "test.go", c)
		if !got.Match {
			t.Errorf("No Match!")
		}
		if len(got.LineNumbers) >= len(got.Content) {
			t.Errorf("Wrong amount of lines %d %d ", len(got.LineNumbers), len(got.Content))

		}
	})
}

func Test_fullFilterWithFgetChannel(t *testing.T) {
	c := fget.Get(".", 10)
	found := false
	for val := range Process("func", c, 10){
		if val.Match{
			found = true
		}
	}
	t.Run("full text in folder", func(t *testing.T) {
		if !found {
			t.Errorf("Full search in folder should match func!")
		}

	})
}