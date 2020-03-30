package lexer

import "fmt"

// Position present token position.
type Position struct {
	Filename string
	Offset   int
	Line     int
	Column   int
}

func (pos *Position) String() string {
	s := pos.Filename
	if s == "" {
		s = "<input>"
	}
	return fmt.Sprintf("%s line: %d, column: %d", s, pos.Line, pos.Column)
}
