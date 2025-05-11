package component

import (
	"gooey/common"
)

type createButtonOptions struct {
	Label                   string
	LabelFont               string
	LabelColor              *common.Color
	LabelSize               int32
	BackgroundColor         *common.Color
	BackgroundColorHover    *common.Color
	BackgroundColorPressed  *common.Color
	BackgroundColorDisabled *common.Color
	Roundness               int32
	OnClick                 func()
	ComponentOptions        []CreateComponentOption
}

type CreateButtonOption func(*createButtonOptions)

func newCreateButtonOptions() *createButtonOptions {
	return &createButtonOptions{
		Label:                   "Button",
		LabelFont:               "Arial",
		LabelColor:              &common.Color{Red: 0, Green: 0, Blue: 0},
		LabelSize:               12,
		BackgroundColor:         &common.Color{Red: 240, Green: 240, Blue: 240},
		BackgroundColorHover:    &common.Color{Red: 200, Green: 200, Blue: 200},
		BackgroundColorPressed:  &common.Color{Red: 150, Green: 150, Blue: 150},
		BackgroundColorDisabled: &common.Color{Red: 200, Green: 200, Blue: 200},
		Roundness:               0,
		OnClick:                 nil,
		ComponentOptions:        nil,
	}
}

// ButtonLabelOpt sets the label of the button.
//
// Parameters:
//   - label: The label to set for the button.
func ButtonLabelOpt(label string) CreateButtonOption {
	return func(opts *createButtonOptions) {
		opts.Label = label
	}
}

// ButtonLabelFontOpt sets the font of the button label.
//
// Parameters:
//   - font: The font to set for the button label.
func ButtonLabelFontOpt(font string) CreateButtonOption {
	return func(opts *createButtonOptions) {
		opts.LabelFont = font
	}
}

// ButtonLabelColorOpt sets the color of the button label.
//
// Parameters:
//   - color: The color to set for the button label.
func ButtonLabelColorOpt(color *common.Color) CreateButtonOption {
	return func(opts *createButtonOptions) {
		opts.LabelColor = color
	}
}

// ButtonLabelSizeOpt sets the size of the button label.
//
// Parameters:
//   - size: The size to set for the button label.
func ButtonLabelSizeOpt(size int32) CreateButtonOption {
	return func(opts *createButtonOptions) {
		opts.LabelSize = size
	}
}

// ButtonBackgroundColorOpt sets the background color of the button.
//
// Parameters:
//   - color: The background color to set for the button.
func ButtonBackgroundColorOpt(color *common.Color) CreateButtonOption {
	return func(opts *createButtonOptions) {
		opts.BackgroundColor = color
	}
}

// ButtonBackgroundColorHoverOpt sets the background color of the button when hovered.
//
// Parameters:
//   - color: The background color to set for the button when hovered.
func ButtonBackgroundColorHoverOpt(color *common.Color) CreateButtonOption {
	return func(opts *createButtonOptions) {
		opts.BackgroundColorHover = color
	}
}

// ButtonBackgroundColorPressedOpt sets the background color of the button when pressed.
//
// Parameters:
//   - color: The background color to set for the button when pressed.
func ButtonBackgroundColorPressedOpt(color *common.Color) CreateButtonOption {
	return func(opts *createButtonOptions) {
		opts.BackgroundColorPressed = color
	}
}

// ButtonBackgroundColorDisabledOpt sets the background color of the button when disabled.
//
// Parameters:
//   - color: The background color to set for the button when disabled.
func ButtonBackgroundColorDisabledOpt(color *common.Color) CreateButtonOption {
	return func(opts *createButtonOptions) {
		opts.BackgroundColorDisabled = color
	}
}

// ButtonRoundnessOpt sets the roundness of the button.
//
// Parameters:
//   - roundness: The roundness to set for the button, this is a value clamped between 0 and 100.
func ButtonRoundnessOpt(roundness int32) CreateButtonOption {
	return func(opts *createButtonOptions) {
		opts.Roundness = max(0, min(100, roundness))
	}
}

// ButtonOnClickOpt sets the function to be called when the button is clicked.
//
// Parameters:
//   - onClick: The function to set for the button click event.
func ButtonOnClickOpt(onClick func()) CreateButtonOption {
	return func(opts *createButtonOptions) {
		opts.OnClick = onClick
	}
}

// ButtonComponentOptionsOpt sets the component options for the button.
//
// Parameters:
//   - options: The component options to set for the button.
func ButtonComponentOptionsOpt(options ...CreateComponentOption) CreateButtonOption {
	return func(opts *createButtonOptions) {
		opts.ComponentOptions = options
	}
}
