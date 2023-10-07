package main

import (
	"image/png"
	"os"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/pixelgl"
)

var gopherimg *pixel.Sprite

func gameloop(win *pixelgl.Window) {
	win.Canvas().SetFragmentShader(fragmentShader)

	for !win.Closed() {
		win.Clear(pixel.RGB(0, 0, 0))
		gopherimg.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 325, 348),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	f, err := os.Open("../assets/images/thegopherproject.png")
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	pd := pixel.PictureDataFromImage(img)
	gopherimg = pixel.NewSprite(pd, pd.Bounds())

	gameloop(win)
}

func main() {
	pixelgl.Run(run)
}

var fragmentShader = `
#version 330 core

in vec2  vTexCoords;

out vec4 fragColor;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

void main() {
	// Get our current screen coordinate
	vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;

	// Sum our 3 color channels
	float sum  = texture(uTexture, t).r;
	      sum += texture(uTexture, t).g;
	      sum += texture(uTexture, t).b;

	// Divide by 3, and set the output to the result
	vec4 color = vec4( sum/3, sum/3, sum/3, 1.0);
	fragColor = color;
}
`
