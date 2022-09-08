package filter

import (
	"testing"
)

func TestGet(t *testing.T) {

	t.Run("Test get files", func(t *testing.T) {
		c := getFiles(".", 10)
		cn := 0
		want := 7
		for range c {
			cn++
		}
		if cn != want {
			t.Errorf("Get() = %v, want %v", cn, want)
		}
	})

}
