//go:build linux
// +build linux

package component

import (
	"github.com/Carmen-Shannon/gooey/common"
	"github.com/Carmen-Shannon/gooey/internal/linux"
)

func drawButton(ctx *common.DrawCtx, b Button) {
	x, y := b.Position()
	w, h := b.Size()
    if w <= 0 || h <= 0 {
        return
    }
	radius := (b.Roundness() * min(w, h)) / 200

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

	display := linux.GetDisplay(ctx.Hwnd)
	if display == nil {
		return
	}
	drawable := linux.C_Drawable(ctx.Hdc)

	// Fill background (rounded rectangle if radius > 0)
	if radius > 0 {
		linux.XFillRoundedRect(display, drawable, int(x), int(y), int(w), int(h), int(radius), color)
	} else {
		linux.XFillRect(display, drawable, int(x), int(y), int(w), int(h), color)
	}

	// Draw label
	fontName := b.LabelFont()
	fontSize := b.LabelSize()
	label := b.Label()
	labelColor := b.LabelColor()

	linux.XDrawTextCentered(display, drawable, int(x), int(y), int(w), int(h), fontName, int(fontSize), label, labelColor)
}

// mapBtnCb maps a button callback to a specific button ID.
// This function registers the button's callback function in the Windows API.
//
// Parameters:
//   - b: The Button component to map the callback for.
//   - cbMap: A map of callback functions for the button.
func mapBtnCb(b Button, cbMap map[string]func(any)) {
	linux.RegisterButtonCallback(b.ID(), cbMap)
}

// registerBtnBounds registers the button's bounds in the Windows API.
// This function registers the button's position and size in the Windows API.
//
// Parameters:
//   - b: The Button component to register the bounds for.
func registerBtnBounds(b Button) {
	x, y := b.Position()
	w, h := b.Size()
	linux.RegisterButtonBounds(b.ID(), [4]int32{x, y, w, h})
}
