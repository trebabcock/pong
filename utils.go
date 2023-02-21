package main

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

func MakeRectangle(position Vector2, col color.Color, width float64, height float64) *imdraw.IMDraw {
	imd := imdraw.New(nil)
	imd.Color = col
	imd.Push(pixel.V(position.X-(width/2), position.Y-(height/2))) // bottom left
	imd.Push(pixel.V(position.X+(width/2), position.Y-(height/2))) // bottom right
	imd.Push(pixel.V(position.X+(width/2), position.Y+(height/2))) // top right
	imd.Push(pixel.V(position.X-(width/2), position.Y+(height/2))) // top left
	imd.Polygon(0)

	return imd
}
