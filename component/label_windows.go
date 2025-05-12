//go:build windows
// +build windows

package component

import (
	"github.com/Carmen-Shannon/gooey/common"
	wdws "github.com/Carmen-Shannon/gooey/internal/windows"
)

// drawLabel handles the drawing of the label component on the windows platform.
// This function uses windows-specific APIs to create the rendered component and handles shrinking the font size to fit the bounds of the component.
// It also handles word wrapping and text alignment.
//
// Parameters:
//   - hdc: Handle to the device context where the label will be drawn.
//   - l: The Label component to be drawn.
func drawLabel(ctx *common.DrawCtx, l Label) {
	x, y := l.Position()
	w, h := l.Size()
	var rect = [4]int32{x, y, x + w, y + h}

	fontName := l.Font()
	fontSize := l.TextSize()
	text := l.Text()

	// Only shrink font size if text is too wide and word wrap is off
	if !l.WordWrap() {
		hFont := wdws.CreateFont(-fontSize, fontName)
		textWidth, _ := wdws.MeasureText(ctx.Hdc, hFont, text)
		for textWidth > w && fontSize > 12 {
			fontSize--
			hFont = wdws.CreateFont(-fontSize, fontName)
			textWidth, _ = wdws.MeasureText(ctx.Hdc, hFont, text)
		}
	}

	hFont := wdws.CreateFont(-fontSize, fontName)
	if hFont != 0 {
		oldFont := wdws.SelectObject(ctx.Hdc, hFont)
		defer func() {
			wdws.SelectObject(ctx.Hdc, oldFont)
		}()
	}

	wdws.SetTextColor(ctx.Hdc, l.Color())
	wdws.SetBkMode(ctx.Hdc, wdws.BK_TRANSPARENT)

	var alignment uint32
	switch l.TextAlignment() {
	case LeftAlign:
		alignment = wdws.DT_LEFT
	case CenterAlign:
		alignment = wdws.DT_CENTER
	case RightAlign:
		alignment = wdws.DT_RIGHT
	}

	format := alignment
	if l.WordWrap() {
		format |= wdws.DT_WORDBREAK
		// Only expand height if needed, never shrink
		calcRect := rect
		wdws.DrawText(ctx.Hdc, text, &calcRect, format|wdws.DT_CALCRECT)
		newHeight := calcRect[3] - calcRect[1]
		if newHeight > h {
			if lbl, ok := l.(*label); ok {
				lbl.SetSize(w, newHeight)
				rect[3] = rect[1] + newHeight
			}
		}
	} else {
		format |= wdws.DT_VCENTER | wdws.DT_SINGLELINE
	}

	wdws.DrawText(ctx.Hdc, text, &rect, format)
}
