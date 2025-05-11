//go:build windows
// +build windows

package wdws

import (
	"gooey/common"
	"unsafe"

	"golang.org/x/sys/windows"
)

// CreateWindow creates a window with the specified parameters using the CreateWindowExW function.
// It takes a variadic list of options to configure the window creation process.
//
// Note: If the following parameters are not set, the function will return an error:
//   - ClassName: The name of the window class (required)
//   - Style: The style of the window (required)
//   - Instance: The handle to the application instance (required)
//
// Parameters:
//   - options: A variadic list of CreateWindowBuilderOpt functions to configure the window creation process
//
// Returns:
//   - windows.Handle: The handle to the created window on success, or 0 if the function fails
//   - error: An error object if the function fails, or nil on success
func CreateWindow(options ...CreateWindowBuilderOpt) (windows.Handle, error) {
	opts := newCwBuilderOpts()
	for _, opt := range options {
		opt(opts)
	}
	if err := opts.validate(); err != nil {
		return 0, err
	}

	ret, _, err := procCreateWindowExW.Call(
		uintptr(opts.ExStyle),
		uintptr(unsafe.Pointer(opts.ClassName)),
		uintptr(unsafe.Pointer(opts.WindowName)),
		uintptr(opts.Style),
		uintptr(opts.X),
		uintptr(opts.Y),
		uintptr(opts.Width),
		uintptr(opts.Height),
		uintptr(opts.Parent),
		uintptr(opts.Menu),
		uintptr(opts.Instance),
		uintptr(opts.Param),
	)
	if ret == 0 {
		return 0, err
	}

	return windows.Handle(ret), nil
}

// RegisterClassExW registers a window class for use in creating windows.
// It takes a pointer to a WdsWndClass structure and returns an atom (class identifier) or an error.
//
// Parameters:
//   - wndClass: Pointer to a WdsWndClass structure that contains the class information
//
// Returns:
//   - uint16: The atom (class identifier) on success, or 0 if the function fails
//   - error: An error object if the function fails, or nil on success
func RegisterClassExW(options ...RegisterClassOpt) (uint16, error) {
	opts := newRegisterClassOpts()
	for _, opt := range options {
		opt(opts) // Apply each option to the opts struct
	}

	menuName, err := windows.UTF16PtrFromString(opts.MenuName)
	if err != nil {
		menuName = nil
	}
	className, err := windows.UTF16PtrFromString(opts.ClassName)
	if err != nil {
		className = nil
	}

	wndClass := &wdsWndClass{
		CbSize:        opts.Size,
		Style:         opts.Style,
		LpfnWndProc:   opts.Procedure,
		CbClsExtra:    opts.ExtraBytesStruct,
		CbWndExtra:    opts.ExtraBytesWindow,
		HInstance:     windows.Handle(opts.InstanceHandle),
		HIcon:         windows.Handle(opts.IconHandle),
		HCursor:       windows.Handle(opts.CursorHandle),
		HbrBackground: windows.Handle(opts.BackgroundHandle),
		LpszMenuName:  menuName,
		LpszClassName: className,
		HIconSm:       windows.Handle(opts.SmallIconHandle),
	}

	ret, _, err := procRegisterClassExW.Call(
		uintptr(unsafe.Pointer(wndClass)),
	)
	if ret == 0 {
		return 0, err // Return error if registration fails
	}

	return uint16(ret), nil // Return the atom (class identifier) on success
}

// ShowWindow changes the visibility and state of a window. It takes a window handle and a command to show the window.
//
// Parameters:
//   - hWnd: The handle to the window to be shown
//   - nCmdShow: The command that specifies how the window should be shown
//
// Returns:
//   - bool: true if the function succeeds, false otherwise
func ShowWindow(hWnd windows.Handle, nCmdShow ShowWindowCmd) bool {
	ret, _, _ := procShowWindow.Call(
		uintptr(hWnd),
		uintptr(nCmdShow),
	)
	return ret != 0
}

// SetWindowPos sets the position and size of a window. It takes a window handle, insertion order, position, size, and flags.
//
// Parameters:
//   - hwnd: The handle to the window to be positioned
//   - hwndInsertAfter: The handle to the window to precede the positioned window in the Z order
//   - x: The new position of the left side of the window
//   - y: The new position of the top side of the window
//   - cx: The new width of the window
//   - cy: The new height of the window
//   - uFlags: The window sizing and positioning flags
//
// Returns:
//   - bool: true if the function succeeds, false otherwise
func SetWindowPos(hwnd uintptr, hwndInsertAfter uintptr, x, y, cx, cy int32, uFlags uint32) bool {
	ret, _, _ := procSetWindowPos.Call(
		hwnd,
		hwndInsertAfter,
		uintptr(x), uintptr(y), uintptr(cx), uintptr(cy),
		uintptr(uFlags),
	)
	return ret != 0
}

// TranslateMessage wraps the Win32 TranslateMessage function
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-translatemessage
//
// Parameters:
//   - msg: Pointer to the message to be translated
//
// Returns:
//   - bool: true if the message was translated, false otherwise
func TranslateMessage(msg *Msg) bool {
	ret, _, _ := procTranslateMessage.Call(uintptr(unsafe.Pointer(msg)))
	return ret != 0
}

// DispatchMessage wraps the Win32 DispatchMessageW function
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-dispatchmessagew
//
// Parameters:
//   - msg: Pointer to the message to be dispatched
//
// Returns:
//   - uintptr: The result of the message processing
func DispatchMessage(msg *Msg) uintptr {
	ret, _, _ := procDispatchMessageW.Call(uintptr(unsafe.Pointer(msg)))
	return ret
}

// SendMessageW wraps the Win32 SendMessageW function
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessagew
//
// Parameters:
//   - hwnd: Handle to the window to which the message is sent
//   - msg: The message to be sent
//   - wParam: Additional message-specific information
//   - lParam: Additional message-specific information
//
// Returns:
//   - uintptr: The result of the message processing
func SendMessageW(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procSendMessageW.Call(hwnd, uintptr(msg), wParam, lParam)
	return ret
}

// InvalidateRect wraps the Win32 InvalidateRect function.
// This function marks a rectangle in the client area of a window as needing to be redrawn.
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-invalidaterect
//
// Parameters:
//   - windowHandle: Handle to the window to be invalidated
//   - rect: Pointer to a rectangle structure that specifies the area to be invalidated
//   - erase: Boolean value that specifies whether the background should be erased
func InvalidateRect(windowHandle windows.Handle, rect *[4]int32, erase bool) error {
	err, _, _ := procInvalidateRect.Call(uintptr(windowHandle), uintptr(unsafe.Pointer(rect)), uintptr(BoolToBit(erase)))
	if err != 0 {
		return windows.GetLastError()
	}
	return nil
}

// SetCursor wraps the Win32 SetCursor function
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setcursor
//
// Parameters:
//   - hCursor: Handle to the cursor to be set
func SetCursor(hCursor windows.Handle) {
	procSetCursor.Call(uintptr(hCursor))
}

// LoadArrowCursor loads the default arrow cursor wrapping the LoadCursorW function
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-loadcursorw
//
// Returns:
//   - windows.Handle: Handle to the loaded cursor
func LoadArrowCursor() windows.Handle {
	h, _, _ := procLoadCursorW.Call(0, uintptr(IDC_ARROW))
	return windows.Handle(h)
}

// LoadIBeamCursor loads the I-beam cursor wrapping the LoadCursorW function
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-loadcursorw
//
// Returns:
//   - windows.Handle: Handle to the loaded cursor
func LoadIBeamCursor() windows.Handle {
	h, _, _ := procLoadCursorW.Call(0, uintptr(IDC_BEAM))
	return windows.Handle(h)
}

// BeginPaint wraps the Win32 BeginPaint function
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-beginpaint
//
// Parameters:
//   - windowHandle: Handle to the window to be painted
//   - paint: Pointer to a Paint structure that receives information about the painting
func BeginPaint(windowHandle windows.Handle, paint *Paint) {
	procBeginPaint.Call(uintptr(windowHandle), uintptr(unsafe.Pointer(paint)))
}

// EndPaint wraps the Win32 EndPaint function
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-endpaint
//
// Parameters:
//   - windowHandle: Handle to the window to be painted
//   - paint: Pointer to a Paint structure that contains information about the painting
func EndPaint(windowHandle windows.Handle, paint *Paint) {
	procEndPaint.Call(uintptr(windowHandle), uintptr(unsafe.Pointer(paint)))
}

// DeleteObject deletes a GDI object and frees the associated resources.
// It is a wrapper around the Windows API DeleteObject function.
// https://learn.microsoft.com/en-us/windows/win32/api/gdi32/nf-gdi32-deleteobject
//
// Parameters:
//   - h: Handle to the GDI object to be deleted
func DeleteObject(h windows.Handle) {
	procDeleteObject.Call(uintptr(h))
}

// SelectObject selects a GDI object into the specified device context (DC).
// It is a wrapper around the Windows API SelectObject function.
// https://learn.microsoft.com/en-us/windows/win32/api/gdi32/nf-gdi32-selectobject
//
// Parameters:
//   - hdc: Handle to the device context (DC)
//   - hgdiobj: Handle to the GDI object to be selected
//
// Returns:
//   - windows.Handle: Handle to the previously selected GDI object
func SelectObject(hdc uintptr, hgdiobj windows.Handle) windows.Handle {
	ret, _, _ := procSelectObject.Call(hdc, uintptr(hgdiobj))
	return windows.Handle(ret)
}

// MeasureText measures the dimensions of a text string using the specified font.
// It is a wrapper around the Windows API GetTextExtentPoint32W function.
//
// Parameters:
//   - hdc: Handle to the device context (DC)
//   - font: Handle to the font to be used for measuring the text
//   - text: The text string to be measured
//
// Returns:
//   - int32: The width of the text
//   - int32: The height of the text
func MeasureText(hdc uintptr, font windows.Handle, text string) (int32, int32) {
	oldFont := SelectObject(hdc, font)
	defer SelectObject(hdc, oldFont)
	utf16Str, _ := windows.UTF16PtrFromString(text)
	var size Point
	procGetTextExtentPoint32W.Call(
		hdc,
		uintptr(unsafe.Pointer(utf16Str)),
		uintptr(len([]rune(text))),
		uintptr(unsafe.Pointer(&size)),
	)
	return size.X, size.Y
}

// FillBackground fills the background of a window with the specified color using the FillRect function.
//
// Parameters:
//   - windowHandle: Handle to the window to be filled
//   - color: The color to fill the background with
//   - p: Pointer to a Paint structure that contains information about the painting
func FillBackground(windowHandle windows.Handle, color *common.Color, p *Paint) {
	brush := CreateSolidBrush(color)
	defer procDeleteObject.Call(uintptr(brush))
	procFillRect.Call(p.Hdc, uintptr(unsafe.Pointer(&p.RcPaint)), uintptr(brush))
}

// GetMessage wraps the Win32 GetMessageW function
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getmessagew
//
// Parameters:
//   - options: A variadic list of GetMessageOpt functions to configure the message retrieval process
//
// Returns:
//   - int32: The message identifier on success, or -1 if the function fails
//   - error: An error object if the function fails, or nil on success
func GetMessage(options ...GetMessageOpt) (int32, error) {
	opts := newGetMessageOpts()
	for _, opt := range options {
		opt(opts)
	}
	ret, _, err := procGetMessageW.Call(
		uintptr(unsafe.Pointer(opts.Message)),
		uintptr(opts.WindowHandle),
		uintptr(opts.MessageFilterMin),
		uintptr(opts.MessageFilterMax),
	)
	if int32(ret) == -1 {
		return -1, err
	}
	return int32(ret), nil
}

// GetMouseState retrieves the current state of the mouse cursor.
// It returns the x and y coordinates of the cursor, as well as the state of the left, right, and middle mouse buttons.
// It uses the GetCursorPos and ScreenToClient functions to obtain the cursor position relative to the window.
// It also uses the GetKeyState function to determine the state of the mouse buttons.
//
// Parameters:
//   - windowHandle: The handle to the window to which the cursor position is relative
//
// Returns:
//   - x: The x coordinate of the cursor
//   - y: The y coordinate of the cursor
//   - leftDown: true if the left mouse button is down, false otherwise
//   - rightDown: true if the right mouse button is down, false otherwise
//   - middleDown: true if the middle mouse button is down, false otherwise
func GetMouseState(windowHandle windows.Handle) (x, y int32, leftDown, rightDown, middleDown bool) {
	var pt Point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	ProcScreenToClient.Call(uintptr(windowHandle), uintptr(unsafe.Pointer(&pt)))

	ld, _, _ := procGetKeyState.Call(VK_LBUTTON)
	rd, _, _ := procGetKeyState.Call(VK_RBUTTON)
	md, _, _ := procGetKeyState.Call(VK_MBUTTON)

	leftDown = (ld & 0x8000) != 0
	rightDown = (rd & 0x8000) != 0
	middleDown = (md & 0x8000) != 0

	return pt.X, pt.Y, leftDown, rightDown, middleDown
}

func GetKeyState(vk int32) int16 {
	ret, _, _ := procGetKeyState.Call(uintptr(vk))
	return int16(ret)
}

// GetDC retrieves the device context (DC) for a specified window.
// It is a wrapper around the Windows API GetDC function.
//
// Parameters:
//   - hwnd: Handle to the window for which the DC is to be retrieved
//
// Returns:
//   - uintptr: Handle to the device context (DC)
func GetDC(hwnd windows.Handle) uintptr {
	hdc, _, _ := procGetDC.Call(uintptr(hwnd))
	return hdc
}

// ReleaseDC releases the device context (DC) for a specified window.
// It is a wrapper around the Windows API ReleaseDC function.
//
// Parameters:
//   - hwnd: Handle to the window for which the DC is to be released
//   - hdc: Handle to the device context (DC) to be released
func ReleaseDC(hwnd windows.Handle, hdc uintptr) {
	procReleaseDC.Call(uintptr(hwnd), hdc)
}

// handlePaint handles the WM_PAINT message for a window.
//
// Parameters:
//   - hwnd: Handle to the window to be painted
func handlePaint(hwnd windows.Handle) {
	p := &Paint{}
	BeginPaint(hwnd, p)
	defer EndPaint(hwnd, p)

	// Get window size
	var rect [4]int32
	copy(rect[:], p.RcPaint[:])
	width := rect[2] - rect[0]
	height := rect[3] - rect[1]

	// Create memory DC and bitmap
	hdc := p.Hdc
	hdcMem, _, _ := procCreateCompatibleDC.Call(hdc)
	hbmMem, _, _ := procCreateCompatibleBitmap.Call(hdc, uintptr(width), uintptr(height))
	old := SelectObject(hdcMem, windows.Handle(hbmMem))

	bgColor := GetWindowColor(uintptr(hwnd))
	brush := CreateSolidBrush(bgColor)
	defer DeleteObject(brush)
	procFillRect.Call(hdcMem, uintptr(unsafe.Pointer(&rect)), uintptr(brush))

	// Draw everything to memory DC
	cb := getDrawCallback(uintptr(hwnd))
	if cb != nil {
		cb(hdcMem)
	}

	// BitBlt memory DC to window DC
	procBitBlt.Call(
		hdc, 0, 0, uintptr(width), uintptr(height),
		hdcMem, 0, 0, OP_SRCCOPY,
	)

	// Cleanup
	SelectObject(hdcMem, old)
	procDeleteObject.Call(hbmMem)
	procDeleteDC.Call(hdcMem)
}

// local helper functions \\

// BoolToBit converts a boolean value to an integer (1 for true, 0 for false)
func BoolToBit(b bool) int {
	if b {
		return 1
	}
	return 0
}

// colorToColorRef converts a common.Color to a Windows COLORREF value
func colorToColorRef(c *common.Color) uint32 {
	return uint32(c.Red) | (uint32(c.Green) << 8) | (uint32(c.Blue) << 16)
}

// caretPosFromClick calculates the caret position based on a mouse click.
// It uses a binary search algorithm to find the closest character position in the text string.
// It takes into account the width of the text and the padding around it.
// It returns the calculated caret position as an int32.
//
// Parameters:
//   - hdc: Handle to the device context (DC)
//   - font: Handle to the font used for measuring text
//   - text: The text string to be measured
//   - inputX: The x coordinate of the text input field
//   - clickX: The x coordinate of the mouse click
//   - padding: The padding around the text input field
//
// Returns:
//   - int32: The calculated caret position
func caretPosFromClick(hdc uintptr, font windows.Handle, text string, inputX, clickX, padding int32) int32 {
	clickOffset := clickX - inputX - padding
	if clickOffset <= 0 {
		return 0
	}
	runes := []rune(text)
	low, high := 0, len(runes)
	for low < high {
		mid := (low + high) / 2
		sub := string(runes[:mid])
		w, _ := MeasureText(hdc, font, sub)
		if w < clickOffset {
			low = mid + 1
		} else {
			high = mid
		}
	}
	if low > 0 && low <= len(runes) {
		prevW, _ := MeasureText(hdc, font, string(runes[:low-1]))
		currW, _ := MeasureText(hdc, font, string(runes[:low]))
		midpoint := (prevW + currW) / 2
		if clickOffset < midpoint {
			return int32(low - 1)
		}
	}
	return int32(low)
}

// SetClipboardText sets the clipboard text to the specified string.
// It uses the OpenClipboard, EmptyClipboard, GlobalAlloc, GlobalLock, GlobalUnlock,
// SetClipboardData, and CloseClipboard functions to manage the clipboard data.
// It allocates memory for the text in UTF-16 format and copies it to the clipboard.
// It also handles the necessary cleanup of resources.
//
// Parameters:
//   - text: The text string to be set in the clipboard
func setClipboardText(text string) {
	r, _, _ := procOpenClipboard.Call(0)
	if r == 0 {
		return
	}
	defer procCloseClipboard.Call()

	procEmptyClipboard.Call()

	utf16, _ := windows.UTF16FromString(text)
	size := len(utf16) * 2

	hMem, _, _ := procGlobalAlloc.Call(GMEM_MOVEABLE, uintptr(size))
	if hMem == 0 {
		return
	}
	ptr, _, _ := procGlobalLock.Call(hMem)
	if ptr == 0 {
		return
	}
	for i, v := range utf16 {
		*(*uint16)(unsafe.Pointer(ptr + uintptr(i*2))) = v
	}
	procGlobalUnlock.Call(hMem)
	procSetClipboardData.Call(CF_UNICODETEXT, hMem)
}

// GetClipboardText retrieves the text from the clipboard.
// It uses the OpenClipboard, CloseClipboard, GetClipboardData, GlobalLock,
// GlobalUnlock, and GlobalFree functions to manage the clipboard data.
// It allocates memory for the text in UTF-16 format and copies it to a string.
// It also handles the necessary cleanup of resources.
//
// Returns:
//   - string: The text string retrieved from the clipboard
func getClipboardText() string {
	r, _, _ := procOpenClipboard.Call(0)
	if r == 0 {
		return ""
	}
	defer procCloseClipboard.Call()

	hMem, _, _ := procGetClipboardData.Call(CF_UNICODETEXT)
	if hMem == 0 {
		return ""
	}
	ptr, _, _ := procGlobalLock.Call(hMem)
	if ptr == 0 {
		return ""
	}
	defer procGlobalUnlock.Call(hMem)

	var length int
	for {
		c := *(*uint16)(unsafe.Pointer(ptr + uintptr(length*2)))
		if c == 0 {
			break
		}
		length++
	}
	buf := make([]uint16, length)
	for i := 0; i < length; i++ {
		buf[i] = *(*uint16)(unsafe.Pointer(ptr + uintptr(i*2)))
	}
	return windows.UTF16ToString(buf)
}
