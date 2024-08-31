package atlas

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

func generateImageGradient(bounds image.Rectangle, cTop, cBottom color.Color) *image.RGBA {
	img := image.NewRGBA(bounds)
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(float64(cTop.(color.RGBA).R)*(1-float64(y-bounds.Min.Y)/float64(bounds.Dy())) + float64(cBottom.(color.RGBA).R)*float64(y-bounds.Min.Y)/float64(bounds.Dy())),
				G: uint8(float64(cTop.(color.RGBA).G)*(1-float64(y-bounds.Min.Y)/float64(bounds.Dy())) + float64(cBottom.(color.RGBA).G)*float64(y-bounds.Min.Y)/float64(bounds.Dy())),
				B: uint8(float64(cTop.(color.RGBA).B)*(1-float64(y-bounds.Min.Y)/float64(bounds.Dy())) + float64(cBottom.(color.RGBA).B)*float64(y-bounds.Min.Y)/float64(bounds.Dy())),
				A: uint8(float64(cTop.(color.RGBA).A)*(1-float64(y-bounds.Min.Y)/float64(bounds.Dy())) + float64(cBottom.(color.RGBA).A)*float64(y-bounds.Min.Y)/float64(bounds.Dy())),
			})
		}
	}

	return img
}

func TestAtlas_Clear(t *testing.T) {
	i1 := generateImageGradient(image.Rect(0, 0, 10, 10), color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255})
	i2 := generateImageGradient(image.Rect(0, 0, 10, 10), color.RGBA{0, 0, 255, 255}, color.RGBA{255, 255, 0, 255})

	a := Atlas{}
	g1 := a.MakeGroup()

	// Add our two images to the atlas
	s1 := a.AddImage(i1)
	g1.AddImage(i2)

	a.Pack()

	// Remove one of the images through its group
	a.Clear(g1)

	// Now the atlas texture should be the same as the first image
	tex := a.internal[a.idMap[s1.id].index].Image()
	require.Equal(t, i1.Bounds(), tex.Bounds())

	for i := range i1.Pix {
		require.Equal(t, i1.Pix[i], tex.Pix[i])
	}
}

func TestAtlas_ClearAll(t *testing.T) {
	i1 := generateImageGradient(image.Rect(0, 0, 10, 10), color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255})
	i2 := generateImageGradient(image.Rect(0, 0, 10, 10), color.RGBA{0, 0, 255, 255}, color.RGBA{255, 255, 0, 255})

	a := Atlas{}

	// Add our two images to the atlas
	a.AddImage(i1)
	a.AddImage(i2)

	a.Pack()

	// Remove all of the images
	a.Clear()
}
