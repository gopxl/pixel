package pixel

import (
	"image/color"

	"golang.org/x/exp/constraints"
)

// RGBA represents an alpha-premultiplied RGBA color with components within range [0, 1].
//
// The difference between color.RGBA is that the value range is [0, 1] and the values are floats.
type RGBA color.RGBA

// RGB returns a fully opaque RGBA color with the given RGB values.
//
// A common way to construct a transparent color is to create one with RGB constructor, then
// multiply it by a color obtained from the Alpha constructor.
func RGB(r, g, b float64) RGBA {
	return RGBA{R: uint8(r * 255), G: uint8(g * 255), B: uint8(b * 255), A: 255}
}

// Alpha returns a white RGBA color with the given alpha component.
func Alpha(a float64) RGBA {
	A := uint8(a * 255)
	return RGBA{A, A, A, A}
}

// Add adds color d to color c component-wise and returns the result (the components are not
// clamped).
func (c RGBA) Add(d color.Color) RGBA {
	rgba := ToRGBA(d)
	return RGBA{
		R: c.R + rgba.R,
		G: c.G + rgba.G,
		B: c.B + rgba.B,
		A: c.A + rgba.A,
	}
}

// Sub subtracts color d from color c component-wise and returns the result (the components
// are not clamped).
func (c RGBA) Sub(d color.Color) RGBA {
	rgba := ToRGBA(d)
	return RGBA{
		R: c.R - rgba.R,
		G: c.G - rgba.G,
		B: c.B - rgba.B,
		A: c.A - rgba.A,
	}
}

// Mul multiplies color c by color d component-wise (the components are not clamped).
func (c RGBA) Mul(d color.Color) RGBA {
	r1, g1, b1, a1 := ColorToFloats[float64](c)
	r2, g2, b2, a2 := ColorToFloats[float64](d)
	return FloatsToColor(r1*r2, g1*g2, b1*b2, a1*a2)
}

// Scaled multiplies each component of color c by scale and returns the result.
func (c RGBA) Scaled(scale float64) RGBA {
	r, g, b, a := ColorToFloats[float64](c)
	return RGBA{
		R: uint8(r * scale),
		G: uint8(g * scale),
		B: uint8(b * scale),
		A: uint8(a * scale),
	}
}

// RGBA returns components of the color.
func (c RGBA) RGBA() (r, g, b, a uint32) {
	return color.RGBA(c).RGBA()
}

// ColorToFloats converts a color to float32 or float64 components.
func ColorToFloats[F constraints.Float](c color.Color) (r, g, b, a F) {
	r1, g1, b1, a1 := c.RGBA()
	return F(r1) / 0xffff, F(g1) / 0xffff, F(b1) / 0xffff, F(a1) / 0xffff
}

// FloatsToColor converts float32 or float64 components to a color.
func FloatsToColor[F constraints.Float](r, g, b, a F) RGBA {
	return RGBA{uint8(r * F(255)), uint8(g * F(255)), uint8(b * F(255)), uint8(a * F(255))}
}

// ToRGBA converts a color to RGBA format.
func ToRGBA(c color.Color) RGBA {
	switch c := c.(type) {
	case RGBA:
		return c
	case color.RGBA:
		return RGBA(c)
	default:
		return RGBA(color.RGBAModel.Convert(c).(color.RGBA))
	}
}

// RGBAModel converts colors to RGBA format.
var RGBAModel = color.ModelFunc(rgbaModel)

func rgbaModel(c color.Color) color.Color {
	return ToRGBA(c)
}
