package main

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

// These hold the state of whether we're placing the first or second point of the line.
const (
	clickLineA = iota
	clickLineB
)

var (
	winBounds = pixel.R(0, 0, 1024, 768)

	r = pixel.R(10, 10, 70, 50)
	l = pixel.L(pixel.V(20, 20), pixel.V(100, 30))

	clickLine = clickLineA
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Line collision",
		Bounds: winBounds,
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	for !win.Closed() {
		win.Clear(color.RGBA{R: 23, G: 39, B: 58, A: 125})
		imd.Clear()

		// When mouse left-click, move the rectangle so its' center is at the mouse position
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			rectToMouse := r.Center().To(win.MousePosition())
			r = r.Moved(rectToMouse)
		}

		// When mouse right-click, set either the beginning or end of the line.
		if win.JustPressed(pixelgl.MouseButtonRight) {
			if clickLine == clickLineA {
				// Set the beginning of the line to the mouse position.
				// To make it clearer to the user, set the end position 1 pixel (in each direction) away from the first
				// point.
				l = pixel.L(win.MousePosition(), win.MousePosition().Add(pixel.V(1, 1)))
				clickLine = clickLineB
			} else {
				// Set the end point of the line.
				l = pixel.L(l.A, win.MousePosition())
				clickLine = clickLineA
			}
		}

		// Draw the rectangle.
		imd.Color = color.Black
		imd.Push(r.Min, r.Max)
		imd.Rectangle(3)

		// Draw the line.
		imd.Color = color.RGBA{R: 10, G: 10, B: 250, A: 255}
		imd.Push(l.A, l.B)
		imd.Line(3)

		imd.Color = color.RGBA{R: 250, G: 10, B: 10, A: 255}
		// Draw any intersection points.
		for _, i := range r.IntersectionPoints(l) {
			imd.Push(i)
			imd.Circle(4, 0)
		}

		imd.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
