package component

import "github.com/Carmen-Shannon/gooey/common"

type createSelectorOptions struct {
	Color            *common.Color
	Opacity          float32
	Bounds           common.Rect
	Drawing          bool
	ComponentOptions []CreateComponentOption
}

type CreateSelectorOption func(*createSelectorOptions)

func newCreateSelectorOptions() *createSelectorOptions {
	return &createSelectorOptions{
		Color:   &common.Color{Red: 0, Green: 0, Blue: 0},
		Opacity: 1.0,
		Bounds:  common.Rect{X: 0, Y: 0, W: 100, H: 100},
		Drawing: false,
	}
}

// SelectorColorOpt creates a CreateSelectorOption that sets the color of the selector.
// It takes a pointer to a common.Color struct as an argument.
//
// Parameters:
//   - color: A pointer to a common.Color struct representing the color of the selector.
//
// Returns:
//   - CreateSelectorOption: A function that takes a pointer to createSelectorOptions and sets its Color field.
func SelectorColorOpt(color *common.Color) CreateSelectorOption {
	return func(opts *createSelectorOptions) {
		opts.Color = color
	}
}

// SelectorOpacityOpt creates a CreateSelectorOption that sets the opacity of the selector.
// It takes a float32 value as an argument.
//
// Parameters:
//   - opacity: A float32 value representing the opacity of the selector.
//
// Returns:
//   - CreateSelectorOption: A function that takes a pointer to createSelectorOptions and sets its Opacity field.
func SelectorOpacityOpt(opacity float32) CreateSelectorOption {
	return func(opts *createSelectorOptions) {
		opts.Opacity = opacity
	}
}

// SelectorBoundsOpt creates a CreateSelectorOption that sets the bounds of the selector.
// It takes a common.Rect struct as an argument.
//
// Parameters:
//   - bounds: A common.Rect struct representing the bounds of the selector.
//
// Returns:
//   - CreateSelectorOption: A function that takes a pointer to createSelectorOptions and sets its Bounds field.
func SelectorBoundsOpt(bounds common.Rect) CreateSelectorOption {
	return func(opts *createSelectorOptions) {
		opts.Bounds = bounds
	}
}

// SelectorDrawingOpt creates a CreateSelectorOption that sets whether the selector is drawing or not.
// It takes a boolean value as an argument.
//
// Parameters:
//   - drawing: A boolean value indicating whether the selector is drawing or not.
//
// Returns:
//   - CreateSelectorOption: A function that takes a pointer to createSelectorOptions and sets its Drawing field.
func SelectorDrawingOpt(drawing bool) CreateSelectorOption {
	return func(opts *createSelectorOptions) {
		opts.Drawing = drawing
	}
}

// SelectorComponentOptionsOpt creates a CreateSelectorOption that appends component options to the selector.
// It takes a variadic list of CreateComponentOption functions as arguments.
//
// Parameters:
//   - options: A variadic list of CreateComponentOption functions to append to the selector.
//
// Returns:
//   - CreateSelectorOption: A function that takes a pointer to createSelectorOptions and appends the component options.
func SelectorComponentOptionsOpt(options ...CreateComponentOption) CreateSelectorOption {
	return func(opts *createSelectorOptions) {
		opts.ComponentOptions = append(opts.ComponentOptions, options...)
	}
}
