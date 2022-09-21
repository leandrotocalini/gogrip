package main

type State struct {
	currentBlock BlockInterface
	searchString string
	total        int
	position     int
	blocks       []BlockInterface
}

func (s *State) moveToBlock(position int) {
	if position >= 0 && position < s.total {
		s.position = position
		s.currentBlock = s.blocks[s.position]
	}
}

func (u *State) nextBlock() {
	u.moveToBlock(u.position + 1)
}

func (u *State) previousBlock() {
	u.moveToBlock(u.position - 1)
}
