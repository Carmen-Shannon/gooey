//go:build linux
// +build linux

package component

import (
	"github.com/Carmen-Shannon/gooey/common"
	"github.com/Carmen-Shannon/gooey/internal/linux"
)

func drawSelector(_ *common.DrawCtx, s Selector) {
	state := s.(*selector).state

	// Launch overlay if needed
	if state.Visible && !linux.SelectorOverlayActive() {
		linux.LaunchSelectorOverlayOnThread(s.ID())
	}

	// If selector is visible and drawing, force overlay redraw on every update
	if state.Visible && state.Drawing {
		linux.ForceSelectorOverlayRedraw()
	}

	// If selector is not visible but overlay is active, destroy overlay
	if !state.Visible && linux.SelectorOverlayActive() {
		linux.DestroySelectorOverlay()
	}
}

// registerSelector registers the selector state with the windows drawing context.
//
// Parameters:
//   - s: Selector to register with the windows drawing context.
func registerSelector(componentID uintptr, s Selector) {
	linux.RegisterSelectorState(componentID, s.(*selector).state)
}
