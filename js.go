//go:build js

package main

import (
	"image"
	"image/color"
	"log"

	"codeberg.org/anaseto/gruid"
	js "codeberg.org/anaseto/gruid-js"
	"codeberg.org/anaseto/gruid/tiles"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
)

var driver gruid.Driver

func initDriver() {
	t, err := getTileDrawer()
	if err != nil {
		log.Fatal(err)
	}
	driver = js.NewDriver(js.Config{
		TileManager: t,
	})
}

type TileDrawer struct {
	drawer *tiles.Drawer
}

func getTileDrawer() (*TileDrawer, error) {
	t := &TileDrawer{}
	var err error
	// We get a monospace font TTF.
	font, err := opentype.Parse(gomono.TTF)
	if err != nil {
		return nil, err
	}
	// We retrieve a font face.
	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: 24,
		DPI:  72,
	})
	if err != nil {
		return nil, err
	}
	// We create a new drawer for tiles using the previous face. Note that
	// if more than one face is wanted (such as an italic or bold variant),
	// you would have to create drawers for thoses faces too, and then use
	// the relevant one accordingly in the GetImage method.
	t.drawer, err = tiles.NewDrawer(face)
	if err != nil {
		return nil, err
	}
	return t, nil
}

const (
	ColorBackground          gruid.Color = gruid.ColorDefault // background
	ColorBackgroundSecondary gruid.Color = 1 + 0              // black
)

func (t *TileDrawer) GetImage(c gruid.Cell) image.Image {
	// we use some selenized colors
	fg := image.NewUniform(color.RGBA{0xad, 0xbc, 0xbc, 255})
	bg := image.NewUniform(color.RGBA{0x18, 0x49, 0x56, 255})
	switch c.Style.Fg {
	case ColorBackgroundSecondary:
		fg = image.NewUniform(color.RGBA{0x46, 0x95, 0xf7, 255})
	}
	return t.drawer.Draw(c.Rune, fg, bg)
}

func (t *TileDrawer) TileSize() gruid.Point {
	return t.drawer.Size()
}
