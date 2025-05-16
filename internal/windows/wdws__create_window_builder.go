//go:build windows
// +build windows

package wdws

import (
	"fmt"

	"golang.org/x/sys/windows"
)

type createWindowBuilderOpts struct {
	ExStyle    uint32
	ClassName  *uint16
	WindowName *uint16
	Style      uint32
	X          int32
	Y          int32
	Width      int32
	Height     int32
	Parent     windows.Handle
	Menu       windows.Handle
	Instance   windows.Handle
	Param      uintptr
}

func newCwBuilderOpts() *createWindowBuilderOpts {
	return &createWindowBuilderOpts{
		ExStyle:    0,
		X:          CW_USEDEFAULT,
		Y:          CW_USEDEFAULT,
		Width:      CW_USEDEFAULT,
		Height:     CW_USEDEFAULT,
		Parent:     0,
		Menu:       0,
		WindowName: nil,
		Param:      0,
	}
}

func (c *createWindowBuilderOpts) validate() error {
	if c.ClassName == nil {
		return fmt.Errorf("validation error: ClassName cannot be nil. Windows err: %w", windows.ERROR_INVALID_PARAMETER)
	}
	if c.Style == 0 {
		return fmt.Errorf("validation error: Style cannot be 0. Windows err: %w", windows.ERROR_INVALID_PARAMETER)
	}
	if c.Instance == 0 {
		return fmt.Errorf("validation error: Instance cannot be 0. Windows err: %w", windows.ERROR_INVALID_PARAMETER)
	}
	return nil
}

type CreateWindowBuilderOpt func(*createWindowBuilderOpts)

// CreateWindowOptExStyle sets the extended window style for the window being created.
//
// Parameters:
//   - exStyle: The extended style flags for the window.
//
// Returns:
//   - CreateWindowBuilderOpt: A function that takes a pointer to createWindowBuilderOpts and sets its ExStyle field.
func CreateWindowOptExStyle(exStyle uint32) CreateWindowBuilderOpt {
	return func(opts *createWindowBuilderOpts) {
		opts.ExStyle = exStyle
	}
}

// CreateWindowOptClassName sets the class name for the window being created.
//
// Parameters:
//   - className: A pointer to a UTF-16 string representing the window class name.
//
// Returns:
//   - CreateWindowBuilderOpt: A function that takes a pointer to createWindowBuilderOpts and sets its ClassName field.
func CreateWindowOptClassName(className *uint16) CreateWindowBuilderOpt {
	return func(opts *createWindowBuilderOpts) {
		opts.ClassName = className
	}
}

// CreateWindowOptWindowName sets the window name (title) for the window being created.
//
// Parameters:
//   - windowName: A pointer to a UTF-16 string representing the window name.
//
// Returns:
//   - CreateWindowBuilderOpt: A function that takes a pointer to createWindowBuilderOpts and sets its WindowName field.
func CreateWindowOptWindowName(windowName *uint16) CreateWindowBuilderOpt {
	return func(opts *createWindowBuilderOpts) {
		opts.WindowName = windowName
	}
}

// CreateWindowOptStyle sets the style flags for the window being created.
//
// Parameters:
//   - style: The style flags for the window.
//
// Returns:
//   - CreateWindowBuilderOpt: A function that takes a pointer to createWindowBuilderOpts and sets its Style field.
func CreateWindowOptStyle(style uint32) CreateWindowBuilderOpt {
	return func(opts *createWindowBuilderOpts) {
		opts.Style = style
	}
}

// CreateWindowOptPosition sets the position of the window being created.
//
// Parameters:
//   - x: The x-coordinate of the window.
//   - y: The y-coordinate of the window.
//
// Returns:
//   - CreateWindowBuilderOpt: A function that takes a pointer to createWindowBuilderOpts and sets its X and Y fields.
func CreateWindowOptPosition(x, y int32) CreateWindowBuilderOpt {
	return func(opts *createWindowBuilderOpts) {
		opts.X = x
		opts.Y = y
	}
}

// CreateWindowOptSize sets the size of the window being created.
//
// Parameters:
//   - width: The width of the window.
//   - height: The height of the window.
//
// Returns:
//   - CreateWindowBuilderOpt: A function that takes a pointer to createWindowBuilderOpts and sets its Width and Height fields.
func CreateWindowOptSize(width, height int32) CreateWindowBuilderOpt {
	return func(opts *createWindowBuilderOpts) {
		opts.Width = width
		opts.Height = height
	}
}

// CreateWindowOptParent sets the parent window handle for the window being created.
//
// Parameters:
//   - parent: The handle to the parent window.
//
// Returns:
//   - CreateWindowBuilderOpt: A function that takes a pointer to createWindowBuilderOpts and sets its Parent field.
func CreateWindowOptParent(parent windows.Handle) CreateWindowBuilderOpt {
	return func(opts *createWindowBuilderOpts) {
		opts.Parent = parent
	}
}

// CreateWindowOptMenu sets the menu handle for the window being created.
//
// Parameters:
//   - menu: The handle to the menu to be used by the window.
//
// Returns:
//   - CreateWindowBuilderOpt: A function that takes a pointer to createWindowBuilderOpts and sets its Menu field.
func CreateWindowOptMenu(menu windows.Handle) CreateWindowBuilderOpt {
	return func(opts *createWindowBuilderOpts) {
		opts.Menu = menu
	}
}

// CreateWindowOptInstance sets the instance handle for the window being created.
//
// Parameters:
//   - instance: The handle to the application instance.
//
// Returns:
//   - CreateWindowBuilderOpt: A function that takes a pointer to createWindowBuilderOpts and sets its Instance field.
func CreateWindowOptInstance(instance windows.Handle) CreateWindowBuilderOpt {
	return func(opts *createWindowBuilderOpts) {
		opts.Instance = instance
	}
}

// CreateWindowOptParam sets the additional parameter for the window being created.
//
// Parameters:
//   - param: A pointer to additional data to be passed to the window.
//
// Returns:
//   - CreateWindowBuilderOpt: A function that takes a pointer to createWindowBuilderOpts and sets its Param field.
func CreateWindowOptParam(param uintptr) CreateWindowBuilderOpt {
	return func(opts *createWindowBuilderOpts) {
		opts.Param = param
	}
}
