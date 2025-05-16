package component

import "github.com/Carmen-Shannon/gooey/common"

type selector struct {
	baseComponent
	color   *common.Color
	opacity float32
	drawing bool
	state   *common.SelectorState
}

type Selector interface {
	Component

	// Color returns the color of the selector.
	//
	// Returns:
	//  - *common.Color: The color of the selector.
	Color() *common.Color

	// SetColor sets the color of the selector.
	//
	// Parameters:
	//  - color: *common.Color to set for the selector.
	SetColor(color *common.Color)

	// Opacity returns the opacity of the selector.
	//
	// Returns:
	//  - float32: The opacity of the selector.
	Opacity() float32

	// SetOpacity sets the opacity of the selector.
	//
	// Parameters:
	//  - opacity: float32 to set for the selector.
	SetOpacity(opacity float32)

	// Drawing returns whether the selector is drawing.
	//
	// Returns:
	//  - bool: True if the selector is drawing, false otherwise.
	Drawing() bool

	// SetDrawing sets the drawing state of the selector.
	//
	// Parameters:
	//  - drawing: bool to set for the selector.
	SetDrawing(drawing bool)

	// StartCapture starts the capture state of the selector, allowing it to be sized by the mouse.
	StartCapture()
}

var _ Selector = (*selector)(nil)

// NewSelector creates a new selector component with the specified options.
// It initializes the selector with default values and applies the provided options.
//
// Parameters:
//   - options: A variadic list of CreateSelectorOption functions to customize the selector.
//
// Returns:
//   - Selector: A pointer to the newly created selector component.
func NewSelector(options ...CreateSelectorOption) Selector {
	opts := newCreateSelectorOptions()
	for _, opt := range options {
		opt(opts)
	}
	cOpts := newCreateComponentOptions()
	for _, opt := range opts.ComponentOptions {
		opt(cOpts)
	}

	s := &selector{
		baseComponent: baseComponent{
			id:      cOpts.ID,
			visible: cOpts.Visible,
			enabled: cOpts.Enabled,
			size: struct {
				Width  int32
				Height int32
			}{
				Width:  cOpts.Size.Width,
				Height: cOpts.Size.Height,
			},
			position: struct {
				X int32
				Y int32
			}{
				X: cOpts.Position.X,
				Y: cOpts.Position.Y,
			},
		},
		color:   opts.Color,
		opacity: opts.Opacity,
		drawing: opts.Drawing,
		state: &common.SelectorState{
			ID:      0,
			Drawing: opts.Drawing,
			Bounds: common.Rect{
				X: cOpts.Position.X,
				Y: cOpts.Position.Y,
				W: cOpts.Size.Width,
				H: cOpts.Size.Height,
			},
			Color:    opts.Color,
			Opacity:  opts.Opacity,
			Visible:  cOpts.Visible,
			Blocking: false,
			CbMap:    make(map[string]func(any)),
		},
	}

	cbMap := make(map[string]func(any))
	cbMap["drawing"] = func(drawing any) {
		if drawingBool, ok := drawing.(bool); ok {
			s.drawing = drawingBool
		}
	}
	cbMap["blocking"] = func(blocking any) {
		if blockingBool, ok := blocking.(bool); ok {
			s.state.Blocking = blockingBool
		}
	}
	cbMap["bounds"] = func(bounds any) {
		if boundsRect, ok := bounds.(common.Rect); ok {
			s.position.X = boundsRect.X
			s.position.Y = boundsRect.Y
			s.size.Width = boundsRect.W
			s.size.Height = boundsRect.H
		}
	}
	cbMap["visible"] = func(visible any) {
		if visibleBool, ok := visible.(bool); ok {
			s.visible = visibleBool
		}
	}
	s.state.CbMap = cbMap

	registerSelector(cOpts.ID, s)

	return s
}

func (s *selector) Draw(ctx *common.DrawCtx) {
	drawComponent(s, ctx)
}

func (s *selector) Color() *common.Color {
	return s.color
}

func (s *selector) SetColor(color *common.Color) {
	s.color = color
	s.state.Color = color
}

func (s *selector) Opacity() float32 {
	return s.opacity
}

func (s *selector) SetOpacity(opacity float32) {
	s.opacity = opacity
	s.state.Opacity = opacity
}

func (s *selector) Drawing() bool {
	return s.drawing
}

func (s *selector) SetDrawing(drawing bool) {
	s.drawing = drawing
	s.state.Drawing = drawing
}

func (s *selector) SetVisible(visible bool) {
	s.baseComponent.SetVisible(visible)
	s.state.Visible = visible
}

func (s *selector) StartCapture() {
	s.SetDrawing(true)
	s.SetVisible(true)
}
