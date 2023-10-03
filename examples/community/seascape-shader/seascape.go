package main

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	// Set up window configs
	cfg := pixelgl.WindowConfig{ // Default: 1024 x 768
		Title:  "Golang GLSL",
		Bounds: pixel.R(0, 0, float64(width), float64(height)),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	camVector := win.Bounds().Center()

	bounds := win.Bounds()
	bounds.Max = bounds.Max.ScaledXY(pixel.V(1.0, 1.0))

	// I am putting all shader example initializing stuff here for
	// easier reference to those learning to use this functionality

	fragSource, err := LoadFileToString(filename)

	if err != nil {
		panic(err)
	}

	var uMouse mgl32.Vec4
	var uTime float32

	canvas := win.Canvas()
	uResolution := mgl32.Vec2{float32(win.Bounds().W()), float32(win.Bounds().H())}

	EasyBindUniforms(canvas,
		"uResolution", &uResolution,
		"uTime", &uTime,
		"uMouse", &uMouse,
		"uDrift", &uDrift,
	)

	canvas.SetFragmentShader(fragSource)

	start := time.Now()

	// Game Loop
	for !win.Closed() {
		uTime = float32(time.Since(start).Seconds())
		mpos := win.MousePosition()
		uMouse[0] = float32(mpos.X)
		uMouse[1] = float32(mpos.Y)

		win.Clear(colornames.Black)

		// Drawing to the screen
		canvas.Draw(win, pixel.IM.Moved(camVector))

		win.Update()
	}

}

func main() {
	parseFlags()

	pixelgl.Run(run)
}
