//go:build windows
// +build windows

package wdws

import (
	"fmt"
	"unsafe"

	"github.com/Carmen-Shannon/gooey/common"

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
	_, _, _ = procSetCursor.Call(uintptr(hCursor))
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

// LoadCrossCursor loads the crosshair cursor wrapping the LoadCursorW function
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-loadcursorw
//
// Returns:
//   - windows.Handle: Handle to the loaded cursor
func LoadCrossCursor() windows.Handle {
	h, _, _ := procLoadCursorW.Call(0, uintptr(IDC_CROSS))
	return windows.Handle(h)
}

// BeginPaint wraps the Win32 BeginPaint function
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-beginpaint
//
// Parameters:
//   - windowHandle: Handle to the window to be painted
//   - paint: Pointer to a Paint structure that receives information about the painting
func BeginPaint(windowHandle windows.Handle, paint *Paint) {
	_, _, _ = procBeginPaint.Call(uintptr(windowHandle), uintptr(unsafe.Pointer(paint)))
}

// EndPaint wraps the Win32 EndPaint function
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-endpaint
//
// Parameters:
//   - windowHandle: Handle to the window to be painted
//   - paint: Pointer to a Paint structure that contains information about the painting
func EndPaint(windowHandle windows.Handle, paint *Paint) {
	_, _, _ = procEndPaint.Call(uintptr(windowHandle), uintptr(unsafe.Pointer(paint)))
}

// DeleteObject deletes a GDI object and frees the associated resources.
// It is a wrapper around the Windows API DeleteObject function.
// https://learn.microsoft.com/en-us/windows/win32/api/gdi32/nf-gdi32-deleteobject
//
// Parameters:
//   - h: Handle to the GDI object to be deleted
func DeleteObject(h windows.Handle) {
	_, _, _ = procDeleteObject.Call(uintptr(h))
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
	_, _, _ = procGetTextExtentPoint32W.Call(
		hdc,
		uintptr(unsafe.Pointer(utf16Str)),
		uintptr(len([]rune(text))),
		uintptr(unsafe.Pointer(&size)),
	)
	return size.X, size.Y
}

// DefWindowProc wraps the Win32 DefWindowProcW function
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-defwindowprocw
//
// Parameters:
//   - hwnd: Handle to the window
//   - msg: The message to be processed
//   - wParam: Additional message-specific information
//   - lParam: Additional message-specific information
//
// Returns:
//   - uintptr: The result of the message processing
func DefWindowProc(hwnd windows.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procDefWindowProcW.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam,
	)
	return ret
}

// FillBackground fills the background of a window with the specified color using the FillRect function.
//
// Parameters:
//   - windowHandle: Handle to the window to be filled
//   - color: The color to fill the background with
//   - p: Pointer to a Paint structure that contains information about the painting
func FillBackground(windowHandle windows.Handle, color *common.Color, p *Paint) {
	brush := CreateSolidBrush(color)
	defer DeleteObject(brush)
	FillRect(p.Hdc, p.RcPaint, uintptr(brush))
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
	_, _, _ = procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	succ := ScreenToClient(windowHandle, &pt)
	if !succ {
		fmt.Println("ScreenToClient failed")
		return 0, 0, false, false, false
	}

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
	_, _, _ = procReleaseDC.Call(uintptr(hwnd), hdc)
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
	hdcMem := CreateCompatibleDC(hdc)
	hbmMem, _, _ := procCreateCompatibleBitmap.Call(hdc, uintptr(width), uintptr(height))
	old := SelectObject(hdcMem, windows.Handle(hbmMem))

	bgColor := GetWindowColor(uintptr(hwnd))
	brush := CreateSolidBrush(bgColor)
	defer DeleteObject(brush)
	FillRect(hdcMem, rect, uintptr(brush))

	// Draw everything to memory DC
	cb := getDrawCallback(uintptr(hwnd))
	if cb != nil {
		cb(hdcMem)
	}

	// BitBlt memory DC to window DC
	_, _, _ = procBitBlt.Call(
		hdc, 0, 0, uintptr(width), uintptr(height),
		hdcMem, 0, 0, OP_SRCCOPY,
	)

	// Cleanup
	SelectObject(hdcMem, old)
	DeleteObject(windows.Handle(hbmMem))
	DeleteDC(hdcMem)
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
	defer func() {
		_, _, _ = procCloseClipboard.Call()
	}()

	_, _, _ = procEmptyClipboard.Call()

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
	_, _, _ = procGlobalUnlock.Call(hMem)
	_, _, _ = procSetClipboardData.Call(CF_UNICODETEXT, hMem)
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
	defer func() {
		_, _, _ = procCloseClipboard.Call()
	}()

	hMem, _, _ := procGetClipboardData.Call(CF_UNICODETEXT)
	if hMem == 0 {
		return ""
	}
	ptr, _, _ := procGlobalLock.Call(hMem)
	if ptr == 0 {
		return ""
	}
	defer func() {
		_, _, _ = procGlobalUnlock.Call(hMem)
	}()

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

// CreateCompatibleDC creates a memory device context compatible with the specified device.
//
// Parameters:
//   - hdc: Handle to the device context (DC) to be used as a reference
//
// Returns:
//   - uintptr: Handle to the created memory device context (DC)
func CreateCompatibleDC(hdc uintptr) uintptr {
	ret, _, _ := procCreateCompatibleDC.Call(hdc)
	return ret
}

// GetSystemMetrics retrieves the specified system metric or system configuration setting.
// It uses the GetSystemMetrics function from the Windows API to obtain the value.
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getsystemmetrics
//
// Parameters:
//   - index: The system metric or configuration setting to retrieve
//
// Returns:
//   - int32: The value of the specified system metric or configuration setting
func GetSystemMetrics(index int32) int32 {
	ret, _, _ := procGetSystemMetrics.Call(uintptr(index))
	return int32(ret)
}

// LOWORD is a helper function that extracts the low-order word from a 32-bit value.
// It is used to retrieve the x-coordinate from a LPARAM value.
//
// Parameters:
//   - lparam: The 32-bit value from which to extract the low-order word
//
// Returns:
//   - uint16: The low-order word (x-coordinate)
func LOWORD(lparam uintptr) uint16 {
	return uint16(lparam & 0xFFFF)
}

// HIWORD is a helper function that extracts the high-order word from a 32-bit value.
// It is used to retrieve the y-coordinate from a LPARAM value.
//
// Parameters:
//   - lparam: The 32-bit value from which to extract the high-order word
//
// Returns:
//   - uint16: The high-order word (y-coordinate)
func HIWORD(lparam uintptr) uint16 {
	return uint16((lparam >> 16) & 0xFFFF)
}

// DestroyWindow destroys the specified window and frees the associated resources.
// It uses the DestroyWindow function from the Windows API to destroy the window.
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-destroywindow
//
// Parameters:
//   - hwnd: Handle to the window to be destroyed
//
// Returns:
//   - bool: true if the function succeeds, false otherwise
func DestroyWindow(hwnd windows.Handle) bool {
	ret, _, _ := procDestroyWindow.Call(uintptr(hwnd))
	return ret != 0
}

// GetModuleHandle retrieves a handle to the specified module.
// It uses the GetModuleHandleW function from the Windows API to obtain the handle.
// https://learn.microsoft.com/en-us/windows/win32/api/libloaderapi/nf-libloaderapi-getmodulehandlew
//
// Returns:
//   - windows.Handle: Handle to the specified module
func GetModuleHandle() windows.Handle {
	ret, _, _ := procGetModuleHandle.Call(0)
	return windows.Handle(ret)
}

// GetTransparentBrush retrieves a handle to a transparent brush.
// It creates a solid black brush (RGB(0,0,0)), which is treated as transparent for layered windows.
// It uses the CreateSolidBrush function from the Windows API to create the brush.
//
// Returns:
//   - windows.Handle: Handle to the transparent brush
func GetTransparentBrush() windows.Handle {
	if transparentBrush == 0 {
		// Create a solid black brush (RGB(0,0,0)), which is treated as transparent for layered windows
		transparentBrush = CreateSolidBrush(common.ColorBlack)
	}
	return transparentBrush
}

// UpdateLayeredWindow updates the layered window attributes for a specified window.
// It uses the UpdateLayeredWindow function from the Windows API to perform the update.
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-updatelayeredwindow
//
// Parameters:
//   - hwnd: Handle to the window to be updated
//   - hdcSrc: Handle to the source device context (DC)
//   - pptDst: Pointer to a POINT structure that specifies the destination point
//   - psize: Pointer to a SIZE structure that specifies the size of the destination rectangle
//   - pptSrc: Pointer to a POINT structure that specifies the source point
//   - crKey: The color key for transparency (RGB value)
//   - blendFunc: Pointer to a BLENDFUNCTION structure that specifies the blending function
//   - dwFlags: The flags that specify the layered window attributes
//
// Returns:
//   - bool: true if the function succeeds, false otherwise
func UpdateLayeredWindow(
	hwnd windows.Handle,
	hdcSrc windows.Handle, // <-- changed from uintptr to windows.Handle
	pptDst, psize, pptSrc uintptr,
	crKey uint32,
	blendFunc uintptr,
	dwFlags uint32,
) bool {
	ret, _, _ := procUpdateLayeredWindow.Call(
		uintptr(hwnd),
		0, // hdcDst (screen DC), 0 means use default
		pptDst,
		psize,
		uintptr(hdcSrc), // <-- pass as windows.Handle
		pptSrc,
		uintptr(crKey),
		blendFunc,
		uintptr(dwFlags),
	)
	return ret != 0
}

// Gdi32CreateDIBSection creates a DIB section compatible with the specified device context.
// It uses the CreateDIBSection function from the GDI32 library to create the DIB section.
// https://learn.microsoft.com/en-us/windows/win32/api/gdi32/nf-gdi32-createdibsection
//
// Parameters:
//   - hdc: Handle to the device context (DC)
//   - pbmi: Pointer to a BITMAPINFO structure that specifies the DIB format
//   - usage: The DIB color space type (DIB_RGB_COLORS or DIB_PAL_COLORS)
//   - bits: Pointer to a pointer that receives the address of the DIB section
//   - hSection: Handle to a file mapping object (optional)
//   - offset: Offset to the DIB section in the file mapping object (optional)
//
// Returns:
//   - windows.Handle: Handle to the DIB section
//   - uintptr: Pointer to the DIB section bits
//   - uintptr: Offset to the DIB section in the file mapping object
func Gdi32CreateDIBSection(hdc uintptr, pbmi unsafe.Pointer, usage uint32, bits **byte, hSection uintptr, offset uint32) (windows.Handle, uintptr, error) {
	ret, bitsPtr, lastErr := procCreateDIBSection.Call(
		hdc,
		uintptr(pbmi),
		uintptr(usage),
		uintptr(unsafe.Pointer(bits)),
		hSection,
		uintptr(offset),
	)
	return windows.Handle(ret), bitsPtr, lastErr
}

// GetWindowLong retrieves information about the specified window.
// It is a wrapper around the Windows API GetWindowLong function.
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getwindowlongw
//
// Parameters:
//   - hwnd: Handle to the window whose information is to be retrieved.
//   - index: The zero-based offset to the value to be retrieved. Use constants like GWL_EXSTYLE, GWL_STYLE, etc.
//
// Returns:
//   - int32: The value retrieved from the specified offset for the window.
func GetWindowLong(hwnd windows.Handle, index int32) int32 {
	ret, _, _ := procGetWindowLong.Call(uintptr(hwnd), uintptr(index))
	return int32(ret)
}

// SetWindowLong changes an attribute of the specified window.
// It is a wrapper around the Windows API SetWindowLong function.
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowlongw
//
// Parameters:
//   - hwnd: Handle to the window whose attribute is to be changed.
//   - index: The zero-based offset to the value to be set. Use constants like GWL_EXSTYLE, GWL_STYLE, etc.
//   - newLong: The replacement value.
//
// Returns:
//   - int32: The previous value of the specified offset for the window.
func SetWindowLong(hwnd windows.Handle, index int32, newLong int32) int32 {
	ret, _, _ := procSetWindowLong.Call(uintptr(hwnd), uintptr(index), uintptr(newLong))
	return int32(ret)
}

// SetForegroundWindow brings the specified window to the foreground and activates it.
// It is a wrapper around the Windows API SetForegroundWindow function.
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setforegroundwindow
//
// Parameters:
//   - hwnd: Handle to the window to be brought to the foreground.
//
// Returns:
//   - bool: true if the window was brought to the foreground, false otherwise.
func SetForegroundWindow(hwnd windows.Handle) bool {
	ret, _, _ := procSetForegroundWindow.Call(uintptr(hwnd))
	return ret != 0
}

// SetCapture sets mouse capture to the specified window.
// All mouse input is directed to this window until ReleaseCapture is called or the window is destroyed.
// It is a wrapper around the Windows API SetCapture function.
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setcapture
//
// Parameters:
//   - hwnd: Handle to the window to capture mouse input.
//
// Returns:
//   - bool: true if the function succeeds, false otherwise.
func SetCapture(hwnd windows.Handle) bool {
	ret, _, _ := procSetCapture.Call(uintptr(hwnd))
	return ret != 0
}

// ReleaseCapture releases the mouse capture from a window in the current thread.
// It is a wrapper around the Windows API ReleaseCapture function.
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-releasecapture
//
// Returns:
//   - bool: true if the function succeeds, false otherwise.
func ReleaseCapture() bool {
	ret, _, _ := procReleaseCapture.Call()
	return ret != 0
}

// DeleteDC deletes the specified device context (DC) and frees associated resources.
// It is a wrapper around the Windows API DeleteDC function.
// https://learn.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-deletedc
//
// Parameters:
//   - hdc: Handle to the device context to be deleted.
func DeleteDC(hdc uintptr) {
	_, _, _ = procDeleteDC.Call(hdc)
}

// GetWindowLongPtr retrieves information about the specified window. (Pointer-sized version.)
// It is a wrapper around the Windows API GetWindowLongPtr function.
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getwindowlongptrw
//
// Parameters:
//   - hwnd: Handle to the window whose information is to be retrieved.
//   - index: The zero-based offset to the value to be retrieved.
//
// Returns:
//   - uintptr: The value retrieved from the specified offset for the window.
func GetWindowLongPtr(hwnd windows.Handle, index int32) uintptr {
	ret, _, _ := procGetWindowLongPtr.Call(uintptr(hwnd), uintptr(index))
	return ret
}

// SetWindowLongPtr changes an attribute of the specified window. (Pointer-sized version.)
// It is a wrapper around the Windows API SetWindowLongPtr function.
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowlongptrw
//
// Parameters:
//   - hwnd: Handle to the window whose attribute is to be changed.
//   - index: The zero-based offset to the value to be set.
//   - value: The replacement value.
//
// Returns:
//   - uintptr: The previous value of the specified offset for the window.
func SetWindowLongPtr(hwnd windows.Handle, index int32, value uintptr) uintptr {
	ret, _, _ := procSetWindowLongPtr.Call(uintptr(hwnd), uintptr(index), value)
	return ret
}

// SetTransparentStyle sets or clears WS_EX_TRANSPARENT on a window.
// If transparent is true, sets WS_EX_TRANSPARENT (makes window click-through).
// If transparent is false, removes WS_EX_TRANSPARENT (makes window interactive).
//
// Parameters:
//   - hwnd: Handle to the window to modify.
//   - transparent: true to enable click-through, false to make interactive.
func SetTransparentStyle(hwnd windows.Handle, transparent bool) {
	style := GetWindowLong(hwnd, GWL_EXSTYLE)
	if transparent {
		style |= WS_EX_TRANSPARENT
	} else {
		style &^= WS_EX_TRANSPARENT
	}
	SetWindowLong(hwnd, GWL_EXSTYLE, style)
}

// IsTransparentStyle returns true if WS_EX_TRANSPARENT is set on the window.
//
// Parameters:
//   - hwnd: Handle to the window to check.
//
// Returns:
//   - bool: true if WS_EX_TRANSPARENT is set, false otherwise.
func IsTransparentStyle(hwnd windows.Handle) bool {
	style := GetWindowLong(hwnd, GWL_EXSTYLE)
	return style&WS_EX_TRANSPARENT != 0
}

// GetAncestor retrieves the ancestor (root or parent) window handle.
// It is a wrapper around the Windows API GetAncestor function.
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getancestor
//
// Parameters:
//   - hwnd: Handle to the window whose ancestor is to be retrieved.
//   - gaFlags: The ancestor to be retrieved. Use GA_PARENT (1), GA_ROOT (2), or GA_ROOTOWNER (3).
//
// Returns:
//   - windows.Handle: Handle to the ancestor window.
func GetAncestor(hwnd windows.Handle, gaFlags uint32) windows.Handle {
	ret, _, _ := procGetAncestor.Call(uintptr(hwnd), uintptr(gaFlags))
	return windows.Handle(ret)
}

// ScreenToClient converts screen coordinates to client-area coordinates for the given window.
// It is a wrapper around the Windows API ScreenToClient function.
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-screentoclient
//
// Parameters:
//   - hwnd: Handle to the window whose client area is used for conversion.
//   - pt: Pointer to a Point struct (with X and Y as int32) containing the screen coordinates.
//     On return, the struct is updated with client-area coordinates.
//
// Returns:
//   - bool: true if the function succeeds, false otherwise.
func ScreenToClient(hwnd windows.Handle, pt *Point) bool {
	ret, _, _ := procScreenToClient.Call(uintptr(hwnd), uintptr(unsafe.Pointer(pt)))
	return ret != 0
}
