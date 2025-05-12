//go:build linux
// +build linux

package window

import (
	"errors"
	"sync"
	"time"

	"github.com/Carmen-Shannon/gooey/common"
	"github.com/Carmen-Shannon/gooey/component"
	"github.com/Carmen-Shannon/gooey/internal/linux"
)

func createWindow(options ...NewWindowOption) Window {
	opts := newWindowOption{}
	for _, opt := range options {
		opt(&opts)
	}

	if opts.Title == "" {
		opts.Title = "Gooey"
	}
	if opts.Width == 0 {
		opts.Width = 800
	}
	if opts.Height == 0 {
		opts.Height = 600
	}
	bgColor := common.Color{Red: 255, Green: 255, Blue: 255}
	if opts.BackgroundColor != nil {
		bgColor = *opts.BackgroundColor
	}
	bgPixel := uint32(bgColor.Red)<<16 | uint32(bgColor.Green)<<8 | uint32(bgColor.Blue)

	display := linux.XOpenDisplay()
	if display == nil {
		panic("cannot open X display")
	}
	screen := linux.XDefaultScreen(display)
	root := linux.XRootWindow(display, screen)

	window := linux.XCreateSimpleWindow(
		display,
		root,
		0, 0,
		uint(opts.Width), uint(opts.Height), 1,
		0, bgPixel,
	)

	linux.XSelectInput(display, window,
		linux.ExposureMask|
			linux.StructureNotifyMask|
			linux.ButtonPressMask|
			linux.ButtonReleaseMask|
			linux.PointerMotionMask|
			linux.KeyPressMask|
			linux.KeyReleaseMask,
	)
	linux.RegisterDisplay(uintptr(window), display)
	linux.XStoreName(display, window, opts.Title)
	linux.XMapWindow(display, window)

	w := &wdw{
		mu:              sync.Mutex{},
		ID:              uintptr(window),
		Height:          opts.Height,
		Width:           opts.Width,
		Title:           opts.Title,
		BackgroundColor: bgColor,
		redraw:          make(chan struct{}, 1),
	}

	linux.RegisterDrawCallback(w.ID, func(hdc uintptr) {
		w.DrawComponents(&common.DrawCtx{
			Hwnd: w.ID,
			Hdc:  hdc,
		})
	})
	linux.SetWindowColor(w.ID, opts.BackgroundColor)

	return w
}

func drawComponents(w *wdw, ctx *common.DrawCtx) {
	if linux.IsResizing(w.ID) {
		for _, c := range w.Components {
			c.Draw(ctx)
		}
		return
	}

	// Get mouse position (implement GetMouseState for Linux if needed)
	x, y := linux.GetMouseState(w.ID)
	display := linux.GetDisplay(w.ID)
	window := linux.C_Window(w.ID)

	for _, c := range w.Components {
		if btn, ok := c.(component.Button); ok {
			bx, by := btn.Position()
			bw, bh := btn.Size()
			inside := x >= bx && x < bx+bw && y >= by && y < by+bh
			btn.SetHovered(inside)
		}
		if ti, ok := c.(component.TextInput); ok {
			tx, ty := ti.Position()
			tw, th := ti.Size()
			inside := x >= tx && x < tx+tw && y >= ty && y < ty+th

			linux.SetCustomCursorDraw(inside)
			if ti.Enabled() && inside {
				linux.SetCursor(display, window, linux.LoadIBeamCursor(display))
			}
		}
		c.Draw(ctx)
	}

	if !linux.IsCustomCursorDraw() {
		linux.SetCursor(display, window, linux.LoadArrowCursor(display))
	}
}

func run(w *wdw, refresh int) {
	setWindowDisplay(w, WindowDisplayFlagShow)
	startDrawHandler(w, refresh)

	display := linux.GetDisplay(w.ID)
	if display == nil {
		panic("cannot get X display for window")
	}

	var event linux.C_XEvent

	for {
		select {
		case <-w.redraw:
			// Call your drawing logic directly
			linux.HandlePaint(w.ID, display)
			// Optionally flush to X server
			linux.XFlush(display)
		default:
			linux.XNextEvent(display, &event)
			cont := linux.WindowProc(w.ID, display, &event)
			if !cont {
				return
			}
		}
	}
}

func setWindowDisplay(w *wdw, flag WindowDisplayFlag) error {
	display := linux.GetDisplay(w.ID)
	if display == nil {
		return errors.New("cannot get X display for window")
	}
	window := linux.C_Window(w.ID)

	switch flag {
	case WindowDisplayFlagShow:
		linux.XShowWindow(display, window)
	case WindowDisplayFlagHide:
		linux.XHideWindow(display, window)
	case WindowDisplayFlagMaximize:
		linux.XMaximizeWindow(display, window)
	case WindowDisplayFlagMinimize:
		linux.XMinimizeWindow(display, window)
	default:
		return errors.New("unsupported WindowDisplayFlag")
	}
	return nil
}

func startDrawHandler(w *wdw, fps int) {
	interval := time.Second / time.Duration(fps)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			select {
			case w.redraw <- struct{}{}:
			default:
			}
			// display := linux.GetDisplay(w.ID)
			// if display == nil {
			// 	continue
			// }
			// window := linux.C_Window(w.ID)
			// // Trigger an Expose event by clearing a 1x1 area (does not actually clear, just triggers event)
			// linux.XClearArea(display, window, 0, 0, 0, 0, true)
		}
	}()
}
