package opengl

import (
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
)

type GLColor mgl32.Vec4

func (c GLColor) RGBA() (r, g, b, a uint32) {
	return uint32(c[0] * 0xffff), uint32(c[1] * 0xffff), uint32(c[2] * 0xffff), uint32(c[3] * 0xffff)
}

func ToGLColor(col color.Color) GLColor {
	if c, ok := col.(GLColor); ok {
		return c
	}
	r, g, b, a := col.RGBA()
	return GLColor{
		float32(r) / 0xffff,
		float32(g) / 0xffff,
		float32(b) / 0xffff,
		float32(a) / 0xffff,
	}
}
