//go:build !sdl && !js

package main

import (
	"codeberg.org/anaseto/gruid"
	tcell "codeberg.org/anaseto/gruid-tcell"
	tc "github.com/gdamore/tcell/v2"
)

var driver gruid.Driver

func initDriver() {
	st := styler{}
	dr := tcell.NewDriver(tcell.Config{StyleManager: st})
	//dr.PreventQuit()
	driver = dr
}

// styler implements the tcell.StyleManager interface.
type styler struct{}

func (sty styler) GetStyle(cst gruid.Style) tc.Style {
	st := tc.StyleDefault
	cst.Fg = map16ColorTo256(cst.Fg, true)
	cst.Bg = map16ColorTo256(cst.Bg, false)
	st = st.Background(tc.ColorValid + tc.Color(cst.Bg)).Foreground(tc.ColorValid + tc.Color(cst.Fg))
	return st
}

const (
	Color256Base03  gruid.Color = 234
	Color256Base02  gruid.Color = 235
	Color256Base01  gruid.Color = 240
	Color256Base00  gruid.Color = 241 // for dark on light background
	Color256Base0   gruid.Color = 244
	Color256Base1   gruid.Color = 245
	Color256Base2   gruid.Color = 254
	Color256Base3   gruid.Color = 230
	Color256Yellow  gruid.Color = 136
	Color256Orange  gruid.Color = 166
	Color256Red     gruid.Color = 160
	Color256Magenta gruid.Color = 125
	Color256Violet  gruid.Color = 61
	Color256Blue    gruid.Color = 33
	Color256Cyan    gruid.Color = 37
	Color256Green   gruid.Color = 64
)

const (
	ColorFOV gruid.Color = Color256Base02
	ColorPlayer gruid.Color = Color256Blue
	ColorMonster gruid.Color = Color256Green
)

func map16ColorTo256(c gruid.Color, fg bool) gruid.Color {
	switch c {
	case ColorBackground:
		if fg {
			return Color256Base0
		}
		return Color256Base03
	case ColorBackgroundSecondary:
		return Color256Base02
	case ColorForegroundEmph:
		return Color256Base1
	case ColorForegroundSecondary:
		return Color256Base01
	case ColorYellow:
		return Color256Yellow
	case ColorOrange:
		return Color256Orange
	case ColorRed:
		return Color256Red
	case ColorMagenta:
		return Color256Magenta
	case ColorViolet:
		return Color256Violet
	case ColorBlue:
		return Color256Blue
	case ColorCyan:
		return Color256Cyan
	case ColorGreen:
		return Color256Green
	default:
		return c
	}
}
