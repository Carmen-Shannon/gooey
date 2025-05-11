package component

import "gooey/common"

type textInput struct {
	baseComponent
	value          string
	maxLength      int32
	font           string
	color          *common.Color
	textColor      *common.Color
	textSize       int32
	textAlignment  TextAlignment
	focused        bool
	caretPos       int32
	selectionStart int32
	selectionEnd   int32
	state          *common.TextInputState
}

// NewTextInput creates a new TextInput component with the specified options.
// It accepts a variadic list of CreateTextInputOption functions to customize the text input's properties.
//
// Parameters:
//  - options: A variadic list of CreateTextInputOption functions to customize the text input's properties.
func NewTextInput(options ...CreateTextInputOption) TextInput {
	opts := newCreateTextInputOptions()
	for _, opt := range options {
		opt(opts)
	}
	cOpts := newCreateComponentOptions()
	for _, opt := range opts.ComponentOptions {
		opt(cOpts)
	}

	ti := &textInput{
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
		value:          opts.Value,
		maxLength:      opts.MaxLength,
		font:           opts.Font,
		color:          opts.Color,
		textColor:      opts.TextColor,
		textSize:       opts.TextSize,
		textAlignment:  opts.TextAlignment,
		focused:        false,
		caretPos:       0,
		selectionStart: 0,
		selectionEnd:   0,
		state: &common.TextInputState{
			Value:          opts.Value,
			MaxLength:      opts.MaxLength,
			Font:           common.Font{Name: opts.Font, Size: opts.TextSize},
			SelectionStart: 0,
			SelectionEnd:   0,
			CaretPos:       0,
			Focused:        false,
			Bounds: struct {
				X      int32
				Y      int32
				Width  int32
				Height int32
			}{
				X:      cOpts.Position.X,
				Y:      cOpts.Position.Y,
				Width:  cOpts.Size.Width,
				Height: cOpts.Size.Height,
			},
			CbMap: make(map[string]func(any)),
		},
	}

	cbMap := make(map[string]func(any))
	cbMap["value"] = func(newVal any) {
		if newValStr, ok := newVal.(string); ok && newValStr != ti.value && len(newValStr) <= int(ti.maxLength) {
			ti.value = newValStr
		}
	}
	cbMap["focused"] = func(focused any) {
		if focusedBool, ok := focused.(bool); ok && focusedBool != ti.focused {
			ti.focused = focusedBool
		}
	}
	cbMap["caretPos"] = func(caretPos any) {
		if caretPosInt, ok := caretPos.(int32); ok && caretPosInt != ti.caretPos {
			ti.caretPos = caretPosInt
		}
	}
	cbMap["selectionStart"] = func(selectionStart any) {
		if selectionStartInt, ok := selectionStart.(int32); ok {
			ti.selectionStart = selectionStartInt
		}
	}
	cbMap["selectionEnd"] = func(selectionEnd any) {
		if selectionEndInt, ok := selectionEnd.(int32); ok {
			ti.selectionEnd = selectionEndInt
		}
	}
	ti.state.CbMap = cbMap

	registerTextInput(ti)
	return ti
}

type TextInput interface {
	Component

	// Value returns the value of the text input.
	//
	// Returns:
	//  - string: The value of the text input.
	Value() string

	// SetValue sets the value of the text input.
	// This also manages the internal state of the TextInput component.
	//
	// Parameters:
	//  - value: The value to set for the text input.
	SetValue(value string)

	// MaxLength returns the maximum length of the text input.
	//
	// Returns:
	//  - int32: The maximum length of the text input.
	MaxLength() int32

	// SetMaxLength sets the maximum length of the text input.
	//
	// Parameters:
	//  - maxLength: The maximum number of characters allowed in the text input.
	SetMaxLength(maxLength int32)

	// Font returns the font name used by the text input.
	//
	// Returns:
	//  - string: The name of the font.
	Font() string

	// SetFont sets the font name for the text input.
	//
	// Parameters:
	//  - font: The name of the font to use.
	SetFont(font string)

	// Color returns the background color of the text input.
	//
	// Returns:
	//  - *common.Color: The background color.
	Color() *common.Color

	// SetColor sets the background color of the text input.
	//
	// Parameters:
	//  - color: The color to set as the background.
	SetColor(color *common.Color)

	// TextColor returns the color of the text.
	//
	// Returns:
	//  - *common.Color: The color of the text.
	TextColor() *common.Color

	// SetTextColor sets the color of the text.
	//
	// Parameters:
	//  - textColor: The color to use for the text.
	SetTextColor(textColor *common.Color)

	// TextSize returns the size of the text.
	//
	// Returns:
	//  - int32: The size of the text in points.
	TextSize() int32

	// SetTextSize sets the size of the text.
	//
	// Parameters:
	//  - textSize: The size of the text in points.
	SetTextSize(textSize int32)

	// TextAlignment returns the alignment of the text within the input.
	//
	// Returns:
	//  - TextAlignment: The alignment of the text.
	TextAlignment() TextAlignment

	// SetTextAlignment sets the alignment of the text within the input.
	//
	// Parameters:
	//  - textAlignment: The alignment to set for the text.
	SetTextAlignment(textAlignment TextAlignment)

	// Focused returns whether the text input currently has focus.
	//
	// Returns:
	//  - bool: true if the text input is focused, false otherwise.
	Focused() bool

	// Caret returns the current caret (cursor) position in the text input.
	//
	// Returns:
	//  - int32: The caret position.
	Caret() int32

	// Selection returns the start and end positions of the current text selection.
	//
	// Returns:
	//  - int32: The start position of the selection.
	//  - int32: The end position of the selection.
	Selection() (int32, int32)
}

var _ TextInput = (*textInput)(nil)

func (ti *textInput) Draw(ctx *common.DrawCtx) {
	drawComponent(ti, ctx)
}

func (ti *textInput) Value() string {
	return ti.value
}

func (ti *textInput) SetValue(value string) {
	ti.value = value
	ti.state.Value = value
}

func (ti *textInput) MaxLength() int32 {
	return ti.maxLength
}

func (ti *textInput) SetMaxLength(maxLength int32) {
	ti.maxLength = maxLength
	ti.state.MaxLength = maxLength
}

func (ti *textInput) Font() string {
	return ti.font
}

func (ti *textInput) SetFont(font string) {
	ti.font = font
	ti.state.Font.Name = font
}

func (ti *textInput) Color() *common.Color {
	return ti.color
}

func (ti *textInput) SetColor(color *common.Color) {
	ti.color = color
}

func (ti *textInput) TextColor() *common.Color {
	return ti.textColor
}

func (ti *textInput) SetTextColor(textColor *common.Color) {
	ti.textColor = textColor
}

func (ti *textInput) TextSize() int32 {
	return ti.textSize
}

func (ti *textInput) SetTextSize(textSize int32) {
	ti.textSize = textSize
}

func (ti *textInput) TextAlignment() TextAlignment {
	return ti.textAlignment
}

func (ti *textInput) SetTextAlignment(textAlignment TextAlignment) {
	ti.textAlignment = textAlignment
}

func (ti *textInput) Focused() bool {
	return ti.focused
}

func (ti *textInput) Caret() int32 {
	return ti.caretPos
}

func (ti *textInput) Selection() (int32, int32) {
	return ti.selectionStart, ti.selectionEnd
}
