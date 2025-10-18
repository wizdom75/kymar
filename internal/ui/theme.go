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
		return color.RGBA{R: 40, G: 44, B: 52, A: 255} // Dark blue-gray background
	case theme.ColorNameButton:
		return color.RGBA{R: 60, G: 64, B: 72, A: 255} // Slightly lighter for buttons
	case theme.ColorNameDisabledButton:
		return color.RGBA{R: 50, G: 54, B: 62, A: 255}
	case theme.ColorNamePrimary:
		return color.RGBA{R: 0, G: 122, B: 255, A: 255} // Blue accent like macOS
	case theme.ColorNameHover:
		return color.RGBA{R: 70, G: 74, B: 82, A: 255}
	case theme.ColorNameFocus:
		return color.RGBA{R: 0, G: 122, B: 255, A: 255}
	case theme.ColorNamePressed:
		return color.RGBA{R: 50, G: 54, B: 62, A: 255}
	case theme.ColorNameSelection:
		return color.RGBA{R: 0, G: 122, B: 255, A: 100} // Semi-transparent blue
	case theme.ColorNameSeparator:
		return color.RGBA{R: 60, G: 64, B: 72, A: 255}
	case theme.ColorNameForeground:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // White text
	case theme.ColorNameDisabled:
		return color.RGBA{R: 128, G: 128, B: 128, A: 255} // Gray for disabled
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 160, G: 160, B: 160, A: 255} // Light gray for placeholders
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 50, G: 54, B: 62, A: 255} // Darker input background
	case theme.ColorNameMenuBackground:
		return color.RGBA{R: 45, G: 49, B: 57, A: 255}
	case theme.ColorNameOverlayBackground:
		return color.RGBA{R: 35, G: 39, B: 47, A: 200}
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
		return 13 // Slightly smaller text for professional look
	case theme.SizeNameCaptionText:
		return 11
	case theme.SizeNameHeadingText:
		return 16
	case theme.SizeNameSubHeadingText:
		return 14
	case theme.SizeNamePadding:
		return 6
	case theme.SizeNameInlineIcon:
		return 16
	case theme.SizeNameScrollBar:
		return 12
	case theme.SizeNameScrollBarSmall:
		return 8
	default:
		return theme.DefaultTheme().Size(name)
	}
}
