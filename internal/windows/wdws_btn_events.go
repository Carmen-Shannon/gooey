//go:build windows
// +build windows

package wdws

// handleButtonCallbacks handles the callbacks for button components.
// It checks if the button is pressed and calls the appropriate callback function.
// It takes an additional argument `doCb` to determine if the callback should be executed.
// This is especially useful for detecting mouse up events.
//
// Parameters:
//   - id: The ID of the button component
//   - found: A boolean indicating whether the button was found
//   - doCb: A boolean indicating whether to execute the callback
func handleButtonCallbacks(id uintptr, found, doCb bool) {
	for cid, cbMap := range buttonCbMap {
		if cid == id && found {
			if cb, ok := cbMap["pressed"]; ok {
				cb(doCb)
			}
			if cb, ok := cbMap["onClick"]; doCb && ok {
				cb(nil)
			}
		} else {
			if cb, ok := cbMap["pressed"]; ok {
				cb(false)
			}
		}
	}
}
