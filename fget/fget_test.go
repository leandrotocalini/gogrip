package fget

import (
	"testing"
)

func TestGet(t *testing.T) {

	t.Run("Test get files", func(t *testing.T) {
		c := Get(".", 10)
		cn := 0
		for range c {
			cn++
		}
		if cn != 2  {
			t.Errorf("Get() = %v, want %v", cn, 0)
		}
	})
	
}
