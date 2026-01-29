//go:build js

package main

import (
	"image"

	"codeberg.org/anaseto/gruid"
	jsd "codeberg.org/anaseto/gruid-js"
)

var driver gruid.Driver

type monochromeTileManager struct{}

// TODO: Finish this function and add a tileset
func (tm *monochromeTileManager) GetImage(gc gruid.Cell) image.Image {
	return nil
}

func (tm *monochromeTileManager) TileSize() gruid.Point {
	return gruid.Point{X: 16, Y: 24}
}

func initDriver() {
	dr := jsd.NewDriver(jsd.Config{
		TileManager: &monochromeTileManager{},
		AppCanvasId: "appcanvas",
	})
	driver = dr
}
