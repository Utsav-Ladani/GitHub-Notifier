package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type myTheme struct{}

var _ fyne.Theme = (*myTheme)(nil)

func (t *myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == "StatusUnread" {
		return color.RGBA{30, 144, 255, 255} // DodgerBlue
	}

	if name == "StatusRead" {
		return color.RGBA{167, 167, 168, 255} // Gray-White
	}

	if name == "Time" || name == "NType" {
		return color.RGBA{120, 120, 119, 255}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (t *myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
