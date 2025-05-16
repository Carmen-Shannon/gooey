package component

type createComponentOptions struct {
	ID   uintptr
	Size struct {
		Width  int32
		Height int32
	}
	Position struct {
		X int32
		Y int32
	}
	Visible bool
	Enabled bool
}

type CreateComponentOption func(*createComponentOptions)

func newCreateComponentOptions() *createComponentOptions {
	return &createComponentOptions{
		ID:       0,
		Size:     struct{ Width, Height int32 }{Width: 100, Height: 100},
		Position: struct{ X, Y int32 }{X: 0, Y: 0},
		Visible:  true,
		Enabled:  true,
	}
}

// ComponentIDOpt sets the ID of the component.
// It takes a uintptr as the ID and returns a CreateComponentOption function.
//
// Parameters:
//   - id: The unique identifier to set for the component.
//
// Returns:
//   - CreateComponentOption: A function that takes a pointer to createComponentOptions
func ComponentIDOpt(id uintptr) CreateComponentOption {
	return func(opts *createComponentOptions) {
		opts.ID = id
	}
}

// ComponentSizeOpt sets the size of the component.
// It takes width and height as int32 values and returns a CreateComponentOption function.
//
// Parameters:
//   - width: The width of the component.
//   - height: The height of the component.
//
// Returns:
//   - CreateComponentOption: A function that takes a pointer to createComponentOptions
func ComponentSizeOpt(width, height int32) CreateComponentOption {
	return func(opts *createComponentOptions) {
		opts.Size.Width = width
		opts.Size.Height = height
	}
}

// ComponentPositionOpt sets the position of the component.
// It takes x and y as int32 values and returns a CreateComponentOption function.
//
// Parameters:
//   - x: The x-coordinate of the component.
//   - y: The y-coordinate of the component.
//
// Returns:
//   - CreateComponentOption: A function that takes a pointer to createComponentOptions
func ComponentPositionOpt(x, y int32) CreateComponentOption {
	return func(opts *createComponentOptions) {
		opts.Position.X = x
		opts.Position.Y = y
	}
}

// ComponentVisibleOpt sets the visibility of the component.
// It takes a bool indicating visibility and returns a CreateComponentOption function.
//
// Parameters:
//   - visible: true to make the component visible, false to hide it.
//
// Returns:
//   - CreateComponentOption: A function that takes a pointer to createComponentOptions
func ComponentVisibleOpt(visible bool) CreateComponentOption {
	return func(opts *createComponentOptions) {
		opts.Visible = visible
	}
}

// ComponentEnabledOpt sets whether the component is enabled.
// It takes a bool indicating enabled state and returns a CreateComponentOption function.
//
// Parameters:
//   - enabled: true to enable the component, false to disable it.
//
// Returns:
//   - CreateComponentOption: A function that takes a pointer to createComponentOptions
func ComponentEnabledOpt(enabled bool) CreateComponentOption {
	return func(opts *createComponentOptions) {
		opts.Enabled = enabled
	}
}
