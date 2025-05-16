package component

import "github.com/Carmen-Shannon/gooey/common"

type createTextInputOptions struct {
	Value            string
	MaxLength        int32
	Font             string
	Color            *common.Color
	TextColor        *common.Color
	TextSize         int32
	TextAlignment    TextAlignment
	ComponentOptions []CreateComponentOption
}

type CreateTextInputOption func(*createTextInputOptions)

func newCreateTextInputOptions() *createTextInputOptions {
	return &createTextInputOptions{
		Value:         "...",
		MaxLength:     3,
		Font:          "Arial",
		Color:         &common.Color{Red: 255, Green: 255, Blue: 255},
		TextColor:     &common.Color{Red: 0, Green: 0, Blue: 0},
		TextSize:      12,
		TextAlignment: CenterAlign,
	}
}

// TextInputValueOpt sets the value of the text input component.
// It takes a string as the value and returns a CreateTextInputOption function.
//
// Parameters:
//   - value: The string value to set for the text input component.
//
// Returns:
//   - CreateTextInputOption: A function that takes a pointer to createTextInputOptions
func TextInputValueOpt(value string) CreateTextInputOption {
	return func(opts *createTextInputOptions) {
		opts.Value = value
	}
}

// TextInputMaxLengthOpt sets the maximum length of the text input component.
// It takes an int32 as the maximum length and returns a CreateTextInputOption function.
//
// Parameters:
//   - maxLength: The maximum number of characters allowed in the text input component.
//
// Returns:
//   - CreateTextInputOption: A function that takes a pointer to createTextInputOptions
func TextInputMaxLengthOpt(maxLength int32) CreateTextInputOption {
	return func(opts *createTextInputOptions) {
		opts.MaxLength = maxLength
	}
}

// TextInputColorOpt sets the background color of the text input component.
// It takes a pointer to a common.Color and returns a CreateTextInputOption function.
//
// Parameters:
//   - color: The background color to set for the text input component.
//
// Returns:
//   - CreateTextInputOption: A function that takes a pointer to createTextInputOptions
func TextInputColorOpt(color *common.Color) CreateTextInputOption {
	return func(opts *createTextInputOptions) {
		opts.Color = color
	}
}

// TextInputTextColorOpt sets the text color of the text input component.
// It takes a pointer to a common.Color and returns a CreateTextInputOption function.
//
// Parameters:
//   - textColor: The text color to set for the text input component.
//
// Returns:
//   - CreateTextInputOption: A function that takes a pointer to createTextInputOptions
func TextInputTextColorOpt(textColor *common.Color) CreateTextInputOption {
	return func(opts *createTextInputOptions) {
		opts.TextColor = textColor
	}
}

// TextInputTextSizeOpt sets the text size of the text input component.
// It takes an int32 as the text size and returns a CreateTextInputOption function.
//
// Parameters:
//   - textSize: The size of the text in points.
//
// Returns:
//   - CreateTextInputOption: A function that takes a pointer to createTextInputOptions
func TextInputTextSizeOpt(textSize int32) CreateTextInputOption {
	return func(opts *createTextInputOptions) {
		opts.TextSize = textSize
	}
}

// TextInputTextAlignmentOpt sets the text alignment of the text input component.
// It takes a TextAlignment value and returns a CreateTextInputOption function.
//
// Parameters:
//   - textAlignment: The alignment to set for the text in the text input component.
//
// Returns:
//   - CreateTextInputOption: A function that takes a pointer to createTextInputOptions
func TextInputTextAlignmentOpt(textAlignment TextAlignment) CreateTextInputOption {
	return func(opts *createTextInputOptions) {
		opts.TextAlignment = textAlignment
	}
}

// TextInputComponentOptionsOpt sets additional component options for the text input component.
// It takes a variadic list of CreateComponentOption and returns a CreateTextInputOption function.
//
// Parameters:
//   - componentOptions: Additional options to apply to the text input component.
//
// Returns:
//   - CreateTextInputOption: A function that takes a pointer to createTextInputOptions
func TextInputComponentOptionsOpt(componentOptions ...CreateComponentOption) CreateTextInputOption {
	return func(opts *createTextInputOptions) {
		opts.ComponentOptions = componentOptions
	}
}
