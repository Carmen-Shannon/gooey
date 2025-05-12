//go:build linux
// +build linux

package component

import (
	"github.com/Carmen-Shannon/gooey/common"
	"github.com/Carmen-Shannon/gooey/internal/linux"
)

func drawTextInput(ctx *common.DrawCtx, ti TextInput) {
	x, y := ti.Position()
	w, h := ti.Size()
	if w <= 0 || h <= 0 {
		return
	}

	display := linux.GetDisplay(ctx.Hwnd)
	if display == nil {
		return
	}
	drawable := linux.C_Drawable(ctx.Hdc)

	// Draw background
	bg := ti.Color()
	linux.XFillRect(display, drawable, int(x), int(y), int(w), int(h), bg)

	// Draw border (simple sunken effect)
	borderColor := &common.Color{Red: 180, Green: 180, Blue: 180}
	linux.XDrawRect(display, drawable, int(x), int(y), int(w), int(h), borderColor)

	// Draw selection highlight if any
	text := ti.Value()
	fontName := ti.Font()
	fontSize := int(ti.TextSize())
	textColor := ti.TextColor()

	selStart, selEnd := ti.Selection()
	if selStart > selEnd {
		selStart, selEnd = selEnd, selStart
	}
	if selStart != selEnd && selStart >= 0 && selEnd <= int32(len([]rune(text))) {
		prefix := []rune(text)[:selStart]
		highlight := []rune(text)[selStart:selEnd]

		prefixStr := string(prefix)
		highlightStr := string(highlight)

		prefixWidth := linux.XTextWidth(display, drawable, fontName, fontSize, prefixStr)
		highlightWidth := linux.XTextWidth(display, drawable, fontName, fontSize, highlightStr)

		highlightColor := &common.Color{Red: 120, Green: 160, Blue: 240}
		linux.XFillRect(display, drawable, int(x)+4+prefixWidth, int(y)+2, highlightWidth, int(h)-4, highlightColor)
	}

	// Draw text
	textRectX := int(x) + 4
	textRectY := int(y) + 2
	textRectW := int(w) - 8
	textRectH := int(h) - 4
	linux.XDrawTextRect(display, drawable, textRectX, textRectY, textRectW, textRectH, fontName, fontSize, text, textColor, linux.ALIGN_LEFT|linux.ALIGN_VCENTER|linux.ALIGN_SINGLELINE)

	// Draw caret if focused and no selection
	if linux.CT.Visible && ti.Focused() && selStart == selEnd {
		caretPos := int(ti.Caret())
		runes := []rune(text)
		if caretPos < 0 {
			caretPos = 0
		}
		if caretPos > len(runes) {
			caretPos = len(runes)
		}
		caretX := int(x) + 4
		if caretPos > 0 {
			caretWidth := linux.XTextWidth(display, drawable, fontName, fontSize, string(runes[:caretPos]))
			caretX += caretWidth
		}
		caretHeight := int(float32(h) * 0.5)
		caretY := int(y) + (int(h)-caretHeight)/2
		caretColor := &common.Color{Red: 0, Green: 0, Blue: 0}
		linux.XFillRect(display, drawable, caretX, caretY, 2, caretHeight, caretColor)
	}
}

// registerTextInput registers the text input component with the windows package so the state can be tracked and updated while rendering or responding to events.
//
// Parameters:
//   - ti: The TextInput component to be registered.
func registerTextInput(ti TextInput) {
	linux.RegisterTextInputState(ti.ID(), ti.(*textInput).state)
}
