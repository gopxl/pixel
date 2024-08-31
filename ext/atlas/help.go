package atlas

import (
	"image"

	// need the following to automatically register for image.decode
	_ "image/jpeg"
	_ "image/png"

	"github.com/gopxl/pixel/v2"
	"golang.org/x/exp/constraints"
)

func area(r image.Rectangle) int {
	return r.Dx() * r.Dy()
}

func rect(x, y, w, h int) image.Rectangle {
	return image.Rect(x, y, x+w, y+h)
}

func pixelRect[T constraints.Integer | constraints.Float](minX, minY, maxX, maxY T) pixel.Rect {
	return pixel.R(float64(minX), float64(minY), float64(maxX), float64(maxY))
}

func image2PixelRect(r image.Rectangle) pixel.Rect {
	return pixelRect(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)
}

// split is the actual algorithm for splitting a given space (by j in spcs) to fit the given width and height.
// Will return an empty rectangle if a space wasn't available
// This function is based on this project (https://github.com/TeamHypersomnia/rectpack2D)
func split(spcs spaces, j int, bw, bh int) (found image.Rectangle, newSpcs spaces) {
	sp := spcs[j]
	spw, sph := sp.Dx(), sp.Dy()
	switch {
	// Perfect match
	case bw == spw && bh == sph:
		found = sp
		spcs = append(spcs[:j], spcs[j+1:]...)
	// Perfect width, split height
	case bw == spw && bh < sph:
		h := sph - bh
		found = rect(sp.Min.X, sp.Min.Y, spw, bh)
		spcs = append(spcs[:j], spcs[j+1:]...)
		spcs = append(spcs, rect(sp.Min.X, sp.Min.Y+bh, spw, h))
	// Perfect height, split width
	case bw < spw && bh == sph:
		w := spw - bw
		found = rect(sp.Min.X, sp.Min.Y, bw, sph)
		spcs = append(spcs[:j], spcs[j+1:]...)
		spcs = append(spcs, rect(sp.Min.X+bw, sp.Min.Y, w, sph))
	// Split both
	case bw < spw && bh < sph:
		w := spw - bw
		h := sph - bh
		found = rect(sp.Min.X, sp.Min.Y, bw, bh)
		var r1, r2 image.Rectangle

		// Maximize the leftover size
		r1 = rect(sp.Min.X+bw, sp.Min.Y, w, bh)
		r2 = rect(sp.Min.X, sp.Min.Y+bh, spw, h)

		spcs = append(spcs[:j], spcs[j+1:]...)
		spcs = append(spcs, r1, r2)
	}
	newSpcs = spcs
	return
}
