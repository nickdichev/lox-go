package parser

// parseError implements error interface.
type parseError struct {
	s string
}

func (p *parseError) Error() string {
	return p.s
}
