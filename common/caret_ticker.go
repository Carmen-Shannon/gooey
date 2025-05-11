package common

import (
	"sync"
	"time"
)

// CaretTicker represents a state value for the caret that is shown in TextInput rendering.
type CaretTicker struct {
	T  *time.Ticker
	Mu sync.Mutex

	Active    bool
	Visible   bool
	FocusedID uintptr
	WindowID  uintptr
}

// NewCaretTicker creates a new instance of CaretTicker with default values.
//
// Returns:
//   - *CaretTicker: A pointer to a new CaretTicker instance with default values.
func NewCaretTicker() *CaretTicker {
	return &CaretTicker{
		Mu:        sync.Mutex{},
		Active:    false,
		Visible:   false,
		FocusedID: 0,
		WindowID:  0,
	}
}

// Start begins the caret ticker for the specified window and focused ID.
// It sets the ticker to toggle the visibility of the caret every 500 milliseconds.
// If the ticker is already active, it does nothing.
//
// Parameters:
//   - windowID: The ID of the window where the caret is displayed.
//   - focusedID: The ID of the focused text input component.
func (c *CaretTicker) Start(windowID, focusedID uintptr) {
	c.Mu.Lock()
	if c.Active {
		c.Mu.Unlock()
		return
	}

	c.Active = true
	c.FocusedID = focusedID
	c.WindowID = windowID
	c.Visible = true
	c.T = time.NewTicker(500 * time.Millisecond)

	c.Mu.Unlock()
	go func() {
		for range c.T.C {
			c.Mu.Lock()
			if !c.Active {
				c.Mu.Unlock()
				return
			}
			c.Visible = !c.Visible
			c.Mu.Unlock()
		}
	}()
}

// Stop stops the caret ticker and resets its state.
// It stops the ticker and sets the Active, Visible, and FocusedID fields to their default values.
func (c *CaretTicker) Stop() {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	if c.T != nil {
		c.T.Stop()
		c.T = nil
	}
	if !c.Active && !c.Visible && c.FocusedID == 0 {
		return
	}
	c.Active = false
	c.Visible = false
	c.FocusedID = 0
}
