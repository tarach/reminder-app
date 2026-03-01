package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// AppTheme provides a modern, warm color palette for the Reminders app.
type AppTheme struct {
	fyne.Theme
}

var _ fyne.Theme = (*AppTheme)(nil)

// App color palette - warm, modern, high contrast
var (
	// Background: medium gray for stronger contrast with white cards
	appBackground = color.NRGBA{R: 0xF8, G: 0xFA, B: 0xFC, A: 0xFF}
	// Card: pure white with subtle warmth
	appCardBackground = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	// Primary accent: teal - fresh and professional
	appPrimary = color.NRGBA{R: 0x0D, G: 0x94, B: 0x88, A: 0xFF}
	// Primary text: dark slate
	appForeground = color.NRGBA{R: 0x1E, G: 0x29, B: 0x3B, A: 0xFF}
	// Secondary text: muted gray
	appDisabled = color.NRGBA{R: 0x94, G: 0xA3, B: 0xB8, A: 0xFF}
	// Button: primary with good contrast
	appButton = color.NRGBA{R: 0x0D, G: 0x94, B: 0x88, A: 0xFF}
	// Day active: soft teal tint
	appDayActive = color.NRGBA{R: 0x14, G: 0xB8, B: 0xA6, A: 0xFF}
	// Day inactive: very light gray
	appDayInactive = color.NRGBA{R: 0xE2, G: 0xE8, B: 0xF0, A: 0xFF}
	// Card border: subtle
	appCardBorder = color.NRGBA{R: 0xE2, G: 0xE8, B: 0xF0, A: 0xFF}
)

var colorOverrides = map[fyne.ThemeColorName]color.Color{
	theme.ColorNameBackground:          appBackground,
	theme.ColorNameForeground:          appForeground,
	theme.ColorNameDisabled:            appDisabled,
	theme.ColorNamePrimary:             appPrimary,
	theme.ColorNameButton:              appButton,
	theme.ColorNameForegroundOnPrimary: color.White,
	theme.ColorNameInputBackground:     appCardBackground,
	theme.ColorNameInputBorder:         appCardBorder,
	theme.ColorNameHeaderBackground:    appBackground,
	theme.ColorNameMenuBackground:      appCardBackground,
	theme.ColorNameOverlayBackground:   appCardBackground,
	theme.ColorNameScrollBar:           color.NRGBA{R: 0x0D, G: 0x94, B: 0x88, A: 0x99},
	theme.ColorNameScrollBarBackground: appDayInactive,
	theme.ColorNameSeparator:           appCardBorder,
	theme.ColorNameShadow:              color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0x15},
}

func (t *AppTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	c, ok := colorOverrides[name]
	if ok {
		return c
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (t *AppTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *AppTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

var sizeOverrides = map[fyne.ThemeSizeName]float32{
	theme.SizeNameText:           14,
	theme.SizeNameHeadingText:    18,
	theme.SizeNameSubHeadingText: 16,
	theme.SizeNameCaptionText:    12,
	theme.SizeNamePadding:        12,
	theme.SizeNameInlineIcon:     20,
	theme.SizeNameScrollBar:      8,
	theme.SizeNameScrollBarSmall: 4,
}

func (t *AppTheme) Size(name fyne.ThemeSizeName) float32 {
	s, ok := sizeOverrides[name]
	if ok {
		return s
	}
	return theme.DefaultTheme().Size(name)
}
