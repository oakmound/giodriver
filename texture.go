package giodriver

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/oakmound/oak/v3/shiny/screen"
)

type Texture struct {
	screen *Screen
	size   image.Point
	rgba   *image.RGBA
}

func (ti *Texture) Size() image.Point {
	return ti.size
}

func (ti *Texture) Bounds() image.Rectangle {
	return image.Rect(0, 0, ti.size.X, ti.size.Y)
}

func (ti *Texture) Upload(dp image.Point, src screen.Image, sr image.Rectangle) {
	rgba := src.RGBA()
	ti.rgba = rgba
}
func (*Texture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {}
func (*Texture) Release()                                             {}
