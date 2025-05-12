//go:build linux
// +build linux

package component

import (
	"github.com/Carmen-Shannon/gooey/common"
	"github.com/Carmen-Shannon/gooey/internal/linux"
)

func drawLabel(ctx *common.DrawCtx, l Label) {
	x, y := l.Position()
	w, h := l.Size()
	if w <= 0 || h <= 0 {
		return
	}
	fontName := l.Font()
	fontSize := l.TextSize()
	text := l.Text()
	color := l.Color()

	display := linux.GetDisplay(ctx.Hwnd)
	if display == nil {
		return
	}
	drawable := linux.C_Drawable(ctx.Hdc)

	// Only shrink font size if text is too wide and word wrap is off
	if !l.WordWrap() {
		for fontSize > 12 {
			textWidth := linux.XTextWidth(display, drawable, fontName, int(fontSize), text)
			if textWidth <= int(w) {
				break
			}
			fontSize--
		}
	}

	alignment := l.TextAlignment()
	var alignFlag int
	switch alignment {
	case LeftAlign:
		alignFlag = linux.ALIGN_LEFT
	case CenterAlign:
		alignFlag = linux.ALIGN_CENTER
	case RightAlign:
		alignFlag = linux.ALIGN_RIGHT
	}

	format := alignFlag
	if l.WordWrap() {
		format |= linux.ALIGN_WORDBREAK
		// Optionally expand height if needed (not implemented here, but you could measure and adjust)
	} else {
		format |= linux.ALIGN_VCENTER | linux.ALIGN_SINGLELINE
	}

	linux.XDrawTextRect(display, drawable, int(x), int(y), int(w), int(h), fontName, int(fontSize), text, color, format)
}
