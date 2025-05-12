//go:build linux
// +build linux

package linux

func handleButtonCallbacks(id uintptr, found, pressed bool) {
	buttonCbMapMu.Lock()
	defer buttonCbMapMu.Unlock()
	for cid, cbMap := range buttonCbMap {
		if cid == id && found {
			if cb, ok := cbMap["pressed"]; ok {
				cb(pressed)
			}
			if cb, ok := cbMap["onClick"]; pressed == false && ok {
				cb(nil)
			}
		} else {
			if cb, ok := cbMap["pressed"]; ok {
				cb(false)
			}
		}
	}
}
