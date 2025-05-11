//go:build windows
// +build windows

package wdws

import "unsafe"

type registerClassOpts struct {
	Size             uint32
	Style            uint32
	Procedure        uintptr
	ExtraBytesStruct int32
	ExtraBytesWindow int32
	InstanceHandle   uintptr
	IconHandle       uintptr
	CursorHandle     uintptr
	BackgroundHandle uintptr
	MenuName         string
	ClassName        string
	SmallIconHandle  uintptr
}

func newRegisterClassOpts() *registerClassOpts {
	return &registerClassOpts{
		Size:             uint32(unsafe.Sizeof(wdsWndClass{})),
		Style:            0,
		Procedure:        0,
		ExtraBytesStruct: 0,
		ExtraBytesWindow: 0,
		InstanceHandle:   0,
		IconHandle:       0,
		CursorHandle:     0,
		BackgroundHandle: 0,
		MenuName:         "",
		ClassName:        "",
	}
}

type RegisterClassOpt func(*registerClassOpts)

// SizeOpt sets the size of the reigster class options.
//
// Parameters:
//   - size: The size of the register class options.
//
// Returns:
//   - RegisterClassOpt: A function that takes a pointer to registerClassOpts and sets its size field.
func SizeOpt(size uint32) RegisterClassOpt {
	return func(opts *registerClassOpts) {
		opts.Size = size
	}
}

// StyleOpt sets the style of the register class options.
//
// Parameters:
//   - style: The style flags for the register class.
//
// Returns:
//   - RegisterClassOpt: A function that takes a pointer to registerClassOpts and sets its Style field.
func StyleOpt(style uint32) RegisterClassOpt {
	return func(opts *registerClassOpts) {
		opts.Style = style
	}
}

// ProcedureOpt sets the window procedure for the register class options.
//
// Parameters:
//   - procedure: The pointer to the window procedure function.
//
// Returns:
//   - RegisterClassOpt: A function that takes a pointer to registerClassOpts and sets its Procedure field.
func ProcedureOpt(procedure uintptr) RegisterClassOpt {
	return func(opts *registerClassOpts) {
		opts.Procedure = procedure
	}
}

// ExtraBytesStructOpt sets the extra bytes for the structure in the register class options.
//
// Parameters:
//   - extraBytesStruct: The number of extra bytes to allocate following the window-class structure.
//
// Returns:
//   - RegisterClassOpt: A function that takes a pointer to registerClassOpts and sets its ExtraBytesStruct field.
func ExtraBytesStructOpt(extraBytesStruct int32) RegisterClassOpt {
	return func(opts *registerClassOpts) {
		opts.ExtraBytesStruct = extraBytesStruct
	}
}

// ExtraBytesWindowOpt sets the extra bytes for the window in the register class options.
//
// Parameters:
//   - extraBytesWindow: The number of extra bytes to allocate following the window instance.
//
// Returns:
//   - RegisterClassOpt: A function that takes a pointer to registerClassOpts and sets its ExtraBytesWindow field.
func ExtraBytesWindowOpt(extraBytesWindow int32) RegisterClassOpt {
	return func(opts *registerClassOpts) {
		opts.ExtraBytesWindow = extraBytesWindow
	}
}

// InstanceHandleOpt sets the instance handle for the register class options.
//
// Parameters:
//   - instanceHandle: The handle to the instance that contains the window procedure.
//
// Returns:
//   - RegisterClassOpt: A function that takes a pointer to registerClassOpts and sets its InstanceHandle field.
func InstanceHandleOpt(instanceHandle uintptr) RegisterClassOpt {
	return func(opts *registerClassOpts) {
		opts.InstanceHandle = instanceHandle
	}
}

// IconHandleOpt sets the icon handle for the register class options.
//
// Parameters:
//   - iconHandle: The handle to the icon to be used by windows of this class.
//
// Returns:
//   - RegisterClassOpt: A function that takes a pointer to registerClassOpts and sets its IconHandle field.
func IconHandleOpt(iconHandle uintptr) RegisterClassOpt {
	return func(opts *registerClassOpts) {
		opts.IconHandle = iconHandle
	}
}

// CursorHandleOpt sets the cursor handle for the register class options.
//
// Parameters:
//   - cursorHandle: The handle to the cursor to be used by windows of this class.
//
// Returns:
//   - RegisterClassOpt: A function that takes a pointer to registerClassOpts and sets its CursorHandle field.
func CursorHandleOpt(cursorHandle uintptr) RegisterClassOpt {
	return func(opts *registerClassOpts) {
		opts.CursorHandle = cursorHandle
	}
}

// BackgroundHandleOpt sets the background brush handle for the register class options.
//
// Parameters:
//   - backgroundHandle: The handle to the background brush.
//
// Returns:
//   - RegisterClassOpt: A function that takes a pointer to registerClassOpts and sets its BackgroundHandle field.
func BackgroundHandleOpt(backgroundHandle uintptr) RegisterClassOpt {
	return func(opts *registerClassOpts) {
		opts.BackgroundHandle = backgroundHandle
	}
}

// MenuNameOpt sets the menu name for the register class options.
//
// Parameters:
//   - menuName: The name of the menu resource to be used by windows of this class.
//
// Returns:
//   - RegisterClassOpt: A function that takes a pointer to registerClassOpts and sets its MenuName field.
func MenuNameOpt(menuName string) RegisterClassOpt {
	return func(opts *registerClassOpts) {
		opts.MenuName = menuName
	}
}

// ClassNameOpt sets the class name for the register class options.
//
// Parameters:
//   - className: The name of the window class.
//
// Returns:
//   - RegisterClassOpt: A function that takes a pointer to registerClassOpts and sets its ClassName field.
func ClassNameOpt(className string) RegisterClassOpt {
	return func(opts *registerClassOpts) {
		opts.ClassName = className
	}
}

// SmallIconHandleOpt sets the small icon handle for the register class options.
//
// Parameters:
//   - smallIconHandle: The handle to the small icon to be used by windows of this class.
//
// Returns:
//   - RegisterClassOpt: A function that takes a pointer to registerClassOpts and sets its SmallIconHandle field.
func SmallIconHandleOpt(smallIconHandle uintptr) RegisterClassOpt {
	return func(opts *registerClassOpts) {
		opts.SmallIconHandle = smallIconHandle
	}
}
