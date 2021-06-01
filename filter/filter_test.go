package filter

import (
	"testing"
)

func Test_fullFilterWithFgetChannel(t *testing.T) {
	found := false
	for val := range FilterPath(".", 3, "func") {
		if len(val.Lines) > 0 {
			found = true
			break
		}
	}
	t.Run("full text in folder", func(t *testing.T) {
		if !found {
			t.Errorf("Full search in folder should match func!")
		}

	})
}
