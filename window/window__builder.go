package window

import "github.com/Carmen-Shannon/gooey/common"

type newWindowOption struct {
	Title           string
	Height          int32
	Width           int32
	Style           uint32
	ClassName       string
	CloseChan       chan struct{}
	BackgroundColor *common.Color
}

type NewWindowOption func(*newWindowOption)

// TitleOpt sets the title of the window.
//
// Parameters:
//  - title: The title to set for the window.
//
// Returns:
//  - NewWindowOption: A function that takes a pointer to newWindowOption and sets the Title field.
func TitleOpt(title string) NewWindowOption {
	return func(opts *newWindowOption) {
		opts.Title = title
	}
}

// HeightOpt sets the height of the window.
//
// Parameters:
//  - height: The height to set for the window.
//
// Returns:
//  - NewWindowOption: A function that takes a pointer to newWindowOption and sets the Height field.
func HeightOpt(height int32) NewWindowOption {
	return func(opts *newWindowOption) {
		opts.Height = height
	}
}

// WidthOpt sets the width of the window.
//
// Parameters:
//  - width: The width to set for the window.
//
// Returns:
//  - NewWindowOption: A function that takes a pointer to newWindowOption and sets the Width field.
func WidthOpt(width int32) NewWindowOption {
	return func(opts *newWindowOption) {
		opts.Width = width
	}
}

// StyleOpt sets the style of the window.
//
// Parameters:
//  - style: The style to set for the window.
//
// Returns:
//  - NewWindowOption: A function that takes a pointer to newWindowOption and sets the Style field.
func StyleOpt(style uint32) NewWindowOption {
	return func(opts *newWindowOption) {
		opts.Style = style
	}
}

// ClassNameOpt sets the class name of the window.
//
// Parameters:
//  - className: The class name to set
//
// Returns:
//  - NewWindowOption: A function that takes a pointer to newWindowOption and sets the ClassName field.
func ClassNameOpt(className string) NewWindowOption {
	return func(opts *newWindowOption) {
		opts.ClassName = className
	}
}

// CloseChanOpt sets the close channel of the window.
//
// Parameters:
//  - closeChan: The close channel to set for the window.
//
// Returns:
//  - NewWindowOption: A function that takes a pointer to newWindowOption and sets the CloseChan field.
func CloseChanOpt(closeChan chan struct{}) NewWindowOption {
	return func(opts *newWindowOption) {
		opts.CloseChan = closeChan
	}
}

// BackgroundColorOpt sets the background color of the window.
//
// Parameters:
//  - color: The background color to set for the window.
//
// Returns:
//  - NewWindowOption: A function that takes a pointer to newWindowOption and sets the BackgroundColor field.
func BackgroundColorOpt(color *common.Color) NewWindowOption {
	return func(opts *newWindowOption) {
		opts.BackgroundColor = color
	}
}
