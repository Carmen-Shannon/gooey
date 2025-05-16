//go:build windows
// +build windows

package component

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/Carmen-Shannon/gooey/common"
	wdws "github.com/Carmen-Shannon/gooey/internal/windows"
	"golang.org/x/sys/windows"
)

var (
	screenW    = wdws.GetSystemMetrics(0)
	screenH    = wdws.GetSystemMetrics(1)
	parentHwnd windows.Handle

	captureBounds  bool
	startX, startY int32
)

// createSelectorOverlayWindow creates a new overlay window for the selector component.
// It registers a new window class and creates a window with the specified properties.
// The window is layered and transparent, allowing for a custom drawing of the selector.
// The function also sets the parent window for the overlay.
// It returns the handle of the created window.
// If the parent window is not specified, it will use the first available parent window.
// The function also handles the drawing of the selector rectangle and its properties.
//
// Parameters:
//   - state: The selector state containing the properties for the overlay window.
//
// Returns:
//   - uintptr: The handle of the created overlay window.
func createSelectorOverlayWindow(s Selector) uintptr {
	state := s.(*selector).state
	className := "GooeySelectorOverlay"
	wndProc := windows.NewCallback(selectorOverlayProc)
	clsName, _ := windows.UTF16PtrFromString(className)
	hInstance := wdws.GetModuleHandle()
	_, _ = wdws.RegisterClassExW(
		wdws.ClassNameOpt(className),
		wdws.ProcedureOpt(wndProc),
		wdws.InstanceHandleOpt(uintptr(hInstance)),
		wdws.StyleOpt(0),
	)

	if parentHwnd == 0 {
		parentHwnd = getAnyParentHwnd()
	}
	hwnd, err := wdws.CreateWindow(
		wdws.CreateWindowOptClassName(clsName),
		wdws.CreateWindowOptWindowName(clsName),
		wdws.CreateWindowOptStyle(wdws.WS_POPUP),
		wdws.CreateWindowOptExStyle(wdws.WS_EX_LAYERED|wdws.WS_EX_TRANSPARENT|wdws.WS_EX_TOPMOST),
		wdws.CreateWindowOptPosition(0, 0),
		wdws.CreateWindowOptSize(int32(screenW), int32(screenH)),
		wdws.CreateWindowOptInstance(hInstance),
		wdws.CreateWindowOptParent(0),
		wdws.CreateWindowOptParam(s.ID()),
	)
	if err != nil {
		fmt.Println("Error creating overlay window:", err)
		return 0
	}

	wdws.ShowWindow(hwnd, wdws.SW_SHOWNORMAL)
	updateSelectorOverlay(hwnd, state)
	return uintptr(hwnd)
}

// updateSelectorOverlay updates the overlay window for the selector component.
// It draws the selector rectangle with the specified properties such as color and opacity.
// The function uses a memory device context (DC) to draw the overlay and then updates the layered window.
// It handles the drawing of the selector rectangle and its properties.
//
// Parameters:
//   - hwnd: The handle of the overlay window.
//   - state: The selector state containing the properties for the overlay window.
func updateSelectorOverlay(hwnd windows.Handle, state *common.SelectorState) {
	// --- Prepare memory DC and 32bpp DIB section ---
	type Point struct{ X, Y int32 }
	type Size struct{ CX, CY int32 }
	type BlendFunction struct {
		BlendOp             byte
		BlendFlags          byte
		SourceConstantAlpha byte
		AlphaFormat         byte
	}

	width, height := screenW, screenH
	hdcScreen := wdws.GetDC(hwnd)
	defer wdws.ReleaseDC(hwnd, hdcScreen)

	// Setup BITMAPINFO for 32bpp DIB section (top-down)
	type BITMAPINFOHEADER struct {
		Size          uint32
		Width         int32
		Height        int32
		Planes        uint16
		BitCount      uint16
		Compression   uint32
		SizeImage     uint32
		XPelsPerMeter int32
		YPelsPerMeter int32
		ClrUsed       uint32
		ClrImportant  uint32
	}
	type BITMAPINFO struct {
		Header BITMAPINFOHEADER
		Colors [1]uint32
	}
	bi := BITMAPINFO{
		Header: BITMAPINFOHEADER{
			Size:        40,
			Width:       int32(width),
			Height:      -int32(height), // negative for top-down
			Planes:      1,
			BitCount:    32,
			Compression: 0, // BI_RGB
		},
	}

	var bitsPtr *byte
	hdcMem := wdws.CreateCompatibleDC(hdcScreen)
	defer wdws.DeleteObject(windows.Handle(hdcMem))

	hBmp, _, _ := wdws.Gdi32CreateDIBSection(hdcMem, unsafe.Pointer(&bi), 0, &bitsPtr, 0, 0)
	if hBmp == 0 {
		fmt.Println("Failed to create DIB section")
		return
	}
	defer wdws.DeleteObject(hBmp)
	oldBmp := wdws.SelectObject(hdcMem, hBmp)
	defer wdws.SelectObject(hdcMem, oldBmp)

	// --- Draw overlay into memory DC ---
	// Fill transparent background
	wdws.FillRect(hdcMem, [4]int32{0, 0, int32(width), int32(height)}, uintptr(wdws.GetTransparentBrush()))

	// Draw the selector rectangle
	x, y, w, h := state.Bounds.X, state.Bounds.Y, state.Bounds.W, state.Bounds.H
	if w < 0 {
		x += w
		w = -w
	}
	if h < 0 {
		y += h
		h = -h
	}
	if w == 0 || h == 0 {
		// Don't draw if width or height is zero
		return
	}
	brush := wdws.CreateSolidBrush(state.Color)
	defer wdws.DeleteObject(brush)
	wdws.FillRect(hdcMem, [4]int32{x, y, x + w, y + h}, uintptr(brush))

	pen := wdws.CreatePen(wdws.PS_SOLID, 2, state.Color)
	oldPen := wdws.SelectObject(hdcMem, pen)
	wdws.DrawRectangle(hdcMem, x, y, x+w, y+h, 0)
	wdws.SelectObject(hdcMem, oldPen)
	wdws.DeleteObject(pen)

	// --- Setup parameters for UpdateLayeredWindow ---
	ptDst := Point{0, 0}
	sz := Size{int32(width), int32(height)}
	ptSrc := Point{0, 0}
	blend := BlendFunction{
		BlendOp:             0, // AC_SRC_OVER
		BlendFlags:          0,
		SourceConstantAlpha: byte(state.Opacity * 255),
		AlphaFormat:         0, // 0 for no per-pixel alpha, 1 for AC_SRC_ALPHA if you want per-pixel
	}

	ok := wdws.UpdateLayeredWindow(
		hwnd,
		windows.Handle(hdcMem),
		uintptr(unsafe.Pointer(&ptDst)),
		uintptr(unsafe.Pointer(&sz)),
		uintptr(unsafe.Pointer(&ptSrc)),
		0, // color key
		uintptr(unsafe.Pointer(&blend)),
		2, // ULW_ALPHA
	)
	if !ok {
		fmt.Println("UpdateLayeredWindow failed")
	}
}

// getAnyParentHwnd retrieves the first available parent window handle.
// It iterates through the window color map and returns the first handle found.
//
// Returns:
//   - windows.Handle: The handle of the first available parent window.
func getAnyParentHwnd() windows.Handle {
	for hwnd := range wdws.WdwColorMap() {
		return windows.Handle(hwnd)
	}
	return 0
}

// selectorOverlayProc is the window procedure for the selector overlay window.
// It only handles destroying the window, and re-drawing the selector when triggered by toggling the selector's drawing mode to true.
//
// Parameters:
//   - hwnd: The handle of the window.
//   - msg: The message being processed.
//   - wParam: Additional message-specific information.
//   - lParam: Additional message-specific information.
//
// Returns:
//   - uintptr: The result of the message processing.
func selectorOverlayProc(hwnd windows.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case wdws.WM_SETCURSOR:
		wdws.SetCursor(wdws.LoadCrossCursor())
		return 1
	case wdws.WM_PAINT:
		id := wdws.GetWindowLongPtr(hwnd, wdws.GWLP_USERDATA)
		state := wdws.GetSelectorState(id)
		if state != nil {
			updateSelectorOverlay(hwnd, state)
		}
		return wdws.DefWindowProc(hwnd, msg, wParam, lParam)
	case wdws.WM_NCCREATE:
		cs := (*wdws.CREATESTRUCT)(unsafe.Pointer(lParam))
		wdws.SetWindowLongPtr(hwnd, wdws.GWLP_USERDATA, cs.LpCreateParams)
	case wdws.WM_LBUTTONDOWN:
		id := wdws.GetWindowLongPtr(hwnd, wdws.GWLP_USERDATA)
		state := wdws.GetSelectorState(id)
		if state != nil && state.Drawing && state.Blocking {
			// Start capturing bounds
			captureBounds = true
			startX = int32(lParam & 0xFFFF)
			startY = int32((lParam >> 16) & 0xFFFF)
			// Optionally, set mouse capture to this window
			wdws.SetCapture(hwnd)
		}
		return 0
	case wdws.WM_MOUSEMOVE:
		if captureBounds {
			id := wdws.GetWindowLongPtr(hwnd, wdws.GWLP_USERDATA)
			state := wdws.GetSelectorState(id)
			if state != nil {
				curX := int32(lParam & 0xFFFF)
				curY := int32((lParam >> 16) & 0xFFFF)
				x, y := startX, startY
				w, h := curX-startX, curY-startY
				// Support dragging in any direction
				if w < 0 {
					x, w = curX, -w
				}
				if h < 0 {
					y, h = curY, -h
				}
				wdws.UpdateSelectorState(id, common.UpdateSelectorBounds(common.Rect{
					X: x, Y: y, W: w, H: h,
				}))
				// Trigger redraw
				_ = wdws.InvalidateRect(hwnd, nil, false)
			}
		}
		return 0

	case wdws.WM_LBUTTONUP:
		if captureBounds {
			id := wdws.GetWindowLongPtr(hwnd, wdws.GWLP_USERDATA)
			state := wdws.GetSelectorState(id)
			if state != nil {
				curX := int32(lParam & 0xFFFF)
				curY := int32((lParam >> 16) & 0xFFFF)
				x, y := startX, startY
				w, h := curX-startX, curY-startY
				if w < 0 {
					x, w = curX, -w
				}
				if h < 0 {
					y, h = curY, -h
				}
				wdws.UpdateSelectorState(id,
					common.UpdateSelectorBounds(common.Rect{X: x, Y: y, W: w, H: h}),
					common.UpdateSelectorBlocking(false),
					common.UpdateSelectorDrawing(false),
					common.UpdateSelectorVisible(false),
				)
				wdws.SetTransparentStyle(hwnd, true)
				_ = wdws.InvalidateRect(hwnd, nil, false)
				captureBounds = false
				_ = wdws.ReleaseCapture()

				parent := parentHwnd
				mouseX, mouseY, _, _, _ := wdws.GetMouseState(parent)
				pt := wdws.Point{X: mouseX, Y: mouseY}
				_ = wdws.ScreenToClient(parent, &pt)
				lParamParent := uintptr(uint32(pt.X)&0xFFFF | (uint32(pt.Y) << 16))
				_ = wdws.SendMessageW(uintptr(parent), wdws.WM_LBUTTONUP, 0, lParamParent)
			}
		}
		return 0
	case wdws.WM_DESTROY:
		return 0
	case wdws.WM_CLOSE:
		_ = wdws.DestroyWindow(hwnd)
		return 0
	}
	return wdws.DefWindowProc(hwnd, msg, wParam, lParam)
}

// LaunchSelectorOverlayOnThread creates and launches the selector overlay window on a separate thread.
// It ensures that the window is created on the same thread as the GUI.
// The function waits for the window to be created before returning the handle.
// It handles the message loop for the overlay window and processes messages until the window is closed.
// It returns the handle of the created overlay window.
// If the window creation fails, it returns 0.
//
// Parameters:
//   - s: The selector component to be launched.
func LaunchSelectorOverlayOnThread(s Selector) windows.Handle {
	var hwnd windows.Handle
	done := make(chan struct{})
	go func() {
		runtime.LockOSThread()
		hwnd = windows.Handle(createSelectorOverlayWindow(s))
		close(done)
		if hwnd == 0 {
			return // failed to create window
		}
		msg := new(wdws.Msg)
		for {
			ret, err := wdws.GetMessage(wdws.MessageOpt(msg), wdws.WindowHandleOpt(hwnd))
			if ret == 0 || err != nil {
				break // WM_QUIT or error
			}
			_ = wdws.TranslateMessage(msg)
			_ = wdws.DispatchMessage(msg)
		}
	}()
	<-done // Wait for window creation
	return hwnd
}
