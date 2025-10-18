package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// CustomDarkTheme provides a professional dark theme for the database client
type CustomDarkTheme struct{}

// Color returns custom colors for the dark theme
func (t *CustomDarkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 32, G: 33, B: 36, A: 255} // Darker, more refined background
	case theme.ColorNameButton:
		return color.RGBA{R: 48, G: 49, B: 54, A: 255}
	case theme.ColorNameDisabledButton:
		return color.RGBA{R: 40, G: 41, B: 45, A: 255}
	case theme.ColorNamePrimary:
		return color.RGBA{R: 0, G: 115, B: 230, A: 255} // Match the selection blue
	case theme.ColorNameHover:
		return color.RGBA{R: 55, G: 56, B: 62, A: 255}
	case theme.ColorNameFocus:
		return color.RGBA{R: 0, G: 115, B: 230, A: 255}
	case theme.ColorNamePressed:
		return color.RGBA{R: 42, G: 43, B: 48, A: 255}
	case theme.ColorNameSelection:
		return color.RGBA{R: 0, G: 115, B: 230, A: 255} // Solid selection color
	case theme.ColorNameSeparator:
		return color.RGBA{R: 50, G: 51, B: 56, A: 255}
	case theme.ColorNameForeground:
		return color.RGBA{R: 240, G: 240, B: 240, A: 255} // Slightly softer white
	case theme.ColorNameDisabled:
		return color.RGBA{R: 100, G: 100, B: 100, A: 255}
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 120, G: 120, B: 120, A: 255}
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 42, G: 43, B: 48, A: 255} // Slightly lighter input background
	case theme.ColorNameMenuBackground:
		return color.RGBA{R: 38, G: 39, B: 43, A: 255}
	case theme.ColorNameOverlayBackground:
		return color.RGBA{R: 28, G: 29, B: 32, A: 220}
	case theme.ColorNameShadow:
		return color.RGBA{R: 0, G: 0, B: 0, A: 100}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

// Font returns fonts for the theme
func (t *CustomDarkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon returns icons for the theme
func (t *CustomDarkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns custom sizes for the theme
func (t *CustomDarkTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 13 // Clean, readable text size
	case theme.SizeNameCaptionText:
		return 11
	case theme.SizeNameHeadingText:
		return 15
	case theme.SizeNameSubHeadingText:
		return 14
	case theme.SizeNamePadding:
		return 4 // Tighter padding for cleaner look
	case theme.SizeNameInlineIcon:
		return 16
	case theme.SizeNameScrollBar:
		return 10
	case theme.SizeNameScrollBarSmall:
		return 6
	default:
		return theme.DefaultTheme().Size(name)
	}
}
