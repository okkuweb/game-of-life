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
}

type action struct {
	Type     actionType
	Location gruid.Point
	Update   updateType
}

type updateType int

const (
	ColorBackground          gruid.Color = gruid.ColorDefault // background
	ColorBackgroundSecondary gruid.Color = 1 + 0              // black
	ColorForeground          gruid.Color = gruid.ColorDefault
	ColorForegroundSecondary gruid.Color = 1 + 7  // white
	ColorForegroundEmph      gruid.Color = 1 + 15 // bright white
	ColorRed                 gruid.Color = 1 + 9  // bright red
	ColorGreen               gruid.Color = 1 + 2
	ColorYellow              gruid.Color = 1 + 3
	ColorBlue                gruid.Color = 1 + 4
	ColorMagenta             gruid.Color = 1 + 5
	ColorCyan                gruid.Color = 1 + 6
	ColorOrange              gruid.Color = 1 + 1  // red
	ColorViolet              gruid.Color = 1 + 12 // bright blue
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
			Title: ui.Text(" Game of Life ").WithStyle(gruid.Style{Fg: ColorBlue}),
			Style: gruid.Style{Fg: ColorViolet},
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

