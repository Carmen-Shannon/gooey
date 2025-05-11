package common

// DrawCtx represents the drawing context for the main window, it allows the window information to pass in an agnostic way to the draw calls.
type DrawCtx struct {
	Hwnd uintptr
	Hdc  uintptr
}
