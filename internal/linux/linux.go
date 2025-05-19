//go:build linux
// +build linux

package linux

/*
#cgo LDFLAGS: -lX11 -lXrender
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <X11/Xutil.h>
#include <X11/cursorfont.h>
#include <X11/extensions/Xrender.h>
#include <stdlib.h>
#define GO_FALSE 0
#define GO_TRUE 1
#define GO_CLIENT_MESSAGE 33
#define GO_SUBSTRUCTURE_REDIRECT_MASK (1L<<20)
#define GO_SUBSTRUCTURE_NOTIFY_MASK (1L<<19)

void gooey_x11_init_threads() { XInitThreads(); }

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
	"runtime"
	"sync"
	"unsafe"

	"github.com/Carmen-Shannon/gooey/common"
)

type C_Window = C.Window
type C_XEvent = C.XEvent
type C_Drawable = C.Drawable
type C_Display = C.Display
type C_char = C.char
type C_KeySym = C.KeySym

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
	selectorStateMap    = make(map[uintptr]*common.SelectorState)
	selectorStateMapMu  sync.Mutex

	overlay          *selectorOverlay
	fallbackSelector *fallbackSelectorState
	argbSelector     *argbOverlay

	// Caret Ticker \\
	CT   = common.NewCaretTicker()
	HLTR = common.NewHighlighter()

	// local event tracking \\
	lastClickTimeMu sync.Mutex
	lastClickTime   = make(map[uintptr]C.Time)
	lastClickPos    = make(map[uintptr][2]int32)
)

const (
	// C Flags
	C_EXPOSE          = 12
	C_CONFIGURENOTIFY = 22
	C_DESTROYNOTIFY   = 17
	C_BUTTONPRESS     = 4
	C_BUTTONRELEASE   = 5
	C_MOTIONNOTIFY    = 6
	C_KEYPRESS        = 2

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

	// Special Keys
	XK_BACKSPACE = 0xff08
	XK_DELETE    = 0xffff

	// Event Listening Logic
	doubleClickThresholdMs = 400
	doubleClickMaxDist     = 4
)

type fallbackSelectorState struct {
	display       *C.Display
	inputWin      C.Window // Fullscreen InputOnly overlay (invisible, intercepts all input)
	selectorWin   C.Window // Top-level InputOutput window for the selector rectangle
	screen        C.int
	active        bool
	selectorDrawn bool
}

type argbOverlay struct {
	display  *C.Display
	window   C.Window
	screen   C.int
	visual   *C.XVisualInfo
	colormap C.Colormap
	active   bool
}

type selectorOverlay struct {
	display *C.Display
	window  C.Window
	screen  C.int
	gc      C.GC
	active  bool
}

func init() {
	C.gooey_x11_init_threads()
}

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
	case C_KEYPRESS:
		if HLTR.TextInputID != 0 {
			keyEvent := (*C.XKeyEvent)(unsafe.Pointer(event))
			ctrlDown := (keyEvent.state & C.ControlMask) != 0
			var keysym C.KeySym
			C.XLookupString(keyEvent, nil, 0, &keysym, nil)
			switch keysym {
			case XK_BACKSPACE:
				handleTextInputBackspace(HLTR.TextInputID)
			case XK_DELETE:
				handleTextInputDelete(HLTR.TextInputID)
			case 0x0063, 0x0043: // 'c' or 'C'
				if ctrlDown {
					handleTextInputCopy(HLTR.TextInputID)
					return true
				}
			case 0x0076, 0x0056: // 'v' or 'V'
				if ctrlDown {
					handleTextInputPaste(HLTR.TextInputID)
					return true
				}
			case 0x0078, 0x0058: // 'x' or 'X'
				if ctrlDown {
					handleTextInputCopy(HLTR.TextInputID)
					handleTextInputBackspace(HLTR.TextInputID)
					return true
				}
			default:
				handleTextInputKeyPress(HLTR.TextInputID, hwnd, event, display)
			}
		}
		return true
	case C_BUTTONPRESS:
		x, y := GetMouseState(hwnd)
		btnId, btnFound := FindButtonAt(x, y)
		handleButtonCallbacks(btnId, btnFound, true)
		tiId, tiFound := FindTextInputAt(x, y)

		ev := (*C.XButtonEvent)(unsafe.Pointer(event))
		dblClk := isDoubleClick(hwnd, ev, x, y)
		handleTextInputClickCallbacks(tiId, tiFound, hwnd, x, dblClk)
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
		HLTR.SuppressSelection = false
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

func isDoubleClick(hwnd uintptr, event *C.XButtonEvent, x, y int32) bool {
	clkTime := event.time

	lastClickTimeMu.Lock()
	prevTime := lastClickTime[hwnd]
	prevPos := lastClickPos[hwnd]
	lastClickTime[hwnd] = clkTime
	lastClickPos[hwnd] = [2]int32{x, y}
	lastClickTimeMu.Unlock()

	if prevTime != 0 && int(clkTime-prevTime) < doubleClickThresholdMs &&
		abs32(int32(x)-int32(prevPos[0])) < doubleClickMaxDist &&
		abs32(int32(y)-int32(prevPos[1])) < doubleClickMaxDist {
		return true
	}
	return false
}

func abs32(a int32) int32 {
	if a < 0 {
		return -a
	}
	return a
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

func handleTextInputKeyPress(id uintptr, hwnd uintptr, event *C.XEvent, display *C.Display) {
	keyEvent := (*C.XKeyEvent)(unsafe.Pointer(event))
	var buf [8]C.char
	var keysym C.KeySym
	n := C.XLookupString(keyEvent, &buf[0], 8, &keysym, nil)
	if n <= 0 {
		return
	}
	input := string(C.GoBytes(unsafe.Pointer(&buf[0]), n))
	runes := []rune(input)
	if len(runes) == 0 {
		return
	}
	ch := runes[0]
	if ch < 32 || ch == 127 {
		return
	}
	handleTextInputChar(id, ch)
}

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

func getScreenSize(display *C.Display, screen C.int) (int, int) {
	width := int(C.XDisplayWidth(display, screen))
	height := int(C.XDisplayHeight(display, screen))
	return width, height
}

// ForceSelectorOverlayRedraw triggers an Expose event on the overlay window to force a redraw.
func ForceSelectorOverlayRedraw() {
	if overlay == nil || !overlay.active {
		return
	}
	var event C.XEvent
	(*(*C.int)(unsafe.Pointer(&event))) = C.Expose
	C.XSendEvent(overlay.display, overlay.window, 0, 0, &event)
	C.XFlush(overlay.display)
}

// Find a 32-bit TrueColor visual
func findARGBVisual(display *C.Display, screen C.int) *C.XVisualInfo {
	var vinfo C.XVisualInfo
	vinfo.depth = 32

	// Set the 'class' field (hack, since cgo does not expose it directly)
	classOffset := unsafe.Offsetof(vinfo.depth) + unsafe.Sizeof(vinfo.depth)
	classPtr := (*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(&vinfo)) + classOffset))
	*classPtr = C.TrueColor

	var n C.int
	visuals := C.XGetVisualInfo(display, C.VisualDepthMask|C.VisualClassMask, &vinfo, &n)
	if visuals != nil && n > 0 {
		visualArr := (*[1 << 16]C.XVisualInfo)(unsafe.Pointer(visuals))[:n:n]
		for i := 0; i < int(n); i++ {
			if visualArr[i].screen == screen {
				return &visualArr[i]
			}
		}
		return &visualArr[0]
	}
	return nil
}

// Launch the ARGB selector overlay on a new thread
func LaunchSelectorOverlayOnThread(sID uintptr) {
	go func() {
		runtime.LockOSThread()
		state := GetSelectorState(sID)
		if state == nil {
			return
		}

		display := C.XOpenDisplay(nil)
		if display == nil {
			panic("Cannot open X display")
		}
		screen := C.XDefaultScreen(display)
		root := C.XRootWindow(display, screen)
		screenW := C.XDisplayWidth(display, screen)
		screenH := C.XDisplayHeight(display, screen)

		visual := findARGBVisual(display, screen)
		if visual == nil {
			C.XCloseDisplay(display)
			panic("No 32-bit TrueColor visual found")
		}

		colormap := C.XCreateColormap(display, root, visual.visual, C.AllocNone)

		var attrs C.XSetWindowAttributes
		attrs.colormap = colormap
		attrs.background_pixel = 0 // fully transparent
		attrs.border_pixel = 0
		attrs.override_redirect = 1

		win := C.XCreateWindow(
			display, root,
			0, 0, C.uint(screenW), C.uint(screenH),
			0,
			visual.depth,
			C.InputOutput,
			visual.visual,
			C.CWColormap|C.CWBackPixel|C.CWBorderPixel|C.CWOverrideRedirect,
			&attrs,
		)
		if win == 0 {
			C.XCloseDisplay(display)
			panic("Failed to create ARGB overlay window")
		}

		C.XMapRaised(display, win)
		C.XFlush(display)

		argbSelector = &argbOverlay{display, win, screen, visual, colormap, true}

		// Select input events
		C.XSelectInput(display, win, C.ExposureMask|C.ButtonPressMask|C.ButtonReleaseMask|C.PointerMotionMask|C.KeyPressMask)

		var (
			startX, startY int32
			curX, curY     int32
			dragging       bool
		)

		// Helper to draw the overlay
		drawOverlay := func(rect common.Rect, color *common.Color, opacity float32) {
			// Use XRender to draw a semi-transparent rectangle
			pictFormat := C.XRenderFindVisualFormat(display, visual.visual)
			pictWin := C.XRenderCreatePicture(display, C.Drawable(win), pictFormat, 0, nil)
			defer C.XRenderFreePicture(display, pictWin)

			// Clear the window (fully transparent)
			C.XClearWindow(display, win)

			// Prepare color with alpha
			alpha := uint16(opacity * 65535)
			renderColor := C.XRenderColor{
				red:   C.ushort(color.Red) * 257,
				green: C.ushort(color.Green) * 257,
				blue:  C.ushort(color.Blue) * 257,
				alpha: C.ushort(alpha),
			}

			// Fill rectangle
			C.XRenderFillRectangle(display, C.PictOpOver, pictWin, &renderColor,
				C.int(rect.X), C.int(rect.Y), C.uint(rect.W), C.uint(rect.H))

			// Draw border (opaque)
			borderColor := C.XRenderColor{
				red:   0,
				green: 0,
				blue:  0,
				alpha: 65535,
			}
			thick := 2
			// Top
			C.XRenderFillRectangle(display, C.PictOpOver, pictWin, &borderColor,
				C.int(rect.X), C.int(rect.Y), C.uint(rect.W), C.uint(thick))
			// Bottom
			C.XRenderFillRectangle(display, C.PictOpOver, pictWin, &borderColor,
				C.int(rect.X), C.int(rect.Y+rect.H-int32(thick)), C.uint(rect.W), C.uint(thick))
			// Left
			C.XRenderFillRectangle(display, C.PictOpOver, pictWin, &borderColor,
				C.int(rect.X), C.int(rect.Y), C.uint(thick), C.uint(rect.H))
			// Right
			C.XRenderFillRectangle(display, C.PictOpOver, pictWin, &borderColor,
				C.int(rect.X+rect.W-int32(thick)), C.int(rect.Y), C.uint(thick), C.uint(rect.H))

			C.XFlush(display)
		}

		// Initial overlay (full screen, transparent)
		drawOverlay(common.Rect{X: 0, Y: 0, W: int32(screenW), H: int32(screenH)}, common.ColorBlack, 0.0)

		for argbSelector != nil && argbSelector.active {
			var event C.XEvent
			C.XNextEvent(display, &event)
			eventType := int(*(*C.int)(unsafe.Pointer(&event)))

			switch eventType {
			case C.Expose:
				currentState := GetSelectorState(sID)
				if currentState != nil && currentState.Visible {
					drawOverlay(currentState.Bounds, currentState.Color, currentState.Opacity)
				}
			case C.ButtonPress:
				buttonEvent := (*C.XButtonEvent)(unsafe.Pointer(&event))
				if buttonEvent.button == 1 && !dragging {
					startX, startY = int32(buttonEvent.x_root), int32(buttonEvent.y_root)
					curX, curY = startX, startY
					dragging = true
					C.XGrabPointer(display, win, 1, C.ButtonPressMask|C.ButtonReleaseMask|C.PointerMotionMask,
						C.GrabModeAsync, C.GrabModeAsync, C.None, C.None, C.CurrentTime)
					UpdateSelectorState(sID, common.UpdateSelectorDrawing(true), common.UpdateSelectorBlocking(true))
					drawOverlay(common.Rect{X: startX, Y: startY, W: 0, H: 0}, state.Color, state.Opacity)
				}
			case C.MotionNotify:
				if dragging {
					motionEvent := (*C.XMotionEvent)(unsafe.Pointer(&event))
					curX, curY = int32(motionEvent.x_root), int32(motionEvent.y_root)
					rectX, rectY := startX, startY
					rectW, rectH := curX-startX, curY-startY
					if rectW < 0 {
						rectX = curX
						rectW = -rectW
					}
					if rectH < 0 {
						rectY = curY
						rectH = -rectH
					}
					newBounds := common.Rect{X: rectX, Y: rectY, W: rectW, H: rectH}
					UpdateSelectorState(sID, common.UpdateSelectorBounds(newBounds))
					drawOverlay(newBounds, state.Color, state.Opacity)
				}
			case C.ButtonRelease:
				buttonEvent := (*C.XButtonEvent)(unsafe.Pointer(&event))
				if buttonEvent.button == 1 && dragging {
					dragging = false
					C.XUngrabPointer(display, C.CurrentTime)
					rectX, rectY := startX, startY
					rectW, rectH := curX-startX, curY-startY
					if rectW < 0 {
						rectX = curX
						rectW = -rectW
					}
					if rectH < 0 {
						rectY = curY
						rectH = -rectH
					}
					finalBounds := common.Rect{X: rectX, Y: rectY, W: rectW, H: rectH}
					UpdateSelectorState(sID,
						common.UpdateSelectorBounds(finalBounds),
						common.UpdateSelectorBlocking(false),
						common.UpdateSelectorDrawing(false),
						common.UpdateSelectorVisible(false),
					)
					drawOverlay(finalBounds, state.Color, state.Opacity)
					argbSelector.active = false
				}
			case C.KeyPress:
				keyEvent := (*C.XKeyEvent)(unsafe.Pointer(&event))
				keysym := C.XLookupKeysym(keyEvent, 0)
				if keysym == C.XK_Escape {
					UpdateSelectorState(sID,
						common.UpdateSelectorBlocking(false),
						common.UpdateSelectorDrawing(false),
						common.UpdateSelectorVisible(false),
					)
					argbSelector.active = false
				}
			}
			if argbSelector == nil || !argbSelector.active {
				break
			}
		}

		C.XDestroyWindow(display, win)
		C.XFreeColormap(display, colormap)
		C.XCloseDisplay(display)
		argbSelector = nil
	}()
}

func CreateSelectorOverlay(x, y, w, h int32) {
	display := C.XOpenDisplay(nil)
	if display == nil {
		panic("Cannot open X display")
	}
	screen := C.XDefaultScreen(display)
	root := C.XRootWindow(display, screen)

	visual := findARGBVisual(display, screen)
	if visual != nil {
		// ...existing ARGB path unchanged...
		colormap := C.XCreateColormap(display, root, visual.visual, C.AllocNone)
		var attrs C.XSetWindowAttributes
		attrs.colormap = colormap
		attrs.override_redirect = 1
		attrs.background_pixel = 0 // fully transparent

		win := C.XCreateWindow(
			display, root,
			C.int(x), C.int(y), C.uint(w), C.uint(h),
			0,
			visual.depth,
			C.InputOutput,
			visual.visual,
			C.CWColormap|C.CWBackPixel|C.CWOverrideRedirect,
			&attrs,
		)
		C.XMapWindow(display, win)
		C.XFlush(display)

		argbSelector = &argbOverlay{
			display:  display,
			window:   win,
			screen:   screen,
			visual:   visual,
			colormap: colormap,
			active:   true,
		}
		fallbackSelector = nil
		return
	}

	// Fallback: 2 top-level windows
	var attrs C.XSetWindowAttributes
	attrs.override_redirect = 1
	attrs.event_mask = C.ButtonPressMask | C.ButtonReleaseMask | C.PointerMotionMask | C.KeyPressMask

	screenW := C.XDisplayWidth(display, screen)
	screenH := C.XDisplayHeight(display, screen)
	inputWin := C.XCreateWindow(
		display, root,
		0, 0, C.uint(screenW), C.uint(screenH),
		0,
		0, // depth
		C.InputOnly,
		C.XDefaultVisual(display, screen),
		C.CWOverrideRedirect|C.CWEventMask,
		&attrs,
	)
	C.XMapRaised(display, inputWin)
	C.XFlush(display)

	fallbackSelector = &fallbackSelectorState{
		display:       display,
		inputWin:      inputWin,
		screen:        screen,
		active:        true,
		selectorDrawn: false,
		selectorWin:   0,
	}
	argbSelector = nil
}

// UpdateSelectorOverlay draws or updates the selector rectangle as a separate top-level window.
func UpdateSelectorOverlay(x, y, w, h int32, color *common.Color, opacity float32) {
	// ARGB path unchanged...
	if argbSelector != nil && argbSelector.active {
		display := argbSelector.display
		win := argbSelector.window
		visual := argbSelector.visual

		// Use XRender to draw a semi-transparent rectangle
		pictFormat := C.XRenderFindVisualFormat(display, visual.visual)
		pictWin := C.XRenderCreatePicture(display, C.Drawable(win), pictFormat, 0, nil)
		defer C.XRenderFreePicture(display, pictWin)

		// Clear the window (fully transparent)
		C.XClearWindow(display, win)

		// Prepare color with alpha
		alpha := uint16(opacity * 65535)
		renderColor := C.XRenderColor{
			red:   C.ushort(color.Red) * 257,
			green: C.ushort(color.Green) * 257,
			blue:  C.ushort(color.Blue) * 257,
			alpha: C.ushort(alpha),
		}

		// Fill rectangle
		C.XRenderFillRectangle(display, C.PictOpOver, pictWin, &renderColor,
			C.int(x), C.int(y), C.uint(w), C.uint(h))

		// Draw border (opaque)
		borderColor := C.XRenderColor{
			red:   0,
			green: 0,
			blue:  0,
			alpha: 65535,
		}
		thick := 2
		// Top
		C.XRenderFillRectangle(display, C.PictOpOver, pictWin, &borderColor,
			C.int(x), C.int(y), C.uint(w), C.uint(thick))
		// Bottom
		C.XRenderFillRectangle(display, C.PictOpOver, pictWin, &borderColor,
			C.int(x), C.int(y+int32(h)-int32(thick)), C.uint(w), C.uint(thick))
		// Left
		C.XRenderFillRectangle(display, C.PictOpOver, pictWin, &borderColor,
			C.int(x), C.int(y), C.uint(thick), C.uint(h))
		// Right
		C.XRenderFillRectangle(display, C.PictOpOver, pictWin, &borderColor,
			C.int(x+int32(w)-int32(thick)), C.int(y), C.uint(thick), C.uint(h))

		C.XFlush(display)
		return
	}

	// Fallback: 2-window path
	if fallbackSelector != nil && fallbackSelector.active {
		display := fallbackSelector.display
		screen := fallbackSelector.screen
		root := C.XRootWindow(display, screen)

		// Destroy previous selector window if it exists
		if fallbackSelector.selectorDrawn && fallbackSelector.selectorWin != 0 {
			C.XUnmapWindow(display, fallbackSelector.selectorWin)
			C.XDestroyWindow(display, fallbackSelector.selectorWin)
			fallbackSelector.selectorWin = 0
			fallbackSelector.selectorDrawn = false
		}

		// Only draw if w/h > 0
		if w > 0 && h > 0 {
			var attrs C.XSetWindowAttributes
			attrs.override_redirect = 1
			attrs.background_pixel = C.ulong((uint32(color.Red) << 16) | (uint32(color.Green) << 8) | uint32(color.Blue))
			attrs.event_mask = 0

			selectorWin := C.XCreateWindow(
				display, root, // Top-level, NOT a child of InputOnly window!
				C.int(x), C.int(y), C.uint(w), C.uint(h),
				0,
				C.CopyFromParent,
				C.InputOutput,
				C.XDefaultVisual(display, screen),
				C.CWBackPixel|C.CWOverrideRedirect,
				&attrs,
			)
			C.XMapRaised(display, selectorWin)
			C.XFlush(display)

			// Draw border (black)
			gc := C.XCreateGC(display, selectorWin, 0, nil)
			C.XSetForeground(display, gc, 0)
			thick := 2
			C.XFillRectangle(display, selectorWin, gc, 0, 0, C.uint(w), C.uint(thick))
			C.XFillRectangle(display, selectorWin, gc, 0, C.int(h)-C.int(thick), C.uint(w), C.uint(thick))
			C.XFillRectangle(display, selectorWin, gc, 0, 0, C.uint(thick), C.uint(h))
			C.XFillRectangle(display, selectorWin, gc, C.int(w)-C.int(thick), 0, C.uint(thick), C.uint(h))
			C.XFreeGC(display, gc)

			fallbackSelector.selectorWin = selectorWin
			fallbackSelector.selectorDrawn = true
		}
	}
}

// DestroySelectorOverlay cleans up both windows in the fallback path.
func DestroySelectorOverlay() {
	if argbSelector != nil {
		C.XUnmapWindow(argbSelector.display, argbSelector.window)
		C.XDestroyWindow(argbSelector.display, argbSelector.window)
		C.XFreeColormap(argbSelector.display, argbSelector.colormap)
		C.XCloseDisplay(argbSelector.display)
		argbSelector = nil
	}
	if fallbackSelector != nil {
		if fallbackSelector.selectorDrawn && fallbackSelector.selectorWin != 0 {
			C.XUnmapWindow(fallbackSelector.display, fallbackSelector.selectorWin)
			C.XDestroyWindow(fallbackSelector.display, fallbackSelector.selectorWin)
		}
		C.XUnmapWindow(fallbackSelector.display, fallbackSelector.inputWin)
		C.XDestroyWindow(fallbackSelector.display, fallbackSelector.inputWin)
		C.XCloseDisplay(fallbackSelector.display)
		fallbackSelector = nil
	}
}

func SelectorOverlayActive() bool {
	return (argbSelector != nil && argbSelector.active) || (fallbackSelector != nil && fallbackSelector.active)
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

// RegisterSelectorState registers the state of a selector control.
// It is used to track the state of the selector control, including its bounds and other properties.
// This is useful for handling selector events and managing the state of the selector control.
//
// Parameters:
//   - componentID: The ID of the component associated with the selector control
//   - state: A pointer to a common.SelectorState struct representing the state of the selector control
func RegisterSelectorState(componentID uintptr, state *common.SelectorState) {
	selectorStateMapMu.Lock()
	defer selectorStateMapMu.Unlock()
	selectorStateMap[componentID] = state
}

// GetSelectorState retrieves the state of a selector control.
// It is used to get the current state of the selector control, including its bounds and other properties.
// This is useful for handling selector events and managing the state of the selector control.
//
// Parameters:
//   - componentID: The ID of the component associated with the selector control
//   - updates: A variadic number of update functions to modify the state
func UpdateSelectorState(componentID uintptr, updates ...common.UpdateSelectorState) {
	selectorStateMapMu.Lock()
	defer selectorStateMapMu.Unlock()
	state, ok := selectorStateMap[componentID]
	if !ok {
		return
	}
	for _, update := range updates {
		update(state)
	}
}

// GetSelectorState retrieves the state of a selector control.
// It is used to get the current state of the selector control, including its bounds and other properties.
//
// Parameters:
//   - componentID: The ID of the component associated with the selector control
//
// Returns:
//   - *common.SelectorState: A pointer to a common.SelectorState struct representing the state of the selector control
func GetSelectorState(componentID uintptr) *common.SelectorState {
	selectorStateMapMu.Lock()
	defer selectorStateMapMu.Unlock()
	return selectorStateMap[componentID]
}

// GetAllSelectorStates retrieves all selector states.
// It is used to get the current states of all selector controls, including their bounds and other properties.
//
// Returns:
//   - []*common.SelectorState: A slice of pointers to common.SelectorState structs representing the states of all selector controls
func GetAllSelectorStates() []*common.SelectorState {
	selectorStateMapMu.Lock()
	defer selectorStateMapMu.Unlock()
	states := make([]*common.SelectorState, 0, len(selectorStateMap))
	for _, state := range selectorStateMap {
		states = append(states, state)
	}
	return states
}
