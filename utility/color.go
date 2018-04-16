package utility

import (
	"image/color"

	"github.com/faiface/pixel"
)

//convenience function for creating an rgba color from a color and an alpha
func ToRGBA(c color.Color, alpha float64) pixel.RGBA {
	rgba := pixel.ToRGBA(c)
	return rgba.Mul(pixel.Alpha(alpha))
}
