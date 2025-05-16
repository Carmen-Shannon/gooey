//go:build windows
// +build windows

package wdws

import (
	"unsafe"

	"github.com/Carmen-Shannon/gooey/common"

	"golang.org/x/sys/windows"
)

// SetTextColor sets the text color for the specified device context.
// This function is a wrapper around the Windows API SetTextColor function.
// https://learn.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-settextcolor
//
// Parameters:
//   - deviceContext: Handle to the device context
//   - color: The color to set, specified as a 32-bit RGB value
//
// Returns:
//   - int32: The previous text color, or GDI_ERROR if the function fails
func SetTextColor(deviceContext uintptr, color *common.Color) int32 {
	clr := colorToColorRef(color)
	ret, _, _ := procSetTextColor.Call(deviceContext, uintptr(clr))
	return int32(ret)
}

// SetBkMode sets the background mode for the specified device context.
// This function is a wrapper around the Windows API SetBkMode function.
// https://learn.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-setbkmode
//
// Parameters:
//   - deviceContext: Handle to the device context
//   - mode: The background mode to set (e.g., TRANSPARENT or OPAQUE)
//
// Returns:
//   - int32: The previous background mode, or GDI_ERROR if the function fails
func SetBkMode(deviceContext uintptr, mode int32) int32 {
	ret, _, _ := procSetBkMode.Call(deviceContext, uintptr(mode))
	return int32(ret)
}

// CreateSolidBrush wraps the Win32 CreateSolidBrush function
// https://learn.microsoft.com/en-us/windows/win32/api/gdi32/nf-gdi32-createsolidbrush
//
// Parameters:
//   - color: The color to create the solid brush with
//
// Returns:
//   - windows.Handle: Handle to the created solid brush
func CreateSolidBrush(color *common.Color) windows.Handle {
	clr := colorToColorRef(color)
	brush, _, _ := procCreateSolidBrush.Call(uintptr(clr))
	return windows.Handle(brush)
}

// CreatePen wraps the Win32 CreatePen function
// https://learn.microsoft.com/en-us/windows/win32/api/gdi32/nf-gdi32-createpen
//
// Parameters:
//   - style: The pen style (e.g., PS_SOLID, PS_DASH)
//   - width: The width of the pen in logical units
//   - color: Pointer to a common.Color struct representing the pen color
//
// Returns:
//   - windows.Handle: Handle to the created pen
func CreatePen(style int32, width int32, color *common.Color) windows.Handle {
	clr := colorToColorRef(color)
	h, _, _ := procCreatePen.Call(uintptr(style), uintptr(width), uintptr(clr))
	return windows.Handle(h)
}

// DrawText draws formatted text in the specified rectangle using the specified format options.
// This function is a wrapper around the Windows API DrawTextW function.
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-drawtextw
//
// Parameters:
//   - deviceContext: Handle to the device context
//   - text: The text to draw, specified as a UTF-16 string
//   - rect: Pointer to a rectangle structure that specifies the area in which to draw the text
//   - format: The format options for drawing the text (e.g., DT_SINGLELINE, DT_CENTER)
//
// Returns:
//   - int32: The height of the drawn text, or GDI_ERROR if the function fails
func DrawText(deviceContext uintptr, text string, rect *[4]int32, format uint32) int32 {
	utf16Str, _ := windows.UTF16PtrFromString(text)
	ret, _, _ := procDrawTextW.Call(deviceContext, uintptr(unsafe.Pointer(utf16Str)), uintptr(len([]rune(text))), uintptr(unsafe.Pointer(rect)), uintptr(format))
	return int32(ret)
}

// DrawRectangle draws a rectangle in the specified device context.
// This function is a wrapper around the Windows API Rectangle function.
// https://learn.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-rectangle
//
// Parameters:
//   - deviceContext: Handle to the device context
//   - dimensions: A structure that specifies the coordinates of the rectangle (left, top, right, bottom)
//
// Returns:
//   - bool: true if the rectangle was drawn successfully, false otherwise
func DrawRectangle(deviceContext uintptr, left, top, right, bottom, roundness int32) bool {
	var ret uintptr
	if roundness > 0 {
		ret, _, _ = procRoundRect.Call(deviceContext, uintptr(left), uintptr(top), uintptr(right), uintptr(bottom), uintptr(roundness), uintptr(roundness))
	} else {
		ret, _, _ = procRectangle.Call(deviceContext, uintptr(left), uintptr(top), uintptr(right), uintptr(bottom))
	}
	return ret != 0
}

// DrawEdge draws a 3D edge around a rectangle in the specified device context.
// This function is a wrapper around the Windows API DrawEdge function.
// https://learn.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-drawedge
//
// Parameters:
//   - deviceContext: Handle to the device context
//   - rect: Pointer to a rectangle structure that specifies the area in which to draw the edge
//   - edge: The type of edge to draw (e.g., BDR_SUNKENOUTER)
//   - flags: Additional flags that specify how the edge is drawn (e.g., BF_RECT)
//
// Returns:
//   - bool: true if the edge was drawn successfully, false otherwise
func DrawEdge(hdc uintptr, rect *[4]int32, edge, flags uint32) bool {
	ret, _, _ := procDrawEdge.Call(
		hdc,
		uintptr(unsafe.Pointer(rect)),
		uintptr(edge),
		uintptr(flags),
	)
	return ret != 0
}

func FillRect(hdc uintptr, rect [4]int32, brush uintptr) {
	r := rect
	_, _, _ = procFillRect.Call(hdc, uintptr(unsafe.Pointer(&r)), brush)
}
