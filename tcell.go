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
	cst.Fg = mapColors(cst.Fg, true)
	cst.Bg = mapColors(cst.Bg, false)
	st = st.Background(tc.ColorValid + tc.Color(cst.Bg)).Foreground(tc.ColorValid + tc.Color(cst.Fg))
	return st
}

func mapColors(c gruid.Color, fg bool) gruid.Color {
	if c == ColorPrimary && fg {
		return 250
	} else if c == ColorPrimary {
		return 233
	} else {
		return c
	}
}

