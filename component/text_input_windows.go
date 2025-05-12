//go:build windows
// +build windows

package component

import (
	"github.com/Carmen-Shannon/gooey/common"
	wdws "github.com/Carmen-Shannon/gooey/internal/windows"
)

// drawTextInput handles drawing the text input component on the windows platform.
// This function uses windows-specific APIs to create the rendered component.
// It handles highlighting, and caret rendering as well.
//
// Parameters:
//   - hdc: Handle to the device context where the text input will be drawn.
//   - ti: The TextInput component to be drawn.
func drawTextInput(ctx *common.DrawCtx, ti TextInput) {
	x, y := ti.Position()
	w, h := ti.Size()

	// Draw background
	bg := ti.Color()
	brush := wdws.CreateSolidBrush(bg)
	oldBrush := wdws.SelectObject(ctx.Hdc, brush)
	wdws.DrawRectangle(ctx.Hdc, x, y, x+w, y+h, 4)
	wdws.SelectObject(ctx.Hdc, oldBrush)
	wdws.DeleteObject(brush)

	// Draw border (optional)
	wdws.DrawEdge(ctx.Hdc, &[4]int32{x, y, x + w, y + h}, wdws.BD_EDGE_SUNKEN, wdws.BF_RECT)

	// Draw selection highlight if any
	text := ti.Value()
	font := wdws.CreateFont(-ti.TextSize(), ti.Font())
	oldFont := wdws.SelectObject(ctx.Hdc, font)
	defer func() {
		wdws.SelectObject(ctx.Hdc, oldFont)
		wdws.DeleteObject(font)
	}()

	wdws.SetBkMode(ctx.Hdc, wdws.BK_TRANSPARENT)
	wdws.SetTextColor(ctx.Hdc, ti.TextColor())

	selStart, selEnd := ti.Selection()
	if selStart > selEnd {
		selStart, selEnd = selEnd, selStart
	}
	if selStart != selEnd && selStart >= 0 && selEnd <= int32(len(text)) {
		prefix := text[:selStart]
		highlight := text[selStart:selEnd]

		prefixWidth, _ := wdws.MeasureText(ctx.Hdc, font, prefix)
		highlightWidth, _ := wdws.MeasureText(ctx.Hdc, font, highlight)

		highlightRect := [4]int32{
			x + 4 + prefixWidth,
			y + 2,
			x + 4 + prefixWidth + highlightWidth,
			y + h - 2,
		}
		highlightBrush := wdws.CreateSolidBrush(&common.Color{Red: 120, Green: 160, Blue: 240})
		oldHighlightBrush := wdws.SelectObject(ctx.Hdc, highlightBrush)
		wdws.FillRect(ctx.Hdc, highlightRect, uintptr(highlightBrush))
		wdws.SelectObject(ctx.Hdc, oldHighlightBrush)
		wdws.DeleteObject(highlightBrush)
	}

	// Draw text
	textRect := [4]int32{x + 4, y + 2, x + w - 4, y + h - 2}
	wdws.DrawText(ctx.Hdc, text, &textRect, wdws.DT_LEFT|wdws.DT_VCENTER|wdws.DT_SINGLELINE)

	// Draw caret if focused
	if ti.Focused() && wdws.CT.Visible && selStart == selEnd {
		caretPos := int(ti.Caret())
		runes := []rune(text)
		if caretPos < 0 {
			caretPos = 0
		}
		if caretPos > len(runes) {
			caretPos = len(runes)
		}
		caretX := x + 4
		if caretPos > 0 {
			caretWidth, _ := wdws.MeasureText(ctx.Hdc, font, string(runes[:caretPos]))
			caretX += caretWidth
		}
		caretHeight := int32(float32(h) * 0.5)
		caretY := y + (h-caretHeight)/2
		wdws.DrawRectangle(ctx.Hdc, caretX, caretY, caretX+2, caretY+caretHeight, 0)
	}
}

// registerTextInput registers the text input component with the windows package so the state can be tracked and updated while rendering or responding to events.
//
// Parameters:
//   - ti: The TextInput component to be registered.
func registerTextInput(ti TextInput) {
	wdws.RegisterTextInputState(ti.ID(), ti.(*textInput).state)
}
