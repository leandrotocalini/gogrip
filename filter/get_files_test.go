package filter

import (
	"testing"
)

func TestGet(t *testing.T) {

	t.Run("Test get files", func(t *testing.T) {
		c := getFiles(".", 10)
		cn := 0
		for range c {
			cn++
		}
		if cn != 6  {
			t.Errorf("Get() = %v, want %v", cn, 6)
		}
	})
	
}
