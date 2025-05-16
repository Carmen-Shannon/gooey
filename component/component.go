package component

import (
	"fmt"

	"github.com/Carmen-Shannon/gooey/common"
)

type baseComponent struct {
	id   uintptr
	size struct {
		Width  int32
		Height int32
	}
	position struct {
		X int32
		Y int32
	}
	visible bool
	enabled bool
}

type TextAlignment int

const (
	LeftAlign TextAlignment = iota
	CenterAlign
	RightAlign
)

// NewComponent creates a new component with the specified options.
// It accepts a variadic list of CreateComponentOption functions to customize the component's properties.
//
// Parameters:
//   - options: A variadic list of CreateComponentOption functions to customize the component's properties.
//
// Returns:
//   - Component: A pointer to the newly created component.
func NewComponent(options ...CreateComponentOption) Component {
	opts := &createComponentOptions{}
	for _, opt := range options {
		opt(opts)
	}

	c := &baseComponent{
		id:      opts.ID,
		visible: opts.Visible,
		enabled: opts.Enabled,
		size: struct {
			Width  int32
			Height int32
		}{
			Width:  opts.Size.Width,
			Height: opts.Size.Height,
		},
		position: struct {
			X int32
			Y int32
		}{
			X: opts.Position.X,
			Y: opts.Position.Y,
		},
	}
	return c
}

type Component interface {
	// ID returns the unique identifier for the component.
	//
	// Returns:
	//  - uintptr: The unique identifier for the component.
	ID() uintptr

	// SetID sets the unique identifier for the component.
	//
	// Parameters:
	//  - id: The unique identifier to set for the component.
	SetID(id uintptr)

	// Size returns the width and height of the component.
	//
	// Returns:
	//  - int32: The width of the component.
	//  - int32: The height of the component.
	Size() (int32, int32)

	// SetSize sets the width and height of the component.
	//
	// Parameters:
	//  - width: The width to set for the component.
	//  - height: The height to set for the component.
	SetSize(width, height int32)

	// Position returns the x and y coordinates of the component.
	//
	// Returns:
	//  - int32: The x coordinate of the component.
	//  - int32: The y coordinate of the component.
	Position() (int32, int32)

	// SetPosition sets the x and y coordinates of the component.
	//
	// Parameters:
	//  - x: The x coordinate to set for the component.
	//  - y: The y coordinate to set for the component.
	SetPosition(x, y int32)

	// Visible returns whether the component is visible or not.
	//
	// Returns:
	//  - bool: True if the component is visible, false otherwise.
	Visible() bool

	// SetVisible sets the visibility of the component.
	//
	// Parameters:
	//  - visible: True to make the component visible, false to hide it.
	SetVisible(visible bool)

	// Enabled returns whether the component is enabled or not.
	//
	// Returns:
	//  - bool: True if the component is enabled, false otherwise.
	Enabled() bool

	// SetEnabled sets the enabled state of the component.
	//
	// Parameters:
	//  - enabled: True to enable the component, false to disable it.
	SetEnabled(enabled bool)

	// Draw draws the component using the provided context.
	//
	// Parameters:
	//  - ctx: The context to use for drawing the component.
	Draw(ctx *common.DrawCtx)
}

var _ Component = (*baseComponent)(nil)

func (c *baseComponent) ID() uintptr {
	return c.id
}

func (c *baseComponent) SetID(id uintptr) {
	c.id = id
}

func (c *baseComponent) Size() (int32, int32) {
	return c.size.Width, c.size.Height
}

func (c *baseComponent) SetSize(width, height int32) {
	c.size.Width = width
	c.size.Height = height
}

func (c *baseComponent) Position() (int32, int32) {
	return c.position.X, c.position.Y
}

func (c *baseComponent) SetPosition(x, y int32) {
	c.position.X = x
	c.position.Y = y
}

func (c *baseComponent) Visible() bool {
	return c.visible
}

func (c *baseComponent) SetVisible(visible bool) {
	c.visible = visible
}

func (c *baseComponent) Enabled() bool {
	return c.enabled
}

func (c *baseComponent) SetEnabled(enabled bool) {
	c.enabled = enabled
}

// Draw should never be called from the base Component level and will panic if done so.
func (c *baseComponent) Draw(ctx *common.DrawCtx) {
	panic("Draw must be implemented by a concrete component type")
}

// drawComponent is a helper function to draw a component based on its type.
// All component types should implement the Draw method that calls this drawComponent function.
//
// Parameters:
//   - c: The component to draw.
//   - ctx: The context to use for drawing the component.
func drawComponent(c any, ctx *common.DrawCtx) {
	switch comp := any(c).(type) {
	case Button:
		if comp.Visible() {
			drawButton(ctx, comp)
		}
	case Label:
		if comp.Visible() {
			drawLabel(ctx, comp)
		}
	case TextInput:
		if comp.Visible() {
			drawTextInput(ctx, comp)
		}
	case Selector:
		drawSelector(ctx, comp)
	default:
		fmt.Println("unsupported component type")
	}
}
