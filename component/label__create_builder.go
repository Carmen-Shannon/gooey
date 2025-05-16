package component

import "github.com/Carmen-Shannon/gooey/common"

type createLabelOptions struct {
	Text             string
	Font             string
	Color            *common.Color
	TextSize         int32
	TextAlignment    TextAlignment
	WordWrap         bool
	ComponentOptions []CreateComponentOption
}

type CreateLabelOption func(*createLabelOptions)

func newCreateLabelOptions() *createLabelOptions {
	return &createLabelOptions{
		Text:          "Label",
		Font:          "Arial",
		Color:         &common.Color{Red: 0, Green: 0, Blue: 0},
		TextSize:      12,
		TextAlignment: CenterAlign,
		WordWrap:      false,
	}
}

// LabelTextOpt sets the text of the label.
//
// Parameters:
//   - text: The text to set for the label.
func LabelTextOpt(text string) CreateLabelOption {
	return func(opts *createLabelOptions) {
		opts.Text = text
	}
}

// LabelFontOpt sets the font of the label.
//
// Parameters:
//   - font: The font to set for the label.
func LabelFontOpt(font string) CreateLabelOption {
	return func(opts *createLabelOptions) {
		opts.Font = font
	}
}

// LabelColorOpt sets the color of the label.
//
// Parameters:
//   - color: The color to set for the label.
func LabelColorOpt(color *common.Color) CreateLabelOption {
	return func(opts *createLabelOptions) {
		opts.Color = color
	}
}

// LabelTextSizeOpt sets the text size of the label.
//
// Parameters:
//   - size: The text size to set for the label.
func LabelTextSizeOpt(size int32) CreateLabelOption {
	return func(opts *createLabelOptions) {
		opts.TextSize = size
	}
}

// LabelTextAlignmentOpt sets the text alignment of the label.
//
// Parameters:
//   - alignment: The text alignment to set for the label.
func LabelTextAlignmentOpt(alignment TextAlignment) CreateLabelOption {
	return func(opts *createLabelOptions) {
		opts.TextAlignment = alignment
	}
}

// LabelWordWrapOpt sets the word wrap option of the label.
//
// Parameters:
//   - wrap: The word wrap option to set for the label.
func LabelWordWrapOpt(wrap bool) CreateLabelOption {
	return func(opts *createLabelOptions) {
		opts.WordWrap = wrap
	}
}

// LabelComponentOptionsOpt sets the component options of the label.
//
// Parameters:
//   - options: The component options to set for the label.
func LabelComponentOptionsOpt(options ...CreateComponentOption) CreateLabelOption {
	return func(opts *createLabelOptions) {
		opts.ComponentOptions = options
	}
}
