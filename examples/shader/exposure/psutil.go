package main

import (
	"io/ioutil"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Pixel Shader utility functions

// EasyBindUniforms does all the work for you, just pass in a
// valid array adhering to format: String, Variable, ...
//
// example:
//
//   var uTimeVar float32
//   var uMouseVar mgl32.Vec4
//
//   EasyBindUniforms(win.GetCanvas(),
// 	     "uTime", &uTimeVar,
// 	     "uMouse", &uMouseVar,
//   )
//
func EasyBindUniforms(c *pixelgl.Canvas, unifs ...interface{}) {
	if len(unifs)%2 != 0 {
		panic("needs to be divisable by 2")
	}
	for i := 0; i < len(unifs); i += 2 {

		c.SetUniform(unifs[i+0].(string), unifs[i+1])
	}
}

// CenterWindow will... center the window
func CenterWindow(win *pixelgl.Window) {
	x, y := pixelgl.PrimaryMonitor().Size()
	width, height := win.Bounds().Size().XY()
	win.SetPos(
		pixel.V(
			x/2-width/2,
			y/2-height/2,
		),
	)
}

// LoadFileToString loads the contents of a file into a string
func LoadFileToString(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
