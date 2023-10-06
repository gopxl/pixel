package main

import (
	"image/png"
	"os"
	"time"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
)

var gopherimg *pixel.Sprite
var imd *imdraw.IMDraw

var uTime, uSpeed float32

func gameloop(win *pixelgl.Window) {
	win.Canvas().SetUniform("uTime", &uTime)
	win.Canvas().SetUniform("uSpeed", &uSpeed)
	uSpeed = 5.0
	win.Canvas().SetFragmentShader(fragmentShader)

	start := time.Now()
	for !win.Closed() {
		win.Clear(pixel.RGB(0, 0, 0))
		gopherimg.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		uTime = float32(time.Since(start).Seconds())
		if win.Pressed(pixelgl.KeyRight) {
			uSpeed += 0.1
		}
		if win.Pressed(pixelgl.KeyLeft) {
			uSpeed -= 0.1
		}
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

in vec2 vTexCoords;
out vec4 fragColor;

uniform sampler2D uTexture;
uniform vec4 uTexBounds;

// custom uniforms
uniform float uSpeed;
uniform float uTime;

void main() {
    vec2 t = vTexCoords / uTexBounds.zw;
	vec3 influence = texture(uTexture, t).rgb;

    if (influence.r + influence.g + influence.b > 0.3) {
		t.y += cos(t.x * 40.0 + (uTime * uSpeed))*0.005;
		t.x += cos(t.y * 40.0 + (uTime * uSpeed))*0.01;
	}

    vec3 col = texture(uTexture, t).rgb;
	fragColor = vec4(col * vec3(0.6, 0.6, 1.2),1.0);
}
`
