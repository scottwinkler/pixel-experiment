package utility

import (
	"image/color"

	"github.com/faiface/pixel"
)

//convenience function for creating an rgba color from a color and an alpha
func ToRGBA(c color.Color, alpha float64) pixel.RGBA {
	//noAlpha := pixel.RGBA{R: 1, G: 1, B: 1, A: 0}
	rgba := pixel.ToRGBA(c)
	rgba.A = alpha
	return rgba
}
