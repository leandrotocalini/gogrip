package filter

import (
	"io"
	"regexp"
	"testing"
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
		got := filterFileChannel(r, "test.go", c)
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
		got := filterFileChannel(r, "test.go", c)
		if !got.Match {
			t.Errorf("No Match!")
		}
		if len(got.LineNumbers) >= len(got.Content) {
			t.Errorf("Wrong amount of lines %d %d ", len(got.LineNumbers), len(got.Content))

		}
	})
}
