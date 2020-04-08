package token

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
	// FIXME: column is wrong.
	// return fmt.Sprintf("%s line: %d, column: %d", s, pos.Line, pos.Column)
	return fmt.Sprintf("%s line: %d", s, pos.Line)
}
