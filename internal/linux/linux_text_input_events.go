//go:build linux
// +build linux

package linux

/*
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <X11/Xutil.h>
#include <X11/cursorfont.h>
#include <stdlib.h>
*/
import "C"
import (
	"os/exec"

	"github.com/Carmen-Shannon/gooey/common"
)

func handleTextInputClickCallbacks(id uintptr, found bool, hwnd uintptr, mouseX int32, doubleClick ...bool) {
	isDoubleClick := len(doubleClick) > 0 && doubleClick[0]

	// Suppress further selection updates if a double-click just occurred
	if HLTR.SuppressSelection && !isDoubleClick {
		return
	}

	if found {
		HLTR.TextInputID = id
		state := GetTextInputState(id)
		if state == nil {
			return
		}
		caretPos := getCaretPosForTextInput(id, hwnd, mouseX)
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
			CT.Start(hwnd, id)
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

// Copy selected text to clipboard using xclip
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
	cmd := exec.Command("xclip", "-selection", "clipboard")
	in, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	if err := cmd.Start(); err != nil {
		return
	}
	in.Write([]byte(text))
	in.Close()
	cmd.Wait()
}

// Paste clipboard text at caret using xclip
func handleTextInputPaste(id uintptr) {
	state := GetTextInputState(id)
	if state == nil {
		return
	}
	out, err := exec.Command("xclip", "-selection", "clipboard", "-o").Output()
	if err != nil {
		return
	}
	clipText := string(out)
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

// handleTextInputBackspace removes the selected text or the character before the caret.
func handleTextInputBackspace(id uintptr) {
	state := GetTextInputState(id)
	if state == nil {
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

	var newVal string
	var newCaret int32
	if start != end {
		// Remove selection
		newVal = string(runes[:start]) + string(runes[end:])
		newCaret = start
	} else if start > 0 {
		// Remove character before caret
		newVal = string(runes[:start-1]) + string(runes[start:])
		newCaret = start - 1
	} else {
		// Nothing to delete
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

// handleTextInputDelete removes the selected text or the character at the caret.
func handleTextInputDelete(id uintptr) {
	state := GetTextInputState(id)
	if state == nil {
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

	var newVal string
	var newCaret int32
	if start != end {
		// Remove selection
		newVal = string(runes[:start]) + string(runes[end:])
		newCaret = start
	} else if end < int32(len(runes)) {
		// Remove character at caret
		newVal = string(runes[:end]) + string(runes[end+1:])
		newCaret = end
	} else {
		// Nothing to delete
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

func updateTextInputSelection(id uintptr, hwnd uintptr, mouseX int32, event string) {
	if HLTR.SuppressSelection {
		return
	}
	caret := getCaretPosForTextInput(id, hwnd, mouseX)
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

func getCaretPosForTextInput(id uintptr, hwnd uintptr, mouseX int32) int32 {
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
	// You may need to implement or adapt a font measurement function for Linux/X11
	fontName := fontInfo.Name
	fontSize := fontInfo.Size

	display := GetDisplay(hwnd)
	window := C_Window(hwnd)
	// Use the window as the drawable for measurement
	padding := int32(4)
	return caretPosFromClickLinux(display, window, fontName, fontSize, text, rect[0], mouseX, padding)
}

func caretPosFromClickLinux(display *C_Display, drawable C_Drawable, fontName string, fontSize int32, text string, inputX, clickX, padding int32) int32 {
	clickOffset := clickX - inputX - padding
	if clickOffset <= 0 {
		return 0
	}
	runes := []rune(text)
	low, high := 0, len(runes)
	for low < high {
		mid := (low + high) / 2
		sub := string(runes[:mid])
		w := int32(XTextWidth(display, drawable, fontName, int(fontSize), sub))
		if w < clickOffset {
			low = mid + 1
		} else {
			high = mid
		}
	}
	if low > 0 && low <= len(runes) {
		prevW := int32(XTextWidth(display, drawable, fontName, int(fontSize), string(runes[:low-1])))
		currW := int32(XTextWidth(display, drawable, fontName, int(fontSize), string(runes[:low])))
		midpoint := (prevW + currW) / 2
		if clickOffset < midpoint {
			return int32(low - 1)
		}
	}
	return int32(low)
}

// getWordBounds calculates the start and end positions of a word in the given text.
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
	for caret > 0 && isSeparator(runes[caret]) {
		caret--
	}
	// If still on a space, no word to select
	if isSeparator(runes[caret]) {
		return caret, caret
	}
	// Find word start
	start := caret
	for start > 0 && !isSeparator(runes[start-1]) {
		start--
	}
	// Find word end
	end := caret + 1
	for end < n && !isSeparator(runes[end]) {
		end++
	}
	return start, end
}

func isSeparator(ch rune) bool {
	return ch == ' ' || ch == '/' || ch == '\\' || ch == '.'
}
