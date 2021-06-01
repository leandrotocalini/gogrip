package filter

import "testing"

type TestScan struct {
	counter *int
}

func (t TestScan) Scan() bool {
	*t.counter++
	if *t.counter > 200 {
		return false
	}
	return true
}

func (t TestScan) Text() string {
	if *t.counter%5 == 0 {
		return "func something"
	}
	return "gbsdfsadf"
}
func TestScanFile(t *testing.T) {
	t.Run("Test get files", func(t *testing.T) {
		cn := 0
		scanner := TestScan{counter: &cn}
		f := scanFile(scanner, "asd", "func")
		if !f.Match {
			t.Errorf("Func not found!")
		}
	})
}
