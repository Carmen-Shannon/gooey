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
#define GO_FALSE 0
#define GO_TRUE 1
#define GO_CLIENT_MESSAGE 33
#define GO_SUBSTRUCTURE_REDIRECT_MASK (1L<<20)
#define GO_SUBSTRUCTURE_NOTIFY_MASK (1L<<19)

void set_client_message_event(
    XEvent *event,
    Display *display,
    Window window,
    Atom message_type,
    long d0, long d1, long d2, long d3, long d4
) {
    event->xclient.type = GO_CLIENT_MESSAGE;
    event->xclient.serial = 0;
    event->xclient.send_event = 1;
    event->xclient.display = display;
    event->xclient.window = window;
    event->xclient.message_type = message_type;
    event->xclient.format = 32;
    event->xclient.data.l[0] = d0;
    event->xclient.data.l[1] = d1;
    event->xclient.data.l[2] = d2;
    event->xclient.data.l[3] = d3;
    event->xclient.data.l[4] = d4;
}
*/
import "C"
import (
	"sync"
	"unsafe"

	"github.com/Carmen-Shannon/gooey/common"
)

type C_Window = C.Window
type C_XEvent = C.XEvent
type C_Drawable = C.Drawable
type C_Display = C.Display

var (
	displayMap          = make(map[uintptr]uintptr)
	displayMapMu        sync.Mutex
	drawCallbackMap     = make(map[uintptr]func(hdc uintptr))
	drawCallbackMu      sync.Mutex
	resizingState       = make(map[uintptr]bool)
	resizingStateMu     sync.Mutex
	customCursorDraw    = false
	customCursorDrawMu  sync.Mutex
	wdwColorMap         = make(map[uintptr]common.Color)
	wdwColorMapMu       sync.Mutex
	buttonBoundsMap     = make(map[uintptr][4]int32)
	buttonBoundsMapMu   sync.Mutex
	buttonCbMap         = make(map[uintptr]map[string]func(any))
	buttonCbMapMu       sync.Mutex
	textInputStateMap   = make(map[uintptr]*common.TextInputState)
	textInputStateMapMu sync.Mutex

	// Caret Ticker \\
	CT   = common.NewCaretTicker()
	HLTR = common.NewHighlighter()
)

const (
	// C Flags
	C_EXPOSE          = 12
	C_CONFIGURENOTIFY = 22
	C_DESTROYNOTIFY   = 17
	C_BUTTONPRESS     = 4
	C_BUTTONRELEASE   = 5
	C_MOTIONNOTIFY    = 6

	// Text Alignment
	ALIGN_LEFT       = 0
	ALIGN_CENTER     = 1
	ALIGN_RIGHT      = 2
	ALIGN_VCENTER    = 4
	ALIGN_WORDBREAK  = 8
	ALIGN_SINGLELINE = 16

	// Window Events
	ExposureMask        = 1 << 15
	StructureNotifyMask = 1 << 17
	ButtonPressMask     = 1 << 2
	ButtonReleaseMask   = 1 << 3
	PointerMotionMask   = 1 << 6
	KeyPressMask        = 1 << 0
	KeyReleaseMask      = 1 << 1
)

func WindowProc(hwnd uintptr, display *C.Display, event *C.XEvent) bool {
	switch EventType(event) {
	case C_EXPOSE:
		HandlePaint(hwnd, display)
		return true
	case C_CONFIGURENOTIFY:
		HandlePaint(hwnd, display)
		return true
	case C_DESTROYNOTIFY:
		XCloseDisplay(display)
		UnregisterDisplay(hwnd)
		return false
	case C_BUTTONPRESS:
		x, y := GetMouseState(hwnd)
		btnId, btnFound := FindButtonAt(x, y)
		handleButtonCallbacks(btnId, btnFound, true)
		tiId, tiFound := FindTextInputAt(x, y)
		handleTextInputClickCallbacks(tiId, tiFound, hwnd, x)
		return true
	case C_BUTTONRELEASE:
		x, y := GetMouseState(hwnd)
		btnId, btnFound := FindButtonAt(x, y)
		handleButtonCallbacks(btnId, btnFound, false)
		tiId, tiFound := FindTextInputAt(x, y)
		if tiFound && HLTR.TextInputID == tiId {
			updateTextInputSelection(tiId, hwnd, x, "end")
			handleTextInputCaretCallbacks(tiId)
		}
		HLTR.Active = false
		return true
	case C_MOTIONNOTIFY:
		x, y := GetMouseState(hwnd)
		if HLTR.Active && HLTR.TextInputID != 0 && !HLTR.SuppressSelection {
			tiId, tiFound := FindTextInputAt(x, y)
			if tiFound && tiId == HLTR.TextInputID {
				updateTextInputSelection(tiId, hwnd, x, "update")
			}
		}
		return true
	default:
	}
	return true
}

// handlePaint for Linux: calls the registered draw callback
func HandlePaint(hwnd uintptr, display *C.Display) {
	cb := getDrawCallback(hwnd)
	if cb == nil {
		return
	}

	window := C_Window(hwnd)

	// Get window size
	var attrs C.XWindowAttributes
	C.XGetWindowAttributes(display, window, &attrs)
	width := int(attrs.width)
	height := int(attrs.height)
	if width <= 0 || height <= 0 {
		return
	}

	// Create off-screen pixmap (double buffer)
	pixmap := C.XCreatePixmap(display, window, C.uint(width), C.uint(height), C.uint(attrs.depth))
	defer C.XFreePixmap(display, pixmap)

	// Create GC for pixmap
	gc := C.XCreateGC(display, pixmap, 0, nil)
	defer C.XFreeGC(display, gc)

	// Fill background (optional: get window color)
	bgColor := GetWindowColor(hwnd)
	var pixel C.ulong = 0xffffff
	if bgColor != nil {
		pixel = C.ulong((uint32(bgColor.Red) << 16) | (uint32(bgColor.Green) << 8) | uint32(bgColor.Blue))
	}
	C.XSetForeground(display, gc, pixel)
	C.XFillRectangle(display, pixmap, gc, 0, 0, C.uint(width), C.uint(height))

	// Draw all components to the pixmap
	// Pass the pixmap as the "hdc" to the draw callback
	cb(uintptr(pixmap))

	// Copy the pixmap to the window in one operation
	C.XCopyArea(display, pixmap, window, gc, 0, 0, C.uint(width), C.uint(height), 0, 0)

	// Flush to ensure drawing is visible
	C.XFlush(display)
}

// LoadArrowCursor returns the standard pointer cursor.
func LoadArrowCursor(display *C.Display) C.Cursor {
	return C.XCreateFontCursor(display, C.XC_left_ptr)
}

// LoadIBeamCursor returns the I-beam (text) cursor.
func LoadIBeamCursor(display *C.Display) C.Cursor {
	return C.XCreateFontCursor(display, C.XC_xterm)
}

// SetCursor sets the cursor for the given window.
func SetCursor(display *C.Display, window C.Window, cursor C.Cursor) {
	C.XDefineCursor(display, window, cursor)
	C.XFlush(display)
}

func XPending(display *C.Display) int {
	return int(C.XPending(display))
}

func XFlush(display *C.Display) {
	C.XFlush(display)
}

func XSelectInput(display *C.Display, window C.Window, eventMask int) {
	C.XSelectInput(display, window, C.long(eventMask))
}

func XOpenDisplay() *C.Display {
	return C.XOpenDisplay(nil)
}

func XDefaultScreen(display *C.Display) C.int {
	return C.XDefaultScreen(display)
}

func XRootWindow(display *C.Display, screen C.int) C.Window {
	return C.XRootWindow(display, screen)
}

func XCreateSimpleWindow(
	display *C.Display,
	parent C.Window,
	x, y int,
	width, height, borderWidth uint,
	border, background uint32,
) C.Window {
	return C.XCreateSimpleWindow(
		display,
		parent,
		C.int(x), C.int(y),
		C.uint(width), C.uint(height), C.uint(borderWidth),
		C.ulong(border), C.ulong(background),
	)
}

func XMapWindow(display *C.Display, window C.Window) {
	C.XMapWindow(display, window)
}

func XStoreName(display *C.Display, window C.Window, name string) {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	C.XStoreName(display, window, cstr)
}

func XNextEvent(display *C.Display, event *C.XEvent) {
	C.XNextEvent(display, event)
}

func XCloseDisplay(display *C.Display) {
	C.XCloseDisplay(display)
}

// Show window
func XShowWindow(display *C.Display, window C.Window) {
	C.XMapWindow(display, window)
}

// Hide window
func XHideWindow(display *C.Display, window C.Window) {
	C.XUnmapWindow(display, window)
}

// Maximize window (EWMH)
func XMaximizeWindow(display *C.Display, window C.Window) {
	setNetWMState(display, window, "_NET_WM_STATE_MAXIMIZED_VERT", true)
	setNetWMState(display, window, "_NET_WM_STATE_MAXIMIZED_HORZ", true)
}

// Minimize window (iconify)
func XMinimizeWindow(display *C.Display, window C.Window) {
	C.XIconifyWindow(display, window, C.XDefaultScreen(display))
}

// Fill a rectangle with a color
func XFillRect(display *C.Display, drawable C.Drawable, x, y, w, h int, color *common.Color) {
	gc := C.XCreateGC(display, drawable, 0, nil)
	defer C.XFreeGC(display, gc)
	pixel := (uint32(color.Red) << 16) | (uint32(color.Green) << 8) | uint32(color.Blue)
	C.XSetForeground(display, gc, C.ulong(pixel))
	C.XFillRectangle(display, drawable, gc, C.int(x), C.int(y), C.uint(w), C.uint(h))
}

// Fill a rounded rectangle with a color
func XFillRoundedRect(display *C.Display, drawable C.Drawable, x, y, w, h, radius int, color *common.Color) {
	gc := C.XCreateGC(display, drawable, 0, nil)
	defer C.XFreeGC(display, gc)
	pixel := (uint32(color.Red) << 16) | (uint32(color.Green) << 8) | uint32(color.Blue)
	C.XSetForeground(display, gc, C.ulong(pixel))
	C.XFillArc(display, drawable, gc, C.int(x), C.int(y), C.uint(radius*2), C.uint(radius*2), 0, 23040)
	C.XFillArc(display, drawable, gc, C.int(x+w-radius*2), C.int(y), C.uint(radius*2), C.uint(radius*2), 0, 23040)
	C.XFillArc(display, drawable, gc, C.int(x), C.int(y+h-radius*2), C.uint(radius*2), C.uint(radius*2), 0, 23040)
	C.XFillArc(display, drawable, gc, C.int(x+w-radius*2), C.int(y+h-radius*2), C.uint(radius*2), C.uint(radius*2), 0, 23040)
	C.XFillRectangle(display, drawable, gc, C.int(x+radius), C.int(y), C.uint(w-2*radius), C.uint(h))
	C.XFillRectangle(display, drawable, gc, C.int(x), C.int(y+radius), C.uint(w), C.uint(h-2*radius))
}

// Draw centered text in a rectangle
func XDrawTextCentered(display *C.Display, drawable C.Drawable, x, y, w, h int, fontName string, fontSize int, text string, color *common.Color) {
	gc := C.XCreateGC(display, drawable, 0, nil)
	defer C.XFreeGC(display, gc)
	pixel := (uint32(color.Red) << 16) | (uint32(color.Green) << 8) | uint32(color.Blue)
	C.XSetForeground(display, gc, C.ulong(pixel))

	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	fontStruct := C.XQueryFont(display, C.XGContextFromGC(gc))
	if fontStruct != nil {
		textWidth := C.XTextWidth(fontStruct, cstr, C.int(len(text)))
		textX := x + (w-int(textWidth))/2
		textY := y + h/2 + fontSize/2 // crude vertical centering
		C.XDrawString(display, drawable, gc, C.int(textX), C.int(textY), cstr, C.int(len(text)))
	}
}

// XTextWidth measures the width of the text in pixels for the given font and size.
func XTextWidth(display *C.Display, drawable C.Drawable, fontName string, fontSize int, text string) int {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	gc := C.XCreateGC(display, drawable, 0, nil)
	defer C.XFreeGC(display, gc)
	fontStruct := C.XQueryFont(display, C.XGContextFromGC(gc))
	if fontStruct != nil {
		return int(C.XTextWidth(fontStruct, cstr, C.int(len(text))))
	}
	return 0
}

// XDrawTextRect draws text in a rectangle with alignment and word wrap.
func XDrawTextRect(display *C.Display, drawable C.Drawable, x, y, w, h int, fontName string, fontSize int, text string, color *common.Color, format int) {
	gc := C.XCreateGC(display, drawable, 0, nil)
	defer C.XFreeGC(display, gc)
	pixel := (uint32(color.Red) << 16) | (uint32(color.Green) << 8) | uint32(color.Blue)
	C.XSetForeground(display, gc, C.ulong(pixel))

	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))

	// Alignment
	textWidth := XTextWidth(display, drawable, fontName, fontSize, text)
	var textX, textY int
	switch format & 3 { // ALIGN_LEFT, ALIGN_CENTER, ALIGN_RIGHT
	case ALIGN_LEFT:
		textX = x
	case ALIGN_CENTER:
		textX = x + (w-textWidth)/2
	case ALIGN_RIGHT:
		textX = x + w - textWidth
	}
	// Vertical alignment (simple)
	textY = y + h/2 + fontSize/2

	// Word wrap not implemented here, but you could split text and draw multiple lines.
	C.XDrawString(display, drawable, gc, C.int(textX), C.int(textY), cstr, C.int(len(text)))
}

// Draw a rectangle border
func XDrawRect(display *C.Display, drawable C.Drawable, x, y, w, h int, color *common.Color) {
	gc := C.XCreateGC(display, drawable, 0, nil)
	defer C.XFreeGC(display, gc)
	pixel := (uint32(color.Red) << 16) | (uint32(color.Green) << 8) | uint32(color.Blue)
	C.XSetForeground(display, gc, C.ulong(pixel))
	C.XDrawRectangle(display, drawable, gc, C.int(x), C.int(y), C.uint(w-1), C.uint(h-1))
}

func XClearArea(display *C.Display, window C.Window, x, y, w, h int, exposures bool) {
	var exp C.Bool
	if exposures {
		exp = 1
	} else {
		exp = 0
	}
	C.XClearArea(display, window, C.int(x), C.int(y), C.uint(w), C.uint(h), exp)
}

// GetMouseState returns the current mouse position relative to the window.
func GetMouseState(hwnd uintptr) (x, y int32) {
	display := GetDisplay(hwnd)
	if display == nil {
		return 0, 0
	}
	window := C_Window(hwnd)
	var root, child C.Window
	var rootX, rootY, winX, winY C.int
	var mask C.uint
	C.XQueryPointer(display, window, &root, &child, &rootX, &rootY, &winX, &winY, &mask)
	return int32(winX), int32(winY)
}

// Helper to send _NET_WM_STATE client messages
func setNetWMState(display *C.Display, window C.Window, state string, add bool) {
	atom := C.XInternAtom(display, C.CString("_NET_WM_STATE"), C.GO_FALSE)
	prop := C.XInternAtom(display, C.CString(state), C.GO_FALSE)
	var action C.long
	if add {
		action = 1 // _NET_WM_STATE_ADD
	} else {
		action = 0 // _NET_WM_STATE_REMOVE
	}
	var event C.XEvent
	C.set_client_message_event(
		&event, display, window, atom,
		action, C.long(prop), 0, 0, 0,
	)
	C.XSendEvent(display, C.XDefaultRootWindow(display), C.GO_FALSE,
		C.GO_SUBSTRUCTURE_REDIRECT_MASK|C.GO_SUBSTRUCTURE_NOTIFY_MASK, &event)
}

func EventType(event *C_XEvent) int {
	return int(*(*C.int)(unsafe.Pointer(event)))
}

func RegisterDisplay(hwnd uintptr, display *C.Display) {
	displayMapMu.Lock()
	defer displayMapMu.Unlock()
	displayMap[hwnd] = uintptr(unsafe.Pointer(display))
}

func GetDisplay(hwnd uintptr) *C.Display {
	displayMapMu.Lock()
	defer displayMapMu.Unlock()
	if ptr, ok := displayMap[hwnd]; ok {
		return (*C.Display)(unsafe.Pointer(ptr))
	}
	return nil
}

func UnregisterDisplay(hwnd uintptr) {
	displayMapMu.Lock()
	defer displayMapMu.Unlock()
	delete(displayMap, hwnd)
}

// RegisterDrawCallback registers a callback function to be called when the window needs to be redrawn.
// It takes a window handle and a callback function as parameters.
//
// Parameters:
//   - hwnd: The handle to the window
//   - cb: The callback function to be called when the window needs to be redrawn
func RegisterDrawCallback(hwnd uintptr, cb func(hdc uintptr)) {
	drawCallbackMu.Lock()
	defer drawCallbackMu.Unlock()
	drawCallbackMap[hwnd] = cb
}

// getDrawCallback retrieves the callback function associated with a window handle.
// It takes a window handle as a parameter and returns the callback function.
//
// Parameters:
//   - hwnd: The handle to the window
//
// Returns:
//   - func(hdc uintptr): The callback function associated with the window handle
func getDrawCallback(hwnd uintptr) func(hdc uintptr) {
	drawCallbackMu.Lock()
	defer drawCallbackMu.Unlock()
	return drawCallbackMap[hwnd]
}

// SetResizingState sets the resizing state for a window handle.
//
// Parameters:
//   - hwnd: The handle to the window
//   - resizing: A boolean value indicating whether the window is being resized
func SetResizingState(hwnd uintptr, resizing bool) {
	resizingStateMu.Lock()
	defer resizingStateMu.Unlock()
	resizingState[hwnd] = resizing
}

// IsResizing returns true if the window is currently being resized.
//
// Parameters:
//   - hwnd: The handle to the window
//
// Returns:
//   - bool: true if the window is being resized, false otherwise
func IsResizing(hwnd uintptr) bool {
	resizingStateMu.Lock()
	defer resizingStateMu.Unlock()
	return resizingState[hwnd]
}

// SetWindowColor sets the background color for a particular window handle.
//
// Parameters:
//   - hwnd: The handle to the window
//   - color: A pointer to a common.Color struct representing the background color
func SetWindowColor(hwnd uintptr, color *common.Color) {
	wdwColorMapMu.Lock()
	defer wdwColorMapMu.Unlock()
	wdwColorMap[hwnd] = *color
}

// GetWindowColor retrieves the background color for a particular window handle.
//
// Parameters:
//   - hwnd: The handle to the window
//
// Returns:
//   - *common.Color: A pointer to a common.Color struct representing the background color, or nil if not set
func GetWindowColor(hwnd uintptr) *common.Color {
	wdwColorMapMu.Lock()
	defer wdwColorMapMu.Unlock()
	color, ok := wdwColorMap[hwnd]
	if !ok {
		return nil
	}
	return &color
}

// RegisterButtonBounds registers the bounds of a button control.
// It is used to track the position and size of the button control.
// This is useful for handling mouse events and determining if the button control is being interacted with.
//
// Parameters:
//   - componentID: The ID of the component associated with the button control
//   - bounds: An array of four integers representing the bounds of the button control
func RegisterButtonBounds(componentID uintptr, bounds [4]int32) {
	buttonBoundsMapMu.Lock()
	defer buttonBoundsMapMu.Unlock()
	buttonBoundsMap[componentID] = bounds
}

// GetButtonBounds retrieves the bounds of a button control.
// It is used to get the position and size of the button control.
// This is useful for handling mouse events and determining if the button control is being interacted with.
//
// Parameters:
//   - componentID: The ID of the component associated with the button control
//
// Returns:
//   - rect: An array of four integers representing the bounds of the button control
//   - ok: A boolean indicating whether the bounds were found
func GetButtonBounds(componentID uintptr) (rect [4]int32, ok bool) {
	buttonBoundsMapMu.Lock()
	defer buttonBoundsMapMu.Unlock()
	rect, ok = buttonBoundsMap[componentID]
	return
}

// FindButtonAt checks if a point (x, y) is within the bounds of any button control.
// It is used to determine if the mouse is over a button control.
// This is useful for handling mouse events and determining if the button control is being interacted with.
//
// Parameters:
//   - x: The x-coordinate of the point
//   - y: The y-coordinate of the point
//
// Returns:
//   - componentID: The ID of the button control if found, 0 otherwise
//   - found: A boolean indicating whether the button control was found
func FindButtonAt(x, y int32) (componentID uintptr, found bool) {
	buttonBoundsMapMu.Lock()
	defer buttonBoundsMapMu.Unlock()
	for id, rect := range buttonBoundsMap {
		if x >= rect[0] && x < rect[0]+rect[2] && y >= rect[1] && y < rect[1]+rect[3] {
			return id, true
		}
	}
	return 0, false
}

// RegisterButtonCallback registers a callback function for button events.
// It is used to handle button click events and other button-related actions.
//
// Parameters:
//   - componentID: The ID of the component associated with the button
//   - cbMap: A map of callback functions for different button events
func RegisterButtonCallback(componentID uintptr, cbMap map[string]func(any)) {
	buttonCbMapMu.Lock()
	defer buttonCbMapMu.Unlock()
	buttonCbMap[componentID] = cbMap
}

// GetButtonCallback retrieves the callback function for a given component ID.
// It is used to retrieve the callback function associated with a specific component ID.
//
// Parameters:
//   - componentID: The ID of the component
//
// Returns:
//   - func(any): The callback function to be called when the button is clicked
func GetButtonCallback(componentID uintptr) map[string]func(any) {
	buttonCbMapMu.Lock()
	defer buttonCbMapMu.Unlock()
	return buttonCbMap[componentID]
}

// RegisterTextInputState registers the state of a text input control.
// It is used to track the state of the text input control, including its bounds and other properties.
// This is useful for handling text input events and managing the state of the text input control.
//
// Parameters:
//   - componentID: The ID of the component associated with the text input control
//   - state: A pointer to a common.TextInputState struct representing the state of the text input control
func RegisterTextInputState(componentID uintptr, state *common.TextInputState) {
	textInputStateMapMu.Lock()
	defer textInputStateMapMu.Unlock()
	textInputStateMap[componentID] = state
}

// GetTextInputState retrieves the state of a text input control.
// It is used to get the current state of the text input control, including its bounds and other properties.
// This is useful for handling text input events and managing the state of the text input control.
//
// Parameters:
//   - componentID: The ID of the component associated with the text input control
//
// Returns:
//   - *common.TextInputState: A pointer to a common.TextInputState struct representing the state of the text input control
func GetTextInputState(componentID uintptr) *common.TextInputState {
	textInputStateMapMu.Lock()
	defer textInputStateMapMu.Unlock()
	return textInputStateMap[componentID]
}

// UpdateTextInputState updates the state of a text input control.
// It takes a component ID and a variadic number of update functions to modify the state.
// This is useful for applying multiple updates to the state of the text input control.
//
// Parameters:
//   - componentID: The ID of the component associated with the text input control
//   - updates: A variadic number of update functions to modify the state
func UpdateTextInputState(componentID uintptr, updates ...common.UpdateTextInputState) {
	textInputStateMapMu.Lock()
	defer textInputStateMapMu.Unlock()
	state, ok := textInputStateMap[componentID]
	if !ok {
		return
	}
	for _, update := range updates {
		update(state)
	}
}

// FindTextInputAt checks if a point (x, y) is within the bounds of any text input control.
// It is used to determine if the mouse is over a text input control.
// This is useful for handling mouse events and determining if the text input control is being interacted with.
//
// Parameters:
//   - x: The x-coordinate of the point
//   - y: The y-coordinate of the point
//
// Returns:
//   - componentID: The ID of the text input control if found, 0 otherwise
//   - found: A boolean indicating whether the text input control was found
func FindTextInputAt(x, y int32) (componentID uintptr, found bool) {
	textInputStateMapMu.Lock()
	defer textInputStateMapMu.Unlock()
	for id, state := range textInputStateMap {
		if x >= state.Bounds.X && x < state.Bounds.X+state.Bounds.Width &&
			y >= state.Bounds.Y && y < state.Bounds.Y+state.Bounds.Height {
			return id, true
		}
	}
	return 0, false
}

// SetCustomCursorDraw enables or disables custom cursor drawing.
// It is used to control whether the custom cursor drawing logic should be applied.
//
// Parameters:
//   - enabled: A boolean value indicating whether custom cursor drawing should be enabled
func SetCustomCursorDraw(enabled bool) {
	customCursorDrawMu.Lock()
	defer customCursorDrawMu.Unlock()
	customCursorDraw = enabled
}

// isCustomCursorDraw checks if custom cursor drawing is enabled.
// It returns true if custom cursor drawing is enabled, false otherwise.
//
// Returns:
//   - bool: true if custom cursor drawing is enabled, false otherwise
func IsCustomCursorDraw() bool {
	customCursorDrawMu.Lock()
	defer customCursorDrawMu.Unlock()
	return customCursorDraw
}
