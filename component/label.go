package component

import "github.com/Carmen-Shannon/gooey/common"

type label struct {
	baseComponent
	text          string
	font          string
	color         *common.Color
	textSize      int32
	textAlignment TextAlignment
	wordWrap      bool
}

// NewLabel creates a new label component.
// It accepts a variadic list of CreateLabelOption functions to customize the label's properties.
//
// Parameters:
//  - options: A variadic list of CreateLabelOption functions to customize the label's properties.
//
// Returns:
//  - Label: A pointer to the newly created label component.
func NewLabel(options ...CreateLabelOption) Label {
	opts := newCreateLabelOptions()
	for _, opt := range options {
		opt(opts)
	}
	cOpts := newCreateComponentOptions()
	for _, opt := range opts.ComponentOptions {
		opt(cOpts)
	}

	l := &label{
		baseComponent: baseComponent{
			id:      cOpts.ID,
			visible: cOpts.Visible,
			enabled: cOpts.Enabled,
			size: struct {
				Width  int32
				Height int32
			}{
				Width:  cOpts.Size.Width,
				Height: cOpts.Size.Height,
			},
			position: struct {
				X int32
				Y int32
			}{
				X: cOpts.Position.X,
				Y: cOpts.Position.Y,
			},
		},
		text:          opts.Text,
		font:          opts.Font,
		color:         opts.Color,
		textSize:      opts.TextSize,
		textAlignment: opts.TextAlignment,
	}
	return l
}

type Label interface {
	Component

	// Text returns the text of the label.
	//
	// Returns:
	//  - string: The text of the label.
	Text() string

	// SetText sets the text of the label.
	//
	// Parameters:
	//  - text: The text to set for the label.
	SetText(text string)

	// Font returns the font of the label.
	//
	// Returns:
	//  - string: The font of the label.
	Font() string

	// SetFont sets the font of the label.
	//
	// Parameters:
	//  - font: The font to set for the label.
	SetFont(font string)

	// Color returns the color of the label.
	//
	// Returns:
	//  - *common.Color: The color of the label.
	Color() *common.Color

	// SetColor sets the color of the label.
	//
	// Parameters:
	//  - color: The pointer to the common.Color to set for the label.
	SetColor(color *common.Color)

	// TextSize returns the text size of the label.
	//
	// Returns:
	//  - int32: The text size of the label.
	TextSize() int32

	// SetTextSize sets the text size of the label.
	//
	// Parameters:
	//  - size: The text size to set for the label.
	SetTextSize(size int32)

	// TextAlignment returns the text alignment of the label.
	//
	// Returns:
	//  - TextAlignment: The text alignment of the label.
	TextAlignment() TextAlignment

	// SetTextAlignment sets the text alignment of the label.
	//
	// Parameters:
	//  - alignment: The text alignment to set for the label.
	SetTextAlignment(alignment TextAlignment)

	// WordWrap returns the word wrap setting of the label.
	//
	// Returns:
	//  - bool: The word wrap setting of the label.
	WordWrap() bool

	// SetWordWrap sets the word wrap setting of the label.
	//
	// Parameters:
	//  - wordWrap: The word wrap setting to set for the label.
	SetWordWrap(wordWrap bool)
}

var _ Label = (*label)(nil)

func (l *label) Draw(ctx *common.DrawCtx) {
	drawComponent(l, ctx)
}

func (l *label) Text() string {
	return l.text
}

func (l *label) SetText(text string) {
	l.text = text
}

func (l *label) Font() string {
	return l.font
}

func (l *label) SetFont(font string) {
	l.font = font
}

func (l *label) Color() *common.Color {
	return l.color
}

func (l *label) SetColor(color *common.Color) {
	l.color = color
}

func (l *label) TextSize() int32 {
	return l.textSize
}

func (l *label) SetTextSize(size int32) {
	l.textSize = size
}

func (l *label) TextAlignment() TextAlignment {
	return l.textAlignment
}

func (l *label) SetTextAlignment(alignment TextAlignment) {
	l.textAlignment = alignment
}

func (l *label) WordWrap() bool {
	return l.wordWrap
}

func (l *label) SetWordWrap(wordWrap bool) {
	l.wordWrap = wordWrap
}
