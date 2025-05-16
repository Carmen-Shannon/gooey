package common

// Color represents an RGB color with red, green, and blue components.
// Each component is an 8-bit unsigned integer, allowing values from 0 to 255.
type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

// Adding a bunch of pre-set color values, can be used anywhere throughout the code.
var (
	ColorRed                  = &Color{Red: 255, Green: 0, Blue: 0}
	ColorGreen                = &Color{Red: 0, Green: 255, Blue: 0}
	ColorBlue                 = &Color{Red: 0, Green: 0, Blue: 255}
	ColorWhite                = &Color{Red: 255, Green: 255, Blue: 255}
	ColorBlack                = &Color{Red: 0, Green: 0, Blue: 0}
	ColorYellow               = &Color{Red: 255, Green: 255, Blue: 0}
	ColorCyan                 = &Color{Red: 0, Green: 255, Blue: 255}
	ColorMagenta              = &Color{Red: 255, Green: 0, Blue: 255}
	ColorGray                 = &Color{Red: 128, Green: 128, Blue: 128}
	ColorDarkGray             = &Color{Red: 64, Green: 64, Blue: 64}
	ColorLightGray            = &Color{Red: 192, Green: 192, Blue: 192}
	ColorOrange               = &Color{Red: 255, Green: 165, Blue: 0}
	ColorPurple               = &Color{Red: 128, Green: 0, Blue: 128}
	ColorPink                 = &Color{Red: 255, Green: 192, Blue: 203}
	ColorBrown                = &Color{Red: 165, Green: 42, Blue: 42}
	ColorLime                 = &Color{Red: 0, Green: 255, Blue: 0}
	ColorNavy                 = &Color{Red: 0, Green: 0, Blue: 128}
	ColorTeal                 = &Color{Red: 0, Green: 128, Blue: 128}
	ColorOlive                = &Color{Red: 128, Green: 128, Blue: 0}
	ColorCoral                = &Color{Red: 255, Green: 127, Blue: 80}
	ColorSalmon               = &Color{Red: 250, Green: 128, Blue: 114}
	ColorTurquoise            = &Color{Red: 64, Green: 224, Blue: 208}
	ColorViolet               = &Color{Red: 238, Green: 130, Blue: 238}
	ColorIndigo               = &Color{Red: 75, Green: 0, Blue: 130}
	ColorGold                 = &Color{Red: 255, Green: 215, Blue: 0}
	ColorKhaki                = &Color{Red: 240, Green: 230, Blue: 140}
	ColorPlum                 = &Color{Red: 221, Green: 160, Blue: 221}
	ColorSlateBlue            = &Color{Red: 106, Green: 90, Blue: 205}
	ColorSlateGray            = &Color{Red: 112, Green: 128, Blue: 144}
	ColorSteelBlue            = &Color{Red: 70, Green: 130, Blue: 180}
	ColorLightBlue            = &Color{Red: 173, Green: 216, Blue: 230}
	ColorLightGreen           = &Color{Red: 144, Green: 238, Blue: 144}
	ColorLightPink            = &Color{Red: 255, Green: 182, Blue: 193}
	ColorLightSalmon          = &Color{Red: 255, Green: 160, Blue: 122}
	ColorLightCoral           = &Color{Red: 240, Green: 128, Blue: 128}
	ColorLightGoldenrodYellow = &Color{Red: 250, Green: 250, Blue: 210}
	ColorLightSlateGray       = &Color{Red: 119, Green: 136, Blue: 153}
	ColorLightSteelBlue       = &Color{Red: 176, Green: 196, Blue: 222}
	ColorLightSeaGreen        = &Color{Red: 32, Green: 178, Blue: 170}
	ColorLightSkyBlue         = &Color{Red: 135, Green: 206, Blue: 250}
	ColorLightSlateBlue       = &Color{Red: 132, Green: 112, Blue: 255}
	ColorLightYellow          = &Color{Red: 255, Green: 255, Blue: 224}
	ColorMediumAquamarine     = &Color{Red: 102, Green: 205, Blue: 170}
	ColorMediumBlue           = &Color{Red: 0, Green: 0, Blue: 205}
	ColorMediumOrchid         = &Color{Red: 186, Green: 85, Blue: 211}
	ColorMediumPurple         = &Color{Red: 147, Green: 112, Blue: 219}
	ColorMediumSeaGreen       = &Color{Red: 60, Green: 179, Blue: 113}
	ColorMediumSlateBlue      = &Color{Red: 123, Green: 104, Blue: 238}
	ColorMediumSpringGreen    = &Color{Red: 0, Green: 250, Blue: 154}
	ColorMediumTurquoise      = &Color{Red: 72, Green: 209, Blue: 204}
	ColorMediumVioletRed      = &Color{Red: 199, Green: 21, Blue: 133}
	ColorMidnightBlue         = &Color{Red: 25, Green: 25, Blue: 112}
	ColorMistyRose            = &Color{Red: 255, Green: 228, Blue: 225}
	ColorMoccasin             = &Color{Red: 255, Green: 228, Blue: 181}
	ColorNavajoWhite          = &Color{Red: 255, Green: 222, Blue: 173}
	ColorOldLace              = &Color{Red: 253, Green: 245, Blue: 230}
	ColorOliveDrab            = &Color{Red: 107, Green: 142, Blue: 35}
	ColorOrangeRed            = &Color{Red: 255, Green: 69, Blue: 0}
	ColorOrchid               = &Color{Red: 218, Green: 112, Blue: 214}
	ColorPaleGoldenrod        = &Color{Red: 238, Green: 232, Blue: 170}
)
