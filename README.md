# gooey #

This module is a GUI framework written in native Go. It requires Go 1.24+ to run.

## Installation ##
Import the module in a project:
```sh
go get github.com/Carmen-Shannon/gooey@latest
```

## Current Support ##
This module uses system calls to the Windows API for Windows UI's. There is plan to add support for Linux distro's and MacOS.

Right now the support is just for Windows and requires the `user32`, `kernel32` and `gdi32` DLLs (native to Windows 10+ installations)

## How To Get Set Up ##
The framework uses an option-builder pattern to create the main window and the child components, start by creating a main window that you want to add components to:
```go
package main

import (
    "github.com/Carmen-Shannon/gooey"
)

func main() {
    w := window.NewWindow(
		window.TitleOpt("Gooey Example"),
		window.WidthOpt(1200),
		window.HeightOpt(800),
		window.BackgroundColorOpt(&common.Color{Red: 79, Green: 71, Blue: 92}),
	)

    // Run MUST be called from the main function, as it locks to it's goroutine.
    w.Run(30) // FPS for the rendering, can customize to fit your needs
}
```

If you want to add components, they are customizable in the same builder-pattern as the main window:
```go
// ... from within the main function

func main() {
	w := window.NewWindow(
		window.TitleOpt("Gooey Example 69420"),
		window.WidthOpt(1200),
		window.HeightOpt(800),
		window.BackgroundColorOpt(&common.Color{Red: 79, Green: 71, Blue: 92}),
	)

	label := component.NewLabel(
		component.LabelTextOpt("Hello, Gooey!"),
		component.LabelColorOpt(common.ColorGreen),
		component.LabelTextSizeOpt(12),
		component.LabelFontOpt("Sans Serif"),
		component.LabelComponentOptionsOpt(
			component.ComponentIDOpt(uintptr(70)),
			component.ComponentSizeOpt(100, 25),
			component.ComponentPositionOpt(100, 150),
			component.ComponentVisibleOpt(false),
			component.ComponentEnabledOpt(true),
		),
	)

	sel := component.NewSelector(
		component.SelectorColorOpt(common.ColorLightGray),
		component.SelectorOpacityOpt(0.5),
		component.SelectorDrawingOpt(false),
		component.SelectorComponentOptionsOpt(
			component.ComponentIDOpt(uintptr(67)),
			component.ComponentVisibleOpt(false),
			component.ComponentEnabledOpt(true),
			component.ComponentSizeOpt(400, 400),
			component.ComponentPositionOpt(400, 400),
		),
	)

	btn := component.NewButton(
		component.ButtonLabelOpt("Click Me!"),
		component.ButtonLabelFontOpt("Arial"),
		component.ButtonLabelColorOpt(&common.Color{Red: 0, Green: 0, Blue: 0}),
		component.ButtonLabelSizeOpt(12),
		component.ButtonBackgroundColorOpt(&common.Color{Red: 240, Green: 240, Blue: 240}),
		component.ButtonBackgroundColorHoverOpt(&common.Color{Red: 200, Green: 200, Blue: 200}),
		component.ButtonBackgroundColorPressedOpt(&common.Color{Red: 252, Green: 3, Blue: 3}),
		component.ButtonRoundnessOpt(50),
		component.ButtonOnClickOpt(func() {
			label.SetVisible(!label.Visible())
			sel.StartCapture()
		}),
		component.ButtonComponentOptionsOpt(
			component.ComponentIDOpt(uintptr(69)),
			component.ComponentSizeOpt(100, 50),
			component.ComponentPositionOpt(100, 100),
			component.ComponentVisibleOpt(true),
			component.ComponentEnabledOpt(true),
		),
	)

	ti := component.NewTextInput(
		component.TextInputMaxLengthOpt(10),
		component.TextInputTextSizeOpt(10),
		component.TextInputTextAlignmentOpt(component.LeftAlign),
		component.TextInputTextColorOpt(&common.Color{Red: 0, Green: 0, Blue: 0}),
		component.TextInputColorOpt(&common.Color{Red: 240, Green: 240, Blue: 240}),
		component.TextInputValueOpt("Edit Me!"),
		component.TextInputComponentOptionsOpt(
			component.ComponentIDOpt(uintptr(68)),
			component.ComponentSizeOpt(100, 36),
			component.ComponentPositionOpt(100, 200),
			component.ComponentVisibleOpt(true),
			component.ComponentEnabledOpt(true),
		),
	)

	w.AddComponent(btn)
	w.AddComponent(label)
	w.AddComponent(ti)
	w.AddComponent(sel)
	w.Run(60)
}
```

Using the above code will get you set up with a pre-configured window with the three current component types rendered.

There are currently no formal docs written beyond the function definitions within each package in this repository.
