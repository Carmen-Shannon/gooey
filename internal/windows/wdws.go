//go:build windows
// +build windows

package wdws

import (
	"fmt"
	"github.com/Carmen-Shannon/gooey/common"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	// Callback Maps \\
	drawCallbackMap     = make(map[uintptr]func(hdc uintptr))
	drawCallbackMu      sync.Mutex
	resizingState       = make(map[uintptr]bool)
	resizingStateMu     sync.Mutex
	wdwColorMap         = make(map[uintptr]common.Color)
	wdwColorMapMu       sync.Mutex
	fontCache           = make(map[string]windows.Handle)
	fontCacheMu         sync.Mutex
	customCursorDraw    = false
	customCursorDrawMu  sync.Mutex
	textInputStateMap   = make(map[uintptr]*common.TextInputState)
	textInputStateMapMu sync.Mutex
	buttonBoundsMap     = make(map[uintptr][4]int32)
	buttonBoundsMapMu   sync.Mutex
	buttonCbMap         = make(map[uintptr]map[string]func(any))
	buttonCbMapMu       sync.Mutex

	// DLLs for Windows API functions \\
	user32   = windows.NewLazySystemDLL("user32.dll")
	gdi32    = windows.NewLazySystemDLL("gdi32.dll")
	kernal32 = windows.NewLazySystemDLL("kernel32.dll")

	// User32 functions \\
	procRegisterClassExW = user32.NewProc("RegisterClassExW")
	procCreateWindowExW  = user32.NewProc("CreateWindowExW")
	procShowWindow       = user32.NewProc("ShowWindow")
	procPostQuitMessage  = user32.NewProc("PostQuitMessage")
	procDefWindowProcW   = user32.NewProc("DefWindowProcW")
	procGetMessageW      = user32.NewProc("GetMessageW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
	procDispatchMessageW = user32.NewProc("DispatchMessageW")
	procDestroyWindow    = user32.NewProc("DestroyWindow")
	procSetCursor        = user32.NewProc("SetCursor")
	procLoadCursorW      = user32.NewProc("LoadCursorW")
	procBeginPaint       = user32.NewProc("BeginPaint")
	procEndPaint         = user32.NewProc("EndPaint")
	procInvalidateRect   = user32.NewProc("InvalidateRect")
	procFillRect         = user32.NewProc("FillRect")
	procDrawTextW        = user32.NewProc("DrawTextW")
	procGetCursorPos     = user32.NewProc("GetCursorPos")
	ProcScreenToClient   = user32.NewProc("ScreenToClient")
	procGetKeyState      = user32.NewProc("GetKeyState")
	procDrawEdge         = user32.NewProc("DrawEdge")
	procSetWindowPos     = user32.NewProc("SetWindowPos")
	procSendMessageW     = user32.NewProc("SendMessageW")
	procGetDC            = user32.NewProc("GetDC")
	procReleaseDC        = user32.NewProc("ReleaseDC")
	procOpenClipboard    = user32.NewProc("OpenClipboard")
	procCloseClipboard   = user32.NewProc("CloseClipboard")
	procEmptyClipboard   = user32.NewProc("EmptyClipboard")
	procSetClipboardData = user32.NewProc("SetClipboardData")
	procGetClipboardData = user32.NewProc("GetClipboardData")

	// GDI32 functions \\
	procCreateSolidBrush       = gdi32.NewProc("CreateSolidBrush")
	procSetTextColor           = gdi32.NewProc("SetTextColor")
	procSetBkMode              = gdi32.NewProc("SetBkMode")
	procRectangle              = gdi32.NewProc("Rectangle")
	procDeleteObject           = gdi32.NewProc("DeleteObject")
	procSelectObject           = gdi32.NewProc("SelectObject")
	procCreateFontW            = gdi32.NewProc("CreateFontW")
	procCreateCompatibleDC     = gdi32.NewProc("CreateCompatibleDC")
	procCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	procBitBlt                 = gdi32.NewProc("BitBlt")
	procDeleteDC               = gdi32.NewProc("DeleteDC")
	procGetTextExtentPoint32W  = gdi32.NewProc("GetTextExtentPoint32W")
	procRoundRect              = gdi32.NewProc("RoundRect")
	procCreatePen              = gdi32.NewProc("CreatePen")

	// Kernal32 Functions \\
	procGlobalAlloc  = kernal32.NewProc("GlobalAlloc")
	procGlobalLock   = kernal32.NewProc("GlobalLock")
	procGlobalUnlock = kernal32.NewProc("GlobalUnlock")

	// Caret Ticker \\
	CT   = common.NewCaretTicker()
	HLTR = common.NewHighlighter()
)

// Custom types for enum purposes
type ShowWindowCmd int32

const (
	// Create window constants
	CW_USEDEFAULT = ^int32(0x7FFFFFFF)

	// Show window constants
	SW_HIDE            ShowWindowCmd = 0  // Hides the window
	SW_SHOWNORMAL      ShowWindowCmd = 1  // Activates and displays the window
	SW_SHOWMINIMIZED   ShowWindowCmd = 2  // Activates and displays the window as minimized
	SW_SHOWMAXIMIZED   ShowWindowCmd = 3  // Activates and displays the window as maximized
	SW_SHOWNOACTIVATE  ShowWindowCmd = 4  // Displays the window without activating it
	SW_SHOW            ShowWindowCmd = 5  // Activates and displays the window
	SW_MINIMIZE        ShowWindowCmd = 6  // Minimizes the window
	SW_SHOWMINNOACTIVE ShowWindowCmd = 7  // Displays the window as minimized without activating it
	SW_SHOWNA          ShowWindowCmd = 8  // Displays the window in its current state without activating it
	SW_RESTORE         ShowWindowCmd = 9  // Activates and displays the window, restoring it if minimized or maximized
	SW_SHOWDEFAULT     ShowWindowCmd = 10 // Sets the show state based on the STARTUPINFO structure
	SW_FORCEMINIMIZE   ShowWindowCmd = 11 // Minimizes the window, even if the thread owning the window is not responding

	// Window Styles
	WS_OVERLAPPEDWINDOW = 0xcf0000
	WS_CHILD            = 0x40000000
	WS_VISIBLE          = 0x10000000
	WS_CLIPCHILDREN     = 0x02000000
	WS_EX_CLIENTEDGE    = 0x00000200
	WS_BORDER           = 0x00800000

	// Edit Control Styles
	ES_LEFT        = 0x0000
	ES_CENTER      = 0x0001
	ES_RIGHT       = 0x0002
	ES_AUTOHSCROLL = 0x0080

	// Window Actions
	WM_CLOSE         = 0x0010
	WM_DESTROY       = 0x0002
	WM_SETCURSOR     = 0x0020
	WM_HTCLIENT      = 1
	WM_PAINT         = 0x000F
	WM_SIZE          = 0x0005
	WM_LBUTTONDOWN   = 0x0201
	WM_LBUTTONUP     = 0x0202
	WM_LBUTTONDBCLK  = 0x0203
	WM_ENTERSIZEMOVE = 0x0231
	WM_EXITSIZEMOVE  = 0x0232
	WM_MOUSEMOVE     = 0x0200
	WM_ERASEBKGND    = 0x0014
	WM_GETTEXT       = 0x000D
	WM_GETTEXTLENGTH = 0x000E
	WM_SETTEXT       = 0x000C
	WM_SETFONT       = 0x0030
	WM_ENABLE        = 0x000A
	WM_COMMAND       = 0x0111
	WM_CHAR          = 0x0102
	WM_KEYDOWN       = 0x0100

	// Notification Codes
	EN_CHANGE = 0x0300

	// Cursor Styles
	IDC_ARROW = 32512
	IDC_BEAM  = 32513

	// System Color Indexes
	COLOR_WINDOW        = 5
	COLOR_WINDOWTEXT    = 8
	COLOR_BTNTEXT       = 18
	COLOR_HIGHLIGHT     = 13
	COLOR_HIGHLIGHTTEXT = 14
	COLOR_GRAYTEXT      = 17

	// DrawText Flags
	DT_TOP                  = 0x00000000
	DT_LEFT                 = 0x00000000
	DT_CENTER               = 0x00000001
	DT_RIGHT                = 0x00000002
	DT_VCENTER              = 0x00000004
	DT_BOTTOM               = 0x00000008
	DT_WORDBREAK            = 0x00000010
	DT_SINGLELINE           = 0x00000020
	DT_EXPANDTABS           = 0x00000040
	DT_TABSTOP              = 0x00000080
	DT_NOCLIP               = 0x00000100
	DT_EXTERNALLEADING      = 0x00000200
	DT_CALCRECT             = 0x00000400
	DT_NOPREFIX             = 0x00000800
	DT_INTERNAL             = 0x00001000
	DT_EDITCONTROL          = 0x00002000
	DT_PATH_ELLIPSIS        = 0x00004000
	DT_END_ELLIPSIS         = 0x00008000
	DT_MODIFYSTRING         = 0x00010000
	DT_RTLREADING           = 0x00020000
	DT_WORD_ELLIPSIS        = 0x00040000
	DT_NOFULLWIDTHCHARBREAK = 0x00080000
	DT_HIDEPREFIX           = 0x00100000
	DT_PREFIXONLY           = 0x00200000

	// Background Modes
	BK_TRANSPARENT = 1
	BK_OPAQUE      = 2

	// Operation Codes
	OP_SRCCOPY = 0x00CC0020

	// Stock Object Indexes
	BLACK_BRUSH = 4
	WHITE_BRUSH = 0
	GRAY_BRUSH  = 2

	// Virtual Keys
	VK_LBUTTON = 0x01
	VK_RBUTTON = 0x02
	VK_MBUTTON = 0x04
	VK_BACK    = 0x08
	VK_DELETE  = 0x2E
	VK_CONTROL = 0x11

	// Border Styles
	BD_EDGE_RAISED       = 0x0004
	BD_EDGE_SUNKEN       = 0x0008
	BD_EDGE_SUNKEN_OUTER = 0x0002
	BD_EDGE_SUNKEN_INNER = 0x0008

	// Border Shapes
	BF_RECT = 0x0F

	// Clipboard Functions
	CF_UNICODETEXT = 13

	// Global Memory Flags
	GMEM_MOVEABLE = 0x0002
)

// WdsWndClass is the struct that defines the WNDCLASSEXW model used in the windows API.
// It includes all fields equivalent to the API specifications.
type wdsWndClass struct {
	CbSize        uint32         // Size of the structure in bytes
	Style         uint32         // Class styles
	LpfnWndProc   uintptr        // Pointer to the window procedure
	CbClsExtra    int32          // Extra bytes to allocate following the structure
	CbWndExtra    int32          // Extra bytes to allocate following the window instance
	HInstance     windows.Handle // Handle to the application instance
	HIcon         windows.Handle // Handle to the class icon
	HCursor       windows.Handle // Handle to the class cursor
	HbrBackground windows.Handle // Handle to the class background brush
	LpszMenuName  *uint16        // Pointer to the resource name of the class menu
	LpszClassName *uint16        // Pointer to the class name
	HIconSm       windows.Handle // Handle to the small icon associated with the class
}

// Msg represents a message sent to a window. It contains information about the message, including the window handle, message identifier, and additional parameters.
// It is used in the message loop to process messages sent to the window.
type Msg struct {
	HWnd    windows.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      Point
}

// Point represents a point in 2D space with X and Y coordinates. It is used to specify the location of a point in the client area of a window.
type Point struct {
	X int32
	Y int32
}

// Paint represents the structure used in the BeginPaint and EndPaint functions. It contains information about the device context (DC) and the area to be painted.
// It is used to manage the painting process in a window.
type Paint struct {
	Hdc         uintptr
	FErase      int32
	RcPaint     [4]int32
	FRestore    int32
	FIncUpdate  int32
	RgbReserved [32]byte
}

// WindowProc is the default window procedure for handling messages sent to a window.
// It processes messages such as WM_DESTROY and calls the default window procedure for other messages.
//
// Parameters:
//   - hwnd: The handle to the window
//   - msg: The message identifier
//   - wParam: Additional message-specific information
//   - lParam: Additional message-specific information
//
// Returns:
//   - uintptr: The result of the message processing
func WindowProc(hwnd windows.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_SETCURSOR:
		if !isCustomCursorDraw() && (uint16(lParam&0xFFFF) == WM_HTCLIENT) {
			SetCursor(LoadArrowCursor())
			return 1
		}
	case WM_SIZE:
		_ = InvalidateRect(hwnd, nil, false)
	case WM_ENTERSIZEMOVE:
		SetResizingState(uintptr(hwnd), true)
	case WM_EXITSIZEMOVE:
		SetResizingState(uintptr(hwnd), false)
	case WM_ERASEBKGND:
		return 1
	case WM_CHAR:
		if HLTR.TextInputID != 0 {
			handleTextInputChar(HLTR.TextInputID, rune(wParam))
		}
		return 0
	case WM_KEYDOWN:
		if HLTR.TextInputID != 0 {
			ctrlDown := (uint16(GetKeyState(VK_CONTROL)) & 0x8000) != 0
			switch wParam {
			case VK_BACK:
				handleTextInputBackspace(HLTR.TextInputID)
			case VK_DELETE:
				handleTextInputDelete(HLTR.TextInputID)
			case 'C', 'c':
				if ctrlDown {
					handleTextInputCopy(HLTR.TextInputID)
				}
			case 'X', 'x':
				if ctrlDown {
					handleTextInputCopy(HLTR.TextInputID)
					handleTextInputBackspace(HLTR.TextInputID)
				}
			case 'V', 'v':
				if ctrlDown {
					handleTextInputPaste(HLTR.TextInputID)
				}
			}
		}
		return 0
	case WM_LBUTTONDOWN:
		x := int32(lParam & 0xFFFF)
		y := int32((lParam >> 16) & 0xFFFF)
		HLTR.SuppressSelection = false
		tiId, tiFound := FindTextInputAt(x, y)
		handleTextInputClickCallbacks(tiId, tiFound, hwnd, x)
		btnId, btnFound := FindButtonAt(x, y)
		handleButtonCallbacks(btnId, btnFound, true)
		return 0
	case WM_LBUTTONDBCLK:
		x := int32(lParam & 0xFFFF)
		y := int32((lParam >> 16) & 0xFFFF)
		tiId, tiFound := FindTextInputAt(x, y)
		handleTextInputClickCallbacks(tiId, tiFound, hwnd, x, true)
		btnId, btnFound := FindButtonAt(x, y)
		handleButtonCallbacks(btnId, btnFound, true)
		return 0
	case WM_MOUSEMOVE:
		if HLTR.Active && HLTR.TextInputID != 0 && !HLTR.SuppressSelection {
			x := int32(lParam & 0xFFFF)
			y := int32((lParam >> 16) & 0xFFFF)
			tiId, tiFound := FindTextInputAt(x, y)
			if tiFound && tiId == HLTR.TextInputID {
				updateTextInputSelection(tiId, hwnd, x, "update")
			}
		}
		return 0
	case WM_LBUTTONUP:
		x := int32(lParam & 0xFFFF)
		y := int32((lParam >> 16) & 0xFFFF)
		btnId, btnFound := FindButtonAt(x, y)
		handleButtonCallbacks(btnId, btnFound, false)
		tiId, tiFound := FindTextInputAt(x, y)
		if HLTR.Active && HLTR.TextInputID != 0 && !HLTR.SuppressSelection {
			if tiFound && HLTR.TextInputID == tiId {
				updateTextInputSelection(tiId, hwnd, x, "end")
				handleTextInputCaretCallbacks(tiId)
			}
			HLTR.Active = false
		}
		return 0
	case WM_PAINT:
		handlePaint(hwnd)
		return 0
	case WM_CLOSE:
		procDestroyWindow.Call(uintptr(hwnd))
		return 0
	case WM_DESTROY:
		procPostQuitMessage.Call(0)
		return 0
	}
	ret, _, _ := procDefWindowProcW.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam,
	)

	return ret
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

// getFontCacheKey generates a unique key for the font cache based on the font height and name.
//
// Parameters:
//   - height: The height of the font
//   - fontName: The name of the font
//
// Returns:
//   - string: A unique key for the font cache
func getFontCacheKey(height int32, fontName string) string {
	return fmt.Sprintf("%d|%s", height, fontName)
}

// CreateFont creates a font with the specified height and font name.
// It is a wrapper around the Windows API CreateFontW function.
// https://learn.microsoft.com/en-us/windows/win32/api/gdi32/nf-gdi32-createfontw
//
// Parameters:
//   - height: The height of the font in logical units
//   - fontName: The name of the font to create
//
// Returns:
//   - windows.Handle: Handle to the created font
func CreateFont(height int32, fontName string) windows.Handle {
	key := getFontCacheKey(height, fontName)
	fontCacheMu.Lock()
	defer fontCacheMu.Unlock()
	if hFont, ok := fontCache[key]; ok && hFont != 0 {
		return hFont
	}
	utf16FontName, _ := windows.UTF16PtrFromString(fontName)
	hFont, _, _ := procCreateFontW.Call(
		uintptr(height), 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		uintptr(unsafe.Pointer(utf16FontName)),
	)
	fontCache[key] = windows.Handle(hFont)
	return windows.Handle(hFont)
}

// CleanupFontCache cleans up the font cache by deleting all cached font handles.
// It is called when the application is exiting to free up resources.
func CleanupFontCache() {
	fontCacheMu.Lock()
	defer fontCacheMu.Unlock()
	for _, hFont := range fontCache {
		if hFont != 0 {
			procDeleteObject.Call(uintptr(hFont))
		}
	}
	fontCache = make(map[string]windows.Handle)
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
func isCustomCursorDraw() bool {
	customCursorDrawMu.Lock()
	defer customCursorDrawMu.Unlock()
	return customCursorDraw
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
