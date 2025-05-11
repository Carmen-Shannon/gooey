//go:build windows
// +build windows

package component

import (
	"github.com/Carmen-Shannon/gooey/common"
	wdws "github.com/Carmen-Shannon/gooey/internal/windows"
)

// drawButton handles the drawing of the button component on the windows platform.
// This function uses windows-specific APIs to create the rendered component.
// It handles the button's background color, pressed state, and label rendering.
// It also handles the button's roundness and label font size.
//
// Parameters:
//   - hdc: Handle to the device context where the button will be drawn.
//   - b: The Button component to be drawn.
func drawButton(hdc uintptr, b Button) {
	x, y := b.Position()
	w, h := b.Size()

	radius := (b.Roundness() * min(w, h)) / 200

	var rect = [4]int32{x, y, x + w, y + h}
	var color *common.Color
	if !b.Enabled() {
		color = b.BackgroundColorDisabled()
	} else if b.Pressed() {
		color = b.BackgroundColorPressed()
	} else if b.Hovered() {
		color = b.BackgroundColorHover()
	} else {
		color = b.BackgroundColor()
	}
	brush := wdws.CreateSolidBrush(color)
	defer wdws.DeleteObject(brush)

	oldBrush := wdws.SelectObject(hdc, brush)
	defer wdws.SelectObject(hdc, oldBrush)

	_ = wdws.DrawRectangle(hdc, x, y, x+w, y+h, radius)

	fontName := b.LabelFont()
	fontSize := b.LabelSize()
	hFont := wdws.CreateFont(-fontSize, fontName)
	if hFont != 0 {
		oldFont := wdws.SelectObject(hdc, hFont)
		defer func() {
			wdws.SelectObject(hdc, oldFont)
			wdws.DeleteObject(hFont)
		}()
	}

	wdws.SetTextColor(hdc, b.LabelColor())
	wdws.SetBkMode(hdc, wdws.BK_TRANSPARENT)
	wdws.DrawText(hdc, b.Label(), &rect, wdws.DT_CENTER|wdws.DT_VCENTER|wdws.DT_SINGLELINE)
}

// mapBtnCb maps a button callback to a specific button ID.
// This function registers the button's callback function in the Windows API.
//
// Parameters:
//   - b: The Button component to map the callback for.
//   - cbMap: A map of callback functions for the button.
func mapBtnCb(b Button, cbMap map[string]func(any)) {
    wdws.RegisterButtonCallback(b.ID(), cbMap)
}

// registerBtnBounds registers the button's bounds in the Windows API.
// This function registers the button's position and size in the Windows API.
//
// Parameters:
//   - b: The Button component to register the bounds for.
func registerBtnBounds(b Button) {
    x, y := b.Position()
    w, h := b.Size()
    wdws.RegisterButtonBounds(b.ID(), [4]int32{x, y, w, h})
}
