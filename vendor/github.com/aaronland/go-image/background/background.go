// Package background provides methods for manipulating an image's background.
package background

import (
	"image"
	"image/color"
	"image/draw"
)

// AddBackground draws 'im' on to a new `image.Image` instance of the same dimensions but with
// a background filled with 'background_colour'.
func AddBackground(im image.Image, bg_colour color.NRGBA) image.Image {

	new_im := image.NewNRGBA(im.Bounds())

	draw.Draw(new_im, new_im.Bounds(), image.NewUniform(bg_colour), image.Point{}, draw.Src)
	draw.Draw(new_im, new_im.Bounds(), im, im.Bounds().Min, draw.Over)

	return new_im
}
