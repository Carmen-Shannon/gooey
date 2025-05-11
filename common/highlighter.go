package common

// Highlighter is the structure used to represent the state of the text input highlighter.
// Since only one text input can be focused at a time, we can use a single instance of this struct to manage the state across events.
type Highlighter struct {
	TextInputID       uintptr
	SelectionStart    int32
	SelectionEnd      int32
	Active            bool
	SuppressSelection bool
}

// NewHighlighter creates a new instance of a Highlighter with default values.
//
// Returns:
//   - *Highlighter: A pointer to a new Highlighter instance with default values.
func NewHighlighter() *Highlighter {
	return &Highlighter{
		TextInputID:       0,
		SelectionStart:    0,
		SelectionEnd:      0,
		Active:            false,
		SuppressSelection: false,
	}
}
