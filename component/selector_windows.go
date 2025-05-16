//go:build windows
// +build windows

package component

import (
	"github.com/Carmen-Shannon/gooey/common"
	wdws "github.com/Carmen-Shannon/gooey/internal/windows"
	"golang.org/x/sys/windows"
)

// drawSelector handles the drawing of the selector component on the windows platform.
// This function uses windows-specific APIs to create the rendered component.
// It handles the selector's color, opacity, and bounds.
//
// Parameters:
//   - ctx: The drawing context for the selector.
//   - s: The Selector component to be drawn.
func drawSelector(_ *common.DrawCtx, s Selector) {
	state := wdws.GetSelectorState(s.ID())

	if state.Visible {
		if state.ID == 0 {
			hwnd := LaunchSelectorOverlayOnThread(s)
			if hwnd != 0 {
				state.ID = uintptr(hwnd)
				wdws.UpdateSelectorState(s.ID(), common.UpdateSelectorID(uintptr(hwnd)))
			}
		}
		// Only make interactive if Drawing && !Blocking
		if state.Visible && state.ID != 0 && state.Drawing && !state.Blocking {
			wdws.SetTransparentStyle(windows.Handle(state.ID), false)
			state.Blocking = true // Now we're blocking input until mouse event
			wdws.UpdateSelectorState(s.ID(), common.UpdateSelectorBlocking(true))
		}
	} else if state.ID != 0 {
		wdws.PostMessage(windows.Handle(state.ID), wdws.WM_CLOSE, 0, 0)
		state.ID = 0
		state.Blocking = false // Reset blocking on hide
		wdws.UpdateSelectorState(s.ID(), common.UpdateSelectorID(0), common.UpdateSelectorBlocking(false))
	}
}

// registerSelector registers the selector state with the windows drawing context.
//
// Parameters:
//   - s: Selector to register with the windows drawing context.
func registerSelector(componentID uintptr, s Selector) {
	wdws.RegisterSelectorState(componentID, s.(*selector).state)
}
