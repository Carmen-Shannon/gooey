package common

// SelectorState represents the state of a selector component.
// It contains information about the selector's ID, drawing state, blocking state, bounds, color, opacity, visibility, and a callback map.
type SelectorState struct {
	ID       uintptr
	Drawing  bool
	Blocking bool
	Bounds   Rect
	Color    *Color
	Opacity  float32
	Visible  bool
	CbMap    map[string]func(any)
}

type UpdateSelectorState func(state *SelectorState)

// UpdateSelectorID updates the ID of the selector.
//
// Parameters:
//   - id: The new ID to set for the selector.
//
// Returns:
//   - UpdateSelectorState: A function that takes a pointer to SelectorState and updates its ID field.
func UpdateSelectorID(id uintptr) UpdateSelectorState {
	return func(state *SelectorState) {
		state.ID = id
	}
}

// UpdateSelectorDrawing updates the active state of the selector to indicate whether it is drawing or not.
//
// Parameters:
//   - drawing: A boolean indicating whether the selector is currently drawing or not.
//
// Returns:
//   - UpdateSelectorState: A function that takes a pointer to SelectorState and updates its Drawing field.
func UpdateSelectorDrawing(drawing bool) UpdateSelectorState {
	return func(state *SelectorState) {
		state.Drawing = drawing
		state.CbMap["drawing"](drawing)
	}
}

// UpdateSelectorBlocking updates the blocking state of the selector.
//
// Parameters:
//   - blocking: A boolean indicating whether the selector is currently blocking or not.
//
// Returns:
//   - UpdateSelectorState: A function that takes a pointer to SelectorState and updates its Blocking field.
func UpdateSelectorBlocking(blocking bool) UpdateSelectorState {
	return func(state *SelectorState) {
		state.Blocking = blocking
		state.CbMap["blocking"](blocking)
	}
}

// UpdateSelectorBounds updates the bounds of the selecbools
//
// Parameters:
//   - bounds: A Rect representing the new bounds of the selector.
//
// Returns:
//   - UpdateSelectorState: A function that takes a pointer to SelectorState and updates its Bounds field.
func UpdateSelectorBounds(bounds Rect) UpdateSelectorState {
	return func(state *SelectorState) {
		state.Bounds = bounds
		state.CbMap["bounds"](bounds)
	}
}

// UpdateSelectorColor updates the color of the selector.
//
// Parameters:
//   - color: A pointer to a Color representing the new color of the selector.
//
// Returns:
//   - UpdateSelectorState: A function that takes a pointer to SelectorState and updates its Color field.
func UpdateSelectorColor(color *Color) UpdateSelectorState {
	return func(state *SelectorState) {
		state.Color = color
	}
}

// UpdateSelectorOpacity updates the opacity of the selector.
//
// Parameters:
//   - opacity: A float32 representing the new opacity of the selector.
//
// Returns:
//   - UpdateSelectorState: A function that takes a pointer to SelectorState and updates its Opacity field.
func UpdateSelectorOpacity(opacity float32) UpdateSelectorState {
	return func(state *SelectorState) {
		state.Opacity = opacity
	}
}

// UpdateSelectorVisible updates the visibility of the selector.
//
// Parameters:
//   - visible: A boolean indicating whether the selector should be visible or not.
//
// Returns:
//   - UpdateSelectorState: A function that takes a pointer to SelectorState and updates its Visible field.
func UpdateSelectorVisible(visible bool) UpdateSelectorState {
	return func(state *SelectorState) {
		state.Visible = visible
		state.CbMap["visible"](visible)
	}
}

// UpdateSelectorCbMap updates the callback map of the selector.
//
// Parameters:
//   - cbMap: A map of string keys to callback functions.
//
// Returns:
//   - UpdateSelectorState: A function that takes a pointer to SelectorState and updates its CbMap field.
func UpdateSelectorCbMap(cbMap map[string]func(any)) UpdateSelectorState {
	return func(state *SelectorState) {
		state.CbMap = cbMap
	}
}
