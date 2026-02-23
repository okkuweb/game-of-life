package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/ui"
)

type options struct {
	width  int
	height int
	speed  int
}

type actionType int

type model struct {
	grid       gruid.Grid
	action     action
	heldAction actionType
	interval   time.Duration
	pause      bool
	opts       options
	ui         *ui.Label
	frame      gruid.Grid
	entities   map[gruid.Point]bool
	uiText     string
}

type action struct {
	Type     actionType
	Location gruid.Point
	Update   updateType
}

type updateType int

const (
	ColorPrimary gruid.Color = gruid.ColorDefault
	ColorLife gruid.Color = 64
	ColorUIBackgroundColor gruid.Color = 234
	ColorUIBorder gruid.Color = 125
)

func main() {
	if runtime.GOOS != "js" {
		InitLogger()
		defer logFile.Close()
		Log("Starting game")
	}
	opts := &options{width: 126, height: 33, speed: 200}
	var gd gruid.Grid
	if runtime.GOOS != "js" {
		gd = gruid.NewGrid(1000, 500)
	} else {
		gd = gruid.NewGrid(opts.width, opts.height)
	}
	entities := make(map[gruid.Point]bool)
	md := &model{grid: gd, pause: true, opts: *opts, entities: entities}

	md.ui = &ui.Label{
		Box: &ui.Box{
			Title: ui.Text(" Game of Life ").WithStyle(gruid.Style{Fg: ColorUIBorder, Bg: ColorUIBackgroundColor}),
			Style: gruid.Style{Fg: ColorUIBorder, Bg: ColorUIBackgroundColor},
		},
	}

	initDriver()

	app := gruid.NewApp(gruid.AppConfig{
		Driver: driver,
		Model:  md,
	})

	if err := app.Start(context.Background()); err != nil {
		fmt.Println(err)
	}
}

