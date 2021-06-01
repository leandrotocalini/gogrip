package filter

type Block struct {
	FirstLine int
	LastLine  int
	Lines     map[int]bool
	Content   []string
	FilePath  string
}
type Key struct {
	FirstLine int
	LastLine  int
}

func findBlockOneWay(doc []string, lineNumber, order int) int {
	found := lineNumber
	for {
		if doc[found] == "" {
			return found
		} else {
			if found > 0 && found < len(doc)-1 {
				found += order
			} else {
				return found
			}
		}
	}

	return found
}

func makeBlocks(f Found) map[Key]*Block {
	blocks := make(map[Key]*Block)
	for _, lineNumber := range f.LineNumbers {
		key := Key{
			FirstLine: findBlockOneWay(f.Content, lineNumber, -1),
			LastLine:  findBlockOneWay(f.Content, lineNumber, 1),
		}
		_, ok := blocks[key]
		if !ok {
			blocks[key] = &Block{
				FirstLine: key.FirstLine,
				LastLine:  key.LastLine,
				Lines:     make(map[int]bool),
				FilePath:  f.FilePath,
				Content:   f.Content[key.FirstLine:key.LastLine]}
		}
		blocks[key].Lines[lineNumber] = true
	}

	return blocks
}
