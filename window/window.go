package window

import (
	"runtime"
	"slices"
	"sync"

	"github.com/Carmen-Shannon/gooey/common"
	"github.com/Carmen-Shannon/gooey/component"
)

type wdw struct {
	mu sync.Mutex

	ID              uintptr
	Height          int32
	Width           int32
	Title           string
	BackgroundColor common.Color
	Components      []component.Component
}

type Window interface {
	// AddComponent adds a component to the window's list of components.
	// It takes a component.Component as a parameter.
	//
	// Parameters:
	//  - c: The component to add to the window.
	AddComponent(c component.Component)

	// DrawComponents draws the components of the window using the provided context.
	// It iterates over the window's components and calls their Draw method.
	//
	// Parameters:
	//  - ctx: The context to use for drawing the components.
	DrawComponents(ctx *common.DrawCtx)

	// GetComponent retrieves a component from the window's list of components by its ID.
	// It takes a uintptr as a parameter and returns the corresponding component.Component.
	//
	// Parameters:
	//  - id: The ID of the component to retrieve.
	//
	// Returns:
	//  - component.Component: The component with the specified ID, or nil if not found.
	GetComponent(id uintptr) component.Component

	// GetID returns the ID of the window.
	// It is a unique identifier for the window instance.
	// It is used to identify the window in various operations.
	//
	// Returns:
	//  - uintptr: The ID of the window.
	GetID() uintptr

	// RemoveComponent removes a component from the window's list of components.
	// It takes a component.Component as a parameter.
	// The component is identified by its ID, and if found, it is removed from the list.
	//
	// Parameters:
	//  - c: The component to remove from the window.
	RemoveComponent(id uintptr)

	// Run starts the window's message loop and begins processing events.
	// It locks the OS thread to ensure that the window runs on the main thread.
	//
	// Note: This function will lock the OS thread so it should be called from the main goroutine.
	// It is responsible for handling window messages and dispatching them to the appropriate components.
	// It will block until the window is closed or an error occurs.
	//
	// Parameters:
	//  - refresh: An integer value that specifies the refresh rate or interval for the window in FPS.
	Run(refresh int)

	// SetWindowDisplay sets the display state of the window.
	// It takes a WindowDisplayFlag to specify the desired display state.
	//
	// Parameters:
	//  - flag: A WindowDisplayFlag that indicates the desired display state.
	//
	// Returns:
	//  - error: An error if the operation fails, or nil if it succeeds.
	SetWindowDisplay(flag WindowDisplayFlag) error
}

var _ Window = (*wdw)(nil)

// NewWindow creates a new window with the specified options.
// It takes a variadic number of NewWindowOption functions to customize the window's properties.
//
// Parameters:
// - options: A variadic list of NewWindowOption functions that modify the window's properties.
//
// Returns:
// - Window: A new window instance with the specified properties.
func NewWindow(options ...NewWindowOption) Window {
	return createWindow(options...)
}

func (w *wdw) AddComponent(c component.Component) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.Components = append(w.Components, c)
}

func (w *wdw) DrawComponents(ctx *common.DrawCtx) {
	drawComponents(w, ctx)
}

func (w *wdw) GetComponent(id uintptr) component.Component {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, comp := range w.Components {
		if comp.ID() == id {
			return comp
		}
	}
	return nil
}

func (w *wdw) GetID() uintptr {
	return w.ID
}

func (w *wdw) RemoveComponent(id uintptr) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for i, comp := range w.Components {
		if comp.ID() == id {
			w.Components = slices.Delete(w.Components, i, i+1)
			return
		}
	}
}

func (w *wdw) Run(refresh int) {
	runtime.LockOSThread()
	run(w, refresh)
}

func (w *wdw) SetWindowDisplay(flag WindowDisplayFlag) error {
	return setWindowDisplay(w, flag)
}
