package component

import (
	"gooey/common"
)

type button struct {
	baseComponent
	label      string
	labelFont  string
	labelColor *common.Color
	labelSize  int32
	onClick    func()
	pressed    bool
	hovered    bool
	bgDefault  *common.Color
	bgHover    *common.Color
	bgPressed  *common.Color
	bgDisabled *common.Color
	roundness  int32
}

// NewButton creates a new Button component with the specified options.
// It accepts a variadic list of CreateButtonOption functions to customize the button's properties.
//
// Parameters:
//   - options: A variadic list of CreateButtonOption functions to customize the button's properties.
//
// Returns:
//   - Button: A pointer to the newly created button component.
func NewButton(options ...CreateButtonOption) Button {
	opts := newCreateButtonOptions()
	for _, opt := range options {
		opt(opts)
	}
	cOpts := newCreateComponentOptions()
	for _, opt := range opts.ComponentOptions {
		opt(cOpts)
	}

	b := &button{
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
		label:      opts.Label,
		labelFont:  opts.LabelFont,
		labelColor: opts.LabelColor,
		labelSize:  opts.LabelSize,
		onClick:    opts.OnClick,
		pressed:    false,
		hovered:    false,
		bgDefault:  opts.BackgroundColor,
		bgHover:    opts.BackgroundColorHover,
		bgPressed:  opts.BackgroundColorPressed,
		bgDisabled: opts.BackgroundColorDisabled,
		roundness:  opts.Roundness,
	}

	cbMap := make(map[string]func(any))
	cbMap["onClick"] = func(_ any) {
		if b.onClick != nil {
			b.onClick()
		}
	}
	cbMap["pressed"] = func(p any) {
		b.pressed = p.(bool)
	}

	mapBtnCb(b, cbMap)
	registerBtnBounds(b)

	return b
}

type Button interface {
	Component

	// Label returns the label of the button.
	//
	// Returns:
	//  - string: The label of the button.
	Label() string

	// SetLabel sets the label of the button.
	//
	// Parameters:
	//  - label: The label to set for the button.
	SetLabel(label string)

	// LabelFont returns the font of the button label.
	//
	// Returns:
	//  - string: The font of the button label.
	LabelFont() string

	// SetLabelFont sets the font of the button label.
	//
	// Parameters:
	//  - labelFont: The font to set for the button label.
	//
	// Returns:
	//  - string: The font of the button label.
	SetLabelFont(labelFont string)

	// LabelColor returns the color of the button label.
	//
	// Returns:
	//  - *common.Color: The color of the button label.
	LabelColor() *common.Color

	// SetLabelColor sets the color of the button label.
	//
	// Parameters:
	//  - labelColor: The color to set for the button label.
	SetLabelColor(labelColor *common.Color)

	// LabelSize returns the size of the button label.
	//
	// Returns:
	//  - int32: The size of the button label.
	LabelSize() int32

	// SetLabelSize sets the size of the button label.
	//
	// Parameters:
	//  - labelSize: The size to set for the button label.
	SetLabelSize(labelSize int32)

	// OnClick returns the function to be called when the button is clicked.
	//
	// Returns:
	//  - func(): The function to be called on button click.
	OnClick() func()

	// SetOnClick sets the function to be called when the button is clicked.
	//
	// Parameters:
	//  - onClick: The function to set for button click.
	SetOnClick(onClick func())

	// Pressed returns whether the button is currently pressed.
	//
	// Returns:
	//  - bool: True if the button is pressed, false otherwise.
	Pressed() bool

	// SetPressed sets the pressed state of the button.
	//
	// Parameters:
	//  - pressed: The pressed state to set for the button.
	SetPressed(pressed bool)

	// Hovered returns whether the button is currently hovered by the mouse.
	//
	// Returns:
	//  - bool: True if the button is hovered, false otherwise.
	Hovered() bool

	// SetHovered sets the hovered state of the button.
	//
	// Parameters:
	//  - hovered: The hovered state to set for the button.
	SetHovered(hovered bool)

	// BackgroundColor returns the default background color of the button.
	//
	// Returns:
	//  - *common.Color: The default background color of the button.
	BackgroundColor() *common.Color

	// SetBackgroundColor sets the default background color of the button.
	//
	// Parameters:
	//  - backgroundColor: The color to set as the default background.
	SetBackgroundColor(backgroundColor *common.Color)

	// BackgroundColorHover returns the background color of the button when hovered.
	//
	// Returns:
	//  - *common.Color: The background color when hovered.
	BackgroundColorHover() *common.Color

	// SetBackgroundColorHover sets the background color of the button when hovered.
	//
	// Parameters:
	//  - backgroundColorHover: The color to set as the hover background.
	SetBackgroundColorHover(backgroundColorHover *common.Color)

	// BackgroundColorPressed returns the background color of the button when pressed.
	//
	// Returns:
	//  - *common.Color: The background color when pressed.
	BackgroundColorPressed() *common.Color

	// SetBackgroundColorPressed sets the background color of the button when pressed.
	//
	// Parameters:
	//  - backgroundColorPressed: The color to set as the pressed background.
	SetBackgroundColorPressed(backgroundColorPressed *common.Color)

	// BackgroundColorDisabled returns the background color of the button when disabled.
	//
	// Returns:
	//  - *common.Color: The background color when disabled.
	BackgroundColorDisabled() *common.Color

	// SetBackgroundColorDisabled sets the background color of the button when disabled.
	//
	// Parameters:
	//  - backgroundColorDisabled: The color to set as the disabled background.
	SetBackgroundColorDisabled(backgroundColorDisabled *common.Color)

	// Roundness returns the roundness (corner radius) of the button.
	//
	// Returns:
	//  - int32: The roundness value of the button.
	Roundness() int32

	// SetRoundness sets the roundness (corner radius) of the button.
	//
	// Parameters:
	//  - roundness: The roundness value to set for the button.
	SetRoundness(roundness int32)
}

var _ Button = (*button)(nil)

func (b *button) Draw(ctx *common.DrawCtx) {
	drawComponent(b, ctx)
}

func (b *button) Label() string {
	return b.label
}

func (b *button) SetLabel(label string) {
	b.label = label
}

func (b *button) OnClick() func() {
	return b.onClick
}

func (b *button) SetOnClick(onClick func()) {
	b.onClick = onClick
}

func (b *button) Pressed() bool {
	return b.pressed
}

func (b *button) SetPressed(pressed bool) {
	b.pressed = pressed
}

func (b *button) Hovered() bool {
	return b.hovered
}

func (b *button) SetHovered(hovered bool) {
	b.hovered = hovered
}

func (b *button) LabelFont() string {
	return b.labelFont
}

func (b *button) SetLabelFont(labelFont string) {
	b.labelFont = labelFont
}

func (b *button) LabelColor() *common.Color {
	return b.labelColor
}

func (b *button) SetLabelColor(labelColor *common.Color) {
	b.labelColor = labelColor
}

func (b *button) LabelSize() int32 {
	return b.labelSize
}

func (b *button) SetLabelSize(labelSize int32) {
	b.labelSize = labelSize
}

func (b *button) BackgroundColor() *common.Color {
	return b.bgDefault
}

func (b *button) SetBackgroundColor(backgroundColor *common.Color) {
	b.bgDefault = backgroundColor
}

func (b *button) BackgroundColorHover() *common.Color {
	return b.bgHover
}

func (b *button) SetBackgroundColorHover(backgroundColorHover *common.Color) {
	b.bgHover = backgroundColorHover
}

func (b *button) BackgroundColorPressed() *common.Color {
	return b.bgPressed
}

func (b *button) SetBackgroundColorPressed(backgroundColorPressed *common.Color) {
	b.bgPressed = backgroundColorPressed
}

func (b *button) BackgroundColorDisabled() *common.Color {
	return b.bgDisabled
}

func (b *button) SetBackgroundColorDisabled(backgroundColorDisabled *common.Color) {
	b.bgDisabled = backgroundColorDisabled
}

func (b *button) Roundness() int32 {
	return b.roundness
}

func (b *button) SetRoundness(roundness int32) {
	b.roundness = roundness
}
