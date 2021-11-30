package giodriver

import "image"

type Image struct {
	screen *Screen
	size   image.Point
	rgba   *image.RGBA
}

func (ii Image) Size() image.Point {
	return ii.size
}

func (ii Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, ii.size.X, ii.size.Y)
}

func (Image) Release() {}

func (ii Image) RGBA() *image.RGBA {
	return ii.rgba
}
