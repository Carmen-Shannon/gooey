//go:build windows
// +build windows

package window

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Carmen-Shannon/gooey/common"
	"github.com/Carmen-Shannon/gooey/component"
	wdws "github.com/Carmen-Shannon/gooey/internal/windows"

	"golang.org/x/sys/windows"
)

// createWindow creates a new window with the specified options.
// It takes a variadic number of NewWindowOption functions to customize the window's properties.
// The function sets default values for the title, style, and class name if they are not provided.
//
// Parameters:
//   - options: A variadic list of NewWindowOption functions that modify the window's properties.
//
// Returns:
//   - Window: A new window instance with the specified properties.
func createWindow(options ...NewWindowOption) Window {
	opts := newWindowOption{}
	for _, opt := range options {
		opt(&opts)
	}

	if opts.Title == "" {
		opts.Title = "Gooey"
	}
	if opts.Style == 0 {
		opts.Style = wdws.WS_OVERLAPPEDWINDOW
	}
	if opts.ClassName == "" {
		opts.ClassName = "GooeyWindow"
	}

	style := ClassStyleFlagByteAlignClient | ClassStyleFlagByteAlignWindow | ClassStyleFlagGlobalClass | ClassStyleFlagDoubleClicks | ClassStyleFlagWidthRedraw | ClassStyleFlagHeightRedraw | ClassStyleFlagClipChildren
	wndProc := windows.NewCallback(wdws.WindowProc)
	clsName, _ := windows.UTF16PtrFromString(opts.ClassName)
	wdwTitle, _ := windows.UTF16PtrFromString(opts.Title)

	brush := wdws.CreateSolidBrush(opts.BackgroundColor)

	_, err := wdws.RegisterClassExW(
		wdws.StyleOpt(uint32(style)),
		wdws.ProcedureOpt(wndProc),
		wdws.InstanceHandleOpt(uintptr(windows.CurrentProcess())),
		wdws.ClassNameOpt(opts.ClassName),
		wdws.BackgroundHandleOpt(uintptr(brush)),
	)
	if err != nil {
		return nil
	}

	wdwHandle, err := wdws.CreateWindow(
		wdws.CreateWindowOptClassName(clsName),
		wdws.CreateWindowOptWindowName(wdwTitle),
		wdws.CreateWindowOptStyle(opts.Style),
		wdws.CreateWindowOptInstance(windows.CurrentProcess()),
		wdws.CreateWindowOptSize(opts.Width, opts.Height),
	)
	if err != nil {
		return nil
	}

	w := &wdw{
		mu:     sync.Mutex{},
		ID:     uintptr(wdwHandle),
		Height: opts.Height,
		Width:  opts.Width,
		Title:  opts.Title,
	}

	wdws.RegisterDrawCallback(uintptr(wdwHandle), func(hdc uintptr) {
		w.DrawComponents(&common.DrawCtx{
			Hwnd: uintptr(wdwHandle),
			Hdc:  hdc,
		})
	})
	wdws.SetWindowColor(uintptr(wdwHandle), opts.BackgroundColor)

	return w
}

// setWindowDisplay sets the display state of the window.
// It takes a pointer to the window and a WindowDisplayFlag as parameters.
// The function uses the ShowWindow function from the Windows API to change the window's visibility.
//
// Parameters:
//   - w: A pointer to the window whose display state is to be set.
//   - flag: A WindowDisplayFlag that indicates the desired display state.
//
// Returns:
//   - error: An error if the operation fails, or nil if it succeeds.
func setWindowDisplay(w *wdw, flag WindowDisplayFlag) error {
	var cmd wdws.ShowWindowCmd
	switch flag {
	case WindowDisplayFlagShow:
		cmd = wdws.SW_SHOWNORMAL
	case WindowDisplayFlagHide:
		cmd = wdws.SW_HIDE
	case WindowDisplayFlagMaximize:
		cmd = wdws.SW_SHOWMAXIMIZED
	case WindowDisplayFlagMinimize:
		cmd = wdws.SW_SHOWMINIMIZED
	default:
		cmd = wdws.SW_SHOWNORMAL
	}
	ok := wdws.ShowWindow(windows.Handle(w.ID), cmd)
	if !ok {
		return windows.ERROR_INVALID_WINDOW_HANDLE
	}
	return nil
}

// getMessageHandler retrieves messages from the message queue for the specified window handle.
// It uses the GetMessage function from the Windows API to retrieve messages.
// The function runs in a loop, processing messages until an error occurs or the window handle is invalid.
// This function MUST remain on the main OS thread otherwise the application will become unresponsive.
//
// Parameters:
//   - msg: A pointer to a wdws.Msg structure that will receive the message.
//   - windowHandle: The handle of the window for which to retrieve messages.
func getMessageHandler(msg *wdws.Msg, windowHandle windows.Handle) {
	for {
		_, err := wdws.GetMessage(wdws.MessageOpt(msg), wdws.WindowHandleOpt(windows.Handle(windowHandle)))
		if err != nil {
			if !strings.Contains(err.Error(), "Invalid window handle.") {
				fmt.Println("error getting message:", err)
			}
			break
		}
		wdws.TranslateMessage(msg)
		wdws.DispatchMessage(msg)
	}
}

// drawComponents draws the components within the window.
// It won't update the rendering of components while re-sizing.
//
// Parameters:
//   - w: A pointer to the window whose components are to be drawn.
//   - ctx: A pointer to a common.DrawCtx structure that contains the drawing context.
func drawComponents(w *wdw, ctx *common.DrawCtx) {
	if wdws.IsResizing(w.ID) {
		for _, c := range w.Components {
			c.Draw(ctx)
		}
		return
	}
	x, y, _, _, _ := wdws.GetMouseState(windows.Handle(w.ID))
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

			wdws.SetCustomCursorDraw(inside)
			if ti.Enabled() && inside {
				wdws.SetCursor(wdws.LoadIBeamCursor())
			}
		}
		c.Draw(ctx)
	}
}

// run starts the window's message loop and begins processing events.
// It locks the OS thread to ensure that the window runs on the main thread.
// It is responsible for handling window messages and dispatching them to the appropriate components.
// It will block until the window is closed or an error occurs.
//
// Parameters:
//   - w: A pointer to the window to run.
//   - refresh: The refresh rate in frames per second (FPS) for the window's drawing.
func run(w *wdw, refresh int) {
	_ = setWindowDisplay(w, WindowDisplayFlagShow)
	startDrawHandler(windows.Handle(w.ID), refresh)

	msg := new(wdws.Msg)
	getMessageHandler(msg, windows.Handle(w.ID))
}

// startDrawHandler starts a goroutine that periodically invalidates the window's client area.
// This triggers a redraw of the window at the specified frames per second (FPS).
// It uses a ticker to create a loop that runs at the specified interval.
// The function takes the window handle and the desired FPS as parameters.
//
// Parameters:
//   - hwnd: The handle of the window to be redrawn.
//   - fps: The desired frames per second (FPS) for the redraw interval.
func startDrawHandler(hwnd windows.Handle, fps int) {
	interval := time.Second / time.Duration(fps)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			_ = wdws.InvalidateRect(hwnd, nil, false)
		}
	}()
}
