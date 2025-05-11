//go:build windows
// +build windows

package wdws

import (
	"github.com/Carmen-Shannon/gooey/common"

	"golang.org/x/sys/windows"
)

// updateTextInputSelection updates the text input selection based on the provided parameters.
// It sets the selection start and end positions based on the mouse click position.
//
// Parameters:
//   - id: The ID of the text input component
//   - windowHandle: The handle to the window containing the text input
//   - mouseX: The X coordinate of the mouse click
//   - event: The event type (e.g., "start", "update", "end")
func updateTextInputSelection(id uintptr, windowHandle windows.Handle, mouseX int32, event string) {
	caret := getCaretPosForTextInput(id, windowHandle, mouseX)
	switch event {
	case "start":
		HLTR.TextInputID = id
		HLTR.Active = true
		HLTR.SelectionStart = caret
		HLTR.SelectionEnd = caret
	case "update":
		if HLTR.Active && HLTR.TextInputID == id {
			HLTR.SelectionEnd = caret
		}
	case "end":
		if HLTR.Active && HLTR.TextInputID == id {
			HLTR.SelectionEnd = caret
			HLTR.Active = false
		}
	}
	handleTextInputSelectionCallbacks(id, HLTR.SelectionStart, HLTR.SelectionEnd)
}

// getCaretPosForTextInput calculates the caret position based on the mouse click position.
// It uses the text input state to determine the bounds and font information.
// The function creates a font and uses the device context to measure the text.
//
// Parameters:
//   - id: The ID of the text input component
//   - windowHandle: The handle to the window containing the text input
//   - mouseX: The X coordinate of the mouse click
//
// Returns:
//   - int32: The calculated caret position
func getCaretPosForTextInput(id uintptr, windowHandle windows.Handle, mouseX int32) int32 {
	state := GetTextInputState(id)
	if state == nil {
		return 0
	}
	rect := [4]int32{
		state.Bounds.X,
		state.Bounds.Y,
		state.Bounds.X + state.Bounds.Width,
		state.Bounds.Y + state.Bounds.Height,
	}
	fontInfo := state.Font
	text := state.Value
	font := CreateFont(-fontInfo.Size, fontInfo.Name)
	defer DeleteObject(font)

	hdc := GetDC(windowHandle)
	defer ReleaseDC(windowHandle, hdc)

	padding := int32(4)
	caret := caretPosFromClick(hdc, font, text, rect[0], mouseX, padding)
	runes := []rune(text)
	if caret < 0 {
		caret = 0
	}
	if caret > int32(len(runes)) {
		caret = int32(len(runes))
	}
	return caret
}

// handleTextInputClickCallbacks handles the callbacks for text input components.
// It checks if the component is focused and calls the appropriate callback function.
//
// Parameters:
//   - id: The ID of the text input component
//   - found: A boolean indicating whether the component was found
func handleTextInputClickCallbacks(id uintptr, found bool, windowHandle windows.Handle, mouseX int32, doubleClick ...bool) {
	isDoubleClick := len(doubleClick) > 0 && doubleClick[0]
	if found {
		HLTR.TextInputID = id
		state := GetTextInputState(id)
		if state == nil {
			return
		}
		caretPos := getCaretPosForTextInput(id, windowHandle, mouseX)
		var selStart, selEnd int32
		if isDoubleClick {
			HLTR.SuppressSelection = true
			selStart, selEnd = getWordBounds(state.Value, caretPos)
		} else {
			selStart, selEnd = caretPos, caretPos
		}
		UpdateTextInputState(id,
			common.UpdateTIFocused(true),
			common.UpdateTICaretPos(selEnd),
			common.UpdateTISelectionStart(selStart),
			common.UpdateTISelectionEnd(selEnd),
		)
		HLTR.SelectionStart = selStart
		HLTR.SelectionEnd = selEnd
		HLTR.Active = true
		if cb, ok := state.CbMap["focused"]; ok {
			cb(true)
		}
		if !CT.Active {
			CT.Start(uintptr(windowHandle), id)
		}
		if cb, ok := state.CbMap["caretPos"]; ok {
			cb(selEnd)
		}
		if cb, ok := state.CbMap["selectionStart"]; ok {
			cb(selStart)
		}
		if cb, ok := state.CbMap["selectionEnd"]; ok {
			cb(selEnd)
		}
	} else {
		HLTR.TextInputID = 0
		HLTR.Active = false
		HLTR.SelectionStart = 0
		HLTR.SelectionEnd = 0
		for id, state := range textInputStateMap {
			UpdateTextInputState(id,
				common.UpdateTIFocused(false),
				common.UpdateTISelectionStart(0),
				common.UpdateTISelectionEnd(0),
			)
			if cb, ok := state.CbMap["focused"]; ok {
				cb(false)
			}
			if cb, ok := state.CbMap["selectionStart"]; ok {
				cb(int32(0))
			}
			if cb, ok := state.CbMap["selectionEnd"]; ok {
				cb(int32(0))
			}
		}
	}
}

// getWordBounds calculates the start and end positions of a word in the given text.
// It moves the caret to the nearest non-space character and finds the word boundaries.
//
// Parameters:
//   - text: The text in which to find the word boundaries
//   - caret: The position of the caret in the text
//
// Returns:
//   - int32: The start position of the word
//   - int32: The end position of the word
func getWordBounds(text string, caret int32) (int32, int32) {
	runes := []rune(text)
	n := int32(len(runes))
	if n == 0 || caret < 0 || caret > n {
		return 0, 0
	}
	// If caret is at the end, move back one to select the last word
	if caret == n {
		caret--
	}
	// If caret is on a space, move left to the nearest non-space
	for caret > 0 && isSeperator(runes[caret]) {
		caret--
	}
	// If still on a space, no word to select
	if isSeperator(runes[caret]) {
		return caret, caret
	}
	// Find word start
	start := caret
	for start > 0 && !isSeperator(runes[start-1]) {
		start--
	}
	// Find word end
	end := caret + 1
	for end < n && !isSeperator(runes[end]) {
		end++
	}
	return start, end
}

func isSeperator(ch rune) bool {
	return ch == ' ' || ch == '/' || ch == '\\' || ch == '.'
}

// handleTextInputSelectionCallbacks handles the selection callbacks for text input components.
// It updates the selection start and end positions in the text input state.
//
// Parameters:
//   - id: The ID of the text input component
//   - start: The start position of the selection
//   - end: The end position of the selection
func handleTextInputSelectionCallbacks(id uintptr, start, end int32) {
	UpdateTextInputState(id,
		common.UpdateTISelectionStart(start),
		common.UpdateTISelectionEnd(end),
	)
	state := GetTextInputState(id)
	if state != nil {
		if cb, ok := state.CbMap["selectionStart"]; ok {
			cb(start)
		}
		if cb, ok := state.CbMap["selectionEnd"]; ok {
			cb(end)
		}
	}
}

// handleTextInputCaretCallbacks handles the caret position callbacks for text input components.
// It updates the caret position in the text input state and calls the appropriate callback function.
//
// Parameters:
//   - id: The ID of the text input component
func handleTextInputCaretCallbacks(id uintptr) {
	state := GetTextInputState(id)
	if state != nil {
		UpdateTextInputState(id, common.UpdateTICaretPos(HLTR.SelectionEnd))
		if cb, ok := state.CbMap["caretPos"]; ok {
			cb(HLTR.SelectionEnd)
		}
	}
}

// handleTextInputChar handles the character input for text input components.
// It updates the text input state with the new value and caret position.
// It also calls the appropriate callback functions for value and caret position changes.
//
// Parameters:
//   - id: The ID of the text input component
//   - ch: The character input
func handleTextInputChar(id uintptr, ch rune) {
	if ch < 32 || ch == 127 {
		return
	}
	state := GetTextInputState(id)
	if state == nil {
		return
	}
	runes := []rune(state.Value)
	caret := HLTR.SelectionEnd
	if caret < 0 {
		caret = 0
	}
	if caret > int32(len(runes)) {
		caret = int32(len(runes))
	}
	start, end := HLTR.SelectionStart, HLTR.SelectionEnd
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if end > int32(len(runes)) {
		end = int32(len(runes))
	}

	var newVal string
	var newCaret int32
	if start != end {
		newVal = string(runes[:start]) + string(ch) + string(runes[end:])
		newCaret = start + int32(len([]rune(string(ch))))
	} else {
		newVal = string(runes[:caret]) + string(ch) + string(runes[caret:])
		newCaret = caret + int32(len([]rune(string(ch))))
	}

	if state.MaxLength > 0 && int32(len([]rune(newVal))) > state.MaxLength {
		return
	}
	UpdateTextInputState(id,
		common.UpdateTIStateValue(newVal),
		common.UpdateTISelectionStart(newCaret),
		common.UpdateTISelectionEnd(newCaret),
		common.UpdateTICaretPos(newCaret),
	)
	HLTR.SelectionStart = newCaret
	HLTR.SelectionEnd = newCaret
	if cb, ok := state.CbMap["value"]; ok {
		cb(newVal)
	}
	if cb, ok := state.CbMap["caretPos"]; ok {
		cb(newCaret)
	}
	handleTextInputSelectionCallbacks(id, newCaret, newCaret)
}

// handleTextInputBackspace handles the backspace key input for text input components.
// It deletes the character before the caret position or the selected text.
// It updates the text input state with the new value and caret position.
// It also calls the appropriate callback functions for value and caret position changes.
//
// Parameters:
//   - id: The ID of the text input component
func handleTextInputBackspace(id uintptr) {
	state := GetTextInputState(id)
	if state == nil {
		return
	}
	val := state.Value
	start, end := HLTR.SelectionStart, HLTR.SelectionEnd
	runes := []rune(val)
	if start != end {
		if start > end {
			start, end = end, start
		}
		newVal := string(runes[:start]) + string(runes[end:])
		UpdateTextInputState(id,
			common.UpdateTIStateValue(newVal),
			common.UpdateTISelectionStart(start),
			common.UpdateTISelectionEnd(start),
			common.UpdateTICaretPos(start),
		)
		if cb, ok := state.CbMap["value"]; ok {
			cb(newVal)
		}
		HLTR.SelectionStart = start
		HLTR.SelectionEnd = start
	} else if end > 0 {
		newVal := string(runes[:end-1]) + string(runes[end:])
		UpdateTextInputState(id,
			common.UpdateTIStateValue(newVal),
			common.UpdateTISelectionStart(end-1),
			common.UpdateTISelectionEnd(end-1),
			common.UpdateTICaretPos(end-1),
		)
		if cb, ok := state.CbMap["value"]; ok {
			cb(newVal)
		}
		HLTR.SelectionStart = end - 1
		HLTR.SelectionEnd = end - 1
	}
	if cb, ok := state.CbMap["caretPos"]; ok {
		cb(HLTR.SelectionEnd)
	}
	handleTextInputSelectionCallbacks(id, HLTR.SelectionStart, HLTR.SelectionEnd)
}

// handleTextInputDelete handles the delete key input for text input components.
// It deletes the character at the caret position or the selected text.
// It updates the text input state with the new value and caret position.
// It also calls the appropriate callback functions for value and caret position changes.
//
// Parameters:
//   - id: The ID of the text input component
func handleTextInputDelete(id uintptr) {
	state := GetTextInputState(id)
	if state == nil {
		return
	}
	val := state.Value
	start, end := HLTR.SelectionStart, HLTR.SelectionEnd
	runes := []rune(val)
	if start != end {
		if start > end {
			start, end = end, start
		}
		newVal := string(runes[:start]) + string(runes[end:])
		UpdateTextInputState(id,
			common.UpdateTIStateValue(newVal),
			common.UpdateTISelectionStart(start),
			common.UpdateTISelectionEnd(start),
			common.UpdateTICaretPos(start),
		)
		if cb, ok := state.CbMap["value"]; ok {
			cb(newVal)
		}
	} else if end < int32(len(runes)) {
		newVal := string(runes[:end]) + string(runes[end+1:])
		UpdateTextInputState(id,
			common.UpdateTIStateValue(newVal),
			common.UpdateTISelectionStart(end),
			common.UpdateTISelectionEnd(end),
			common.UpdateTICaretPos(end),
		)
		if cb, ok := state.CbMap["value"]; ok {
			cb(newVal)
		}
	}
	if cb, ok := state.CbMap["caretPos"]; ok {
		cb(HLTR.SelectionEnd)
	}
	handleTextInputSelectionCallbacks(id, HLTR.SelectionStart, HLTR.SelectionEnd)
}

// handleTextInputPaste handles the paste operation for text input components.
// It retrieves the text from the clipboard and inserts it at the caret position or replaces the selected text.
// It updates the text input state with the new value and caret position.
// It also calls the appropriate callback functions for value and caret position changes.
//
// Parameters:
//   - id: The ID of the text input component
func handleTextInputPaste(id uintptr) {
	state := GetTextInputState(id)
	if state == nil {
		return
	}
	clipText := getClipboardText()
	if clipText == "" {
		return
	}
	runes := []rune(state.Value)
	start, end := HLTR.SelectionStart, HLTR.SelectionEnd
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if end > int32(len(runes)) {
		end = int32(len(runes))
	}

	newVal := string(runes[:start]) + clipText + string(runes[end:])
	if state.MaxLength > 0 && int32(len([]rune(newVal))) > state.MaxLength {
		allowed := state.MaxLength - int32(len([]rune(string(runes[:start])+string(runes[end:]))))
		if allowed < 0 {
			allowed = 0
		}
		clipRunes := []rune(clipText)
		clipText = string(clipRunes[:allowed])
		newVal = string(runes[:start]) + clipText + string(runes[end:])
	}

	newCaret := start + int32(len([]rune(clipText)))
	UpdateTextInputState(id,
		common.UpdateTIStateValue(newVal),
		common.UpdateTISelectionStart(newCaret),
		common.UpdateTISelectionEnd(newCaret),
		common.UpdateTICaretPos(newCaret),
	)
	HLTR.SelectionStart = newCaret
	HLTR.SelectionEnd = newCaret
	if cb, ok := state.CbMap["value"]; ok {
		cb(newVal)
	}
	if cb, ok := state.CbMap["caretPos"]; ok {
		cb(newCaret)
	}
	handleTextInputSelectionCallbacks(id, newCaret, newCaret)
}

// handleTextInputCopy handles the copy operation for text input components.
// It copies the selected text to the clipboard.
// It retrieves the selected text and sets it to the clipboard.
//
// Parameters:
//   - id: The ID of the text input component
func handleTextInputCopy(id uintptr) {
	state := GetTextInputState(id)
	if state == nil {
		return
	}
	start, end := HLTR.SelectionStart, HLTR.SelectionEnd
	if start > end {
		start, end = end, start
	}
	runes := []rune(state.Value)
	if start < 0 || end > int32(len(runes)) || start == end {
		return // nothing to copy
	}
	text := string(runes[start:end])
	setClipboardText(text)
}
