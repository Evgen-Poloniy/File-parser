package entity

import (
	"strings"
)

// Dataclass for parse error
type ParseError struct {
	LineNum  *int
	LineData []string
	Err      error
}

// New type of parse error
type ParseErrors struct {
	Errors []ParseError
}

func (p ParseErrors) Error() string {
	if len(p.Errors) == 0 {
		return ""
	}

	msgs := make([]string, len(p.Errors))
	for i, e := range p.Errors {
		msgs[i] = e.Err.Error()
	}

	return strings.Join(msgs, "; ")
}
