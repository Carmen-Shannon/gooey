package common

// TextInputState represents the state of a text input component.
//
// This is required for managing the values between the UI and the backend.
type TextInputState struct {
	Value          string
	MaxLength      int32
	Font           Font
	SelectionStart int32
	SelectionEnd   int32
	CaretPos       int32
	Focused        bool
	Bounds         struct {
		X      int32
		Y      int32
		Width  int32
		Height int32
	}
	CbMap map[string]func(any)
}

type UpdateTextInputState func(state *TextInputState)

// UpdateTIValue updates the value of the text input state.
//
// Parameters:
//   - value: The new value to set for the text input state.
//
// Returns:
//   - UpdateTextInputState: A function that takes a pointer to TextInputState and updates its value.
func UpdateTIStateValue(value string) UpdateTextInputState {
	return func(state *TextInputState) {
		state.Value = value
		state.CbMap["value"](value)
	}
}

// UpdateTIMaxLength updates the maximum length of the text input state.
//
// Parameters:
//   - maxLength: The maximum number of characters allowed in the text input state.
//
// Returns:
//   - UpdateTextInputState: A function that takes a pointer to TextInputState and updates its MaxLength field.
func UpdateTIMaxLength(maxLength int32) UpdateTextInputState {
	return func(state *TextInputState) {
		state.MaxLength = maxLength
		state.CbMap["maxLength"](maxLength)
	}
}

// UpdateTIFont updates the font of the text input state.
//
// Parameters:
//   - font: The Font to set for the text input state.
//
// Returns:
//   - UpdateTextInputState: A function that takes a pointer to TextInputState and updates its Font field.
func UpdateTIFont(font Font) UpdateTextInputState {
	return func(state *TextInputState) {
		state.Font = font
	}
}

// UpdateTISelectionStart updates the selection start position of the text input state.
//
// Parameters:
//   - selectionStart: The start position of the selection.
//
// Returns:
//   - UpdateTextInputState: A function that takes a pointer to TextInputState and updates its SelectionStart field.
func UpdateTISelectionStart(selectionStart int32) UpdateTextInputState {
	return func(state *TextInputState) {
		state.SelectionStart = selectionStart
		state.CbMap["selectionStart"](selectionStart)
	}
}

// UpdateTISelectionEnd updates the selection end position of the text input state.
//
// Parameters:
//   - selectionEnd: The end position of the selection.
//
// Returns:
//   - UpdateTextInputState: A function that takes a pointer to TextInputState and updates its SelectionEnd field.
func UpdateTISelectionEnd(selectionEnd int32) UpdateTextInputState {
	return func(state *TextInputState) {
		state.SelectionEnd = selectionEnd
		state.CbMap["selectionEnd"](selectionEnd)
	}
}

// UpdateTICaretPos updates the caret (cursor) position of the text input state.
//
// Parameters:
//   - caretPos: The position to set for the caret.
//
// Returns:
//   - UpdateTextInputState: A function that takes a pointer to TextInputState and updates its CaretPos field.
func UpdateTICaretPos(caretPos int32) UpdateTextInputState {
	return func(state *TextInputState) {
		state.CaretPos = caretPos
		state.CbMap["caretPos"](caretPos)
	}
}

// UpdateTIFocused updates the focused state of the text input state.
//
// Parameters:
//   - focused: true if the text input should be focused, false otherwise.
//
// Returns:
//   - UpdateTextInputState: A function that takes a pointer to TextInputState and updates its Focused field.
func UpdateTIFocused(focused bool) UpdateTextInputState {
	return func(state *TextInputState) {
		state.Focused = focused
		state.CbMap["focused"](focused)
	}
}

// UpdateTIBounds updates the bounds (position and size) of the text input state.
//
// Parameters:
//   - x: The x-coordinate of the text input.
//   - y: The y-coordinate of the text input.
//   - width: The width of the text input.
//   - height: The height of the text input.
//
// Returns:
//   - UpdateTextInputState: A function that takes a pointer to TextInputState and updates its Bounds fields.
func UpdateTIBounds(x, y, width, height int32) UpdateTextInputState {
	return func(state *TextInputState) {
		state.Bounds.X = x
		state.Bounds.Y = y
		state.Bounds.Width = width
		state.Bounds.Height = height
		state.CbMap["bounds"](state.Bounds)
	}
}

// UpdateTICbMap updates the callback map of the text input state.
//
// Parameters:
//   - cbMap: The map of callback functions to set for the text input state.
//
// Returns:
//   - UpdateTextInputState: A function that takes a pointer to TextInputState and updates its CbMap field.
func UpdateTICbMap(cbMap map[string]func(any)) UpdateTextInputState {
	return func(state *TextInputState) {
		state.CbMap = cbMap
	}
}
