# gooey

A native Go GUI framework that provides a simple, cross-platform way to create windows and GUI components using system calls. This library uses the option builder pattern to configure windows and components, supporting Windows and Linux (with plans for macOS).

Requires Go 1.24+.

## Installation

Import the module in your project:

```sh
go get github.com/Carmen-Shannon/gooey@latest
```

## Current Support

- **Windows**: Uses Windows API calls (requires `user32`, `kernel32`, and `gdi32` DLLs, native to Windows 10+)
- **Linux**: Uses X11 system calls (planned)
- **macOS**: Planned future support

## Project Structure

```
gooey/
├── main.go                    # Example usage
├── go.mod                     # Go module file
├── README.md                  # This documentation
├── common/                    # Shared types and utilities
│   ├── caret_ticker.go        # Caret blinking logic
│   ├── color.go               # Color definitions and predefined colors
│   ├── draw_context.go        # Drawing context utilities
│   ├── font.go                # Font handling
│   ├── highlighter.go         # Text highlighting
│   ├── rect.go                # Rectangle utilities
│   ├── selector_state.go      # Selector component state
│   └── text_input_state.go    # Text input component state
├── component/                 # GUI components
│   ├── button.go              # Button component interface
│   ├── button__create_builder.go  # Button creation options
│   ├── button_linux.go        # Linux-specific button implementation
│   ├── button_windows.go      # Windows-specific button implementation
│   ├── component.go           # Base component interface and types
│   ├── component__create_builder.go  # Base component creation options
│   ├── label.go               # Label component interface
│   ├── label__create_builder.go  # Label creation options
│   ├── label_linux.go         # Linux-specific label implementation
│   ├── label_windows.go       # Windows-specific label implementation
│   ├── selector.go            # Selector component interface
│   ├── selector__create_builder.go  # Selector creation options
│   ├── selector_linux.go      # Linux-specific selector implementation
│   ├── selector_windows.go    # Windows-specific selector implementation
│   ├── selector_windows_ext.go  # Extended Windows selector
│   ├── text_input.go          # Text input component interface
│   ├── text_input__create_builder.go  # Text input creation options
│   ├── text_input_linux.go    # Linux-specific text input implementation
│   └── text_input_windows.go  # Windows-specific text input implementation
├── internal/                  # Platform-specific internal implementations
│   ├── linux/                 # Linux-specific code
│   │   ├── linux.go
│   │   ├── linux_btn_events.go
│   │   └── linux_text_input_events.go
│   └── windows/               # Windows-specific code
│       ├── wdws.go
│       ├── wdws__create_window_builder.go
│       ├── wdws__get_message_builder.go
│       ├── wdws__register_class_builder.go
│       ├── wdws_btn_events.go
│       ├── wdws_draw_proc_events.go
│       ├── wdws_text_input_events.go
│       ├── wdws_window_proc_events.go
│       └── wdws_window_proc_events.go
└── window/                    # Window management
    ├── window.go              # Window interface
    ├── window__builder.go     # Window creation options
    ├── window__flags.go       # Window flags and constants
    ├── window_linux.go        # Linux-specific window implementation
    └── window_windows.go      # Windows-specific window implementation
```

## Main Usage

The framework uses an option builder pattern to create windows and components. Start by creating a main window, then add components to it.

### Creating a Window

```go
package main

import (
    "github.com/Carmen-Shannon/gooey"
    "github.com/Carmen-Shannon/gooey/common"
    "github.com/Carmen-Shannon/gooey/window"
)

func main() {
    w := window.NewWindow(
        window.TitleOpt("My GUI Application"),
        window.WidthOpt(1200),
        window.HeightOpt(800),
        window.BackgroundColorOpt(&common.Color{Red: 79, Green: 71, Blue: 92}),
    )

    // Add components here...

    // Run the window (blocks until closed)
    w.Run(60) // FPS for rendering
}
```

### Adding Components

Components are created using the builder pattern and added to the window:

```go
// Create a button
btn := component.NewButton(
    component.ButtonLabelOpt("Click Me!"),
    component.ButtonOnClickOpt(func() {
        // Handle click
    }),
    component.ButtonComponentOptionsOpt(
        component.ComponentSizeOpt(100, 50),
        component.ComponentPositionOpt(100, 100),
    ),
)

// Add to window
w.AddComponent(btn)
```

## Components

All components inherit base options from the `Component` interface. Each component type has its own specific options.

### Base Component Options

All components support these options via `ComponentOptionsOpt`:

- `ComponentIDOpt(id uintptr)` - Sets a unique identifier for the component
- `ComponentSizeOpt(width, height int32)` - Sets the component size in pixels
- `ComponentPositionOpt(x, y int32)` - Sets the component position in pixels
- `ComponentVisibleOpt(visible bool)` - Controls component visibility
- `ComponentEnabledOpt(enabled bool)` - Controls if the component is interactive

### Button Component

Interactive button with customizable appearance and click handler.

**Options:**

- `ButtonLabelOpt(label string)` - Sets the button text
- `ButtonLabelFontOpt(font string)` - Sets the font for the button text (e.g., "Arial")
- `ButtonLabelColorOpt(color *common.Color)` - Sets the text color
  - Type: `*common.Color` (pointer to RGB color struct)
- `ButtonLabelSizeOpt(size int32)` - Sets the text size in points
- `ButtonBackgroundColorOpt(color *common.Color)` - Sets the default background color
  - Type: `*common.Color` (pointer to RGB color struct)
- `ButtonBackgroundColorHoverOpt(color *common.Color)` - Sets the background color when hovered
  - Type: `*common.Color` (pointer to RGB color struct)
- `ButtonBackgroundColorPressedOpt(color *common.Color)` - Sets the background color when pressed
  - Type: `*common.Color` (pointer to RGB color struct)
- `ButtonBackgroundColorDisabledOpt(color *common.Color)` - Sets the background color when disabled
  - Type: `*common.Color` (pointer to RGB color struct)
- `ButtonRoundnessOpt(roundness int32)` - Sets corner roundness (0-100)
- `ButtonOnClickOpt(onClick func())` - Sets the click event handler function
- `ButtonComponentOptionsOpt(options ...CreateComponentOption)` - Sets base component options

**Example:**

```go
btn := component.NewButton(
    component.ButtonLabelOpt("Click Me!"),
    component.ButtonLabelFontOpt("Arial"),
    component.ButtonLabelColorOpt(common.ColorBlack),
    component.ButtonLabelSizeOpt(12),
    component.ButtonBackgroundColorOpt(common.ColorLightGray),
    component.ButtonBackgroundColorHoverOpt(common.ColorGray),
    component.ButtonBackgroundColorPressedOpt(common.ColorDarkGray),
    component.ButtonRoundnessOpt(10),
    component.ButtonOnClickOpt(func() {
        // Handle click
    }),
    component.ButtonComponentOptionsOpt(
        component.ComponentSizeOpt(100, 50),
        component.ComponentPositionOpt(100, 100),
        component.ComponentVisibleOpt(true),
        component.ComponentEnabledOpt(true),
    ),
)
```

### Label Component

Displays static text with customizable formatting.

**Options:**

- `LabelTextOpt(text string)` - Sets the label text
- `LabelFontOpt(font string)` - Sets the font (e.g., "Arial")
- `LabelColorOpt(color *common.Color)` - Sets the text color
  - Type: `*common.Color` (pointer to RGB color struct)
- `LabelTextSizeOpt(size int32)` - Sets the text size in points
- `LabelTextAlignmentOpt(alignment TextAlignment)` - Sets text alignment
  - Values: `LeftAlign`, `CenterAlign`, `RightAlign`
- `LabelWordWrapOpt(wrap bool)` - Enables/disables word wrapping
- `LabelComponentOptionsOpt(options ...CreateComponentOption)` - Sets base component options

**Example:**

```go
label := component.NewLabel(
    component.LabelTextOpt("Hello, World!"),
    component.LabelFontOpt("Arial"),
    component.LabelColorOpt(common.ColorBlack),
    component.LabelTextSizeOpt(12),
    component.LabelTextAlignmentOpt(component.CenterAlign),
    component.LabelWordWrapOpt(false),
    component.LabelComponentOptionsOpt(
        component.ComponentSizeOpt(200, 50),
        component.ComponentPositionOpt(100, 100),
        component.ComponentVisibleOpt(true),
        component.ComponentEnabledOpt(true),
    ),
)
```

### Selector Component

Interactive selection rectangle, useful for drawing or selection tools.

**Options:**

- `SelectorColorOpt(color *common.Color)` - Sets the selector color
  - Type: `*common.Color` (pointer to RGB color struct)
- `SelectorOpacityOpt(opacity float32)` - Sets opacity (0.0-1.0)
- `SelectorBoundsOpt(bounds common.Rect)` - Sets initial bounds
- `SelectorDrawingOpt(drawing bool)` - Enables/disables drawing mode
- `SelectorComponentOptionsOpt(options ...CreateComponentOption)` - Sets base component options

**Example:**

```go
sel := component.NewSelector(
    component.SelectorColorOpt(common.ColorBlue),
    component.SelectorOpacityOpt(0.8),
    component.SelectorDrawingOpt(true),
    component.SelectorComponentOptionsOpt(
        component.ComponentSizeOpt(200, 200),
        component.ComponentPositionOpt(100, 100),
        component.ComponentVisibleOpt(true),
        component.ComponentEnabledOpt(true),
    ),
)
```

### TextInput Component

Editable text field with customizable appearance.

**Options:**

- `TextInputValueOpt(value string)` - Sets the initial text value
- `TextInputMaxLengthOpt(maxLength int32)` - Sets maximum character count
- `TextInputColorOpt(color *common.Color)` - Sets background color
  - Type: `*common.Color` (pointer to RGB color struct)
- `TextInputTextColorOpt(textColor *common.Color)` - Sets text color
  - Type: `*common.Color` (pointer to RGB color struct)
- `TextInputTextSizeOpt(textSize int32)` - Sets text size in points
- `TextInputTextAlignmentOpt(textAlignment TextAlignment)` - Sets text alignment
  - Values: `LeftAlign`, `CenterAlign`, `RightAlign`
- `TextInputComponentOptionsOpt(options ...CreateComponentOption)` - Sets base component options

**Example:**

```go
ti := component.NewTextInput(
    component.TextInputValueOpt("Enter text..."),
    component.TextInputMaxLengthOpt(100),
    component.TextInputColorOpt(common.ColorWhite),
    component.TextInputTextColorOpt(common.ColorBlack),
    component.TextInputTextSizeOpt(12),
    component.TextInputTextAlignmentOpt(component.LeftAlign),
    component.TextInputComponentOptionsOpt(
        component.ComponentSizeOpt(200, 30),
        component.ComponentPositionOpt(100, 100),
        component.ComponentVisibleOpt(true),
        component.ComponentEnabledOpt(true),
    ),
)
```

## Window Options

Windows are created with the following options:

- `TitleOpt(title string)` - Sets the window title
- `WidthOpt(width int32)` - Sets window width in pixels
- `HeightOpt(height int32)` - Sets window height in pixels
- `StyleOpt(style uint32)` - Sets window style flags
- `ClassNameOpt(className string)` - Sets window class name
- `CloseChanOpt(closeChan chan struct{})` - Sets a channel to signal window close
- `BackgroundColorOpt(color *common.Color)` - Sets window background color
  - Type: `*common.Color` (pointer to RGB color struct)

## Colors

The framework provides predefined `*common.Color` values in the `common` package:

- `ColorRed`, `ColorGreen`, `ColorBlue`
- `ColorWhite`, `ColorBlack`
- `ColorYellow`, `ColorCyan`, `ColorMagenta`
- `ColorGray`, `ColorDarkGray`, `ColorLightGray`
- And many more...

You can also create custom colors:

```go
customColor := &common.Color{Red: 255, Green: 100, Blue: 50}
```

## Complete Example

```go
package main

import (
    "github.com/Carmen-Shannon/gooey"
    "github.com/Carmen-Shannon/gooey/common"
    "github.com/Carmen-Shannon/gooey/component"
    "github.com/Carmen-Shannon/gooey/window"
)

func main() {
    w := window.NewWindow(
        window.TitleOpt("Gooey Example"),
        window.WidthOpt(1200),
        window.HeightOpt(800),
        window.BackgroundColorOpt(&common.Color{Red: 79, Green: 71, Blue: 92}),
    )

    label := component.NewLabel(
        component.LabelTextOpt("Hello, Gooey!"),
        component.LabelColorOpt(common.ColorGreen),
        component.LabelTextSizeOpt(12),
        component.LabelFontOpt("Arial"),
        component.LabelComponentOptionsOpt(
            component.ComponentSizeOpt(200, 50),
            component.ComponentPositionOpt(100, 100),
            component.ComponentVisibleOpt(true),
        ),
    )

    btn := component.NewButton(
        component.ButtonLabelOpt("Toggle Label"),
        component.ButtonOnClickOpt(func() {
            label.SetVisible(!label.Visible())
        }),
        component.ButtonComponentOptionsOpt(
            component.ComponentSizeOpt(120, 40),
            component.ComponentPositionOpt(100, 200),
        ),
    )

    ti := component.NewTextInput(
        component.TextInputValueOpt("Type here..."),
        component.TextInputComponentOptionsOpt(
            component.ComponentSizeOpt(200, 30),
            component.ComponentPositionOpt(100, 300),
        ),
    )

    w.AddComponent(label)
    w.AddComponent(btn)
    w.AddComponent(ti)

    w.Run(60)
}
```

For more details, see the function documentation in each package.
