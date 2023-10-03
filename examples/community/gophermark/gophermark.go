package main

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"os"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type gopher struct {
	pos pixel.Vec
	vel pixel.Vec
	// angle float32
}

// newGopher creates a gopher with given position and random velocity.
func newGopher(p pixel.Vec) gopher {
	g := gopher{
		pos: p,
	}
	v := pixel.V(1, 0)
	v = v.Rotated(rand.Float64() * 2 * math.Pi)
	g.vel = v.Scaled(rand.Float64()*100 + 50)
	return g
}

// update updates the gophers position and ensures he stays on screen.
func (g *gopher) update(dt float64) {
	g.pos = g.pos.Add(g.vel.Scaled(dt))
	if g.pos.X <= pic.Bounds().W()/2 {
		g.vel.X *= -1
	} else if g.pos.X > win.Bounds().Max.X-pic.Bounds().W()/2 {
		g.vel.X *= -1
	}
	if g.pos.Y <= pic.Bounds().H()/2 {
		g.vel.Y *= -1
	} else if g.pos.Y > win.Bounds().Max.Y-pic.Bounds().H()/2 {
		g.vel.Y *= -1
	}
}

var (
	win     *pixelgl.Window
	pic     pixel.Picture
	sprite  *pixel.Sprite
	gophers []gopher
	frames  = 0
	second  = time.Tick(time.Second)
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {
	rand.Seed(time.Now().UnixNano())
	cfg := pixelgl.WindowConfig{
		Title:  "Gophermark",
		Bounds: pixel.R(0, 0, 1000, 800),
	}
	w, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win = w

	p, err := loadPicture("gopher.png")
	if err != nil {
		panic(err)
	}
	pic = p

	batch := pixel.NewBatch(&pixel.TrianglesData{}, pic)
	sprite = pixel.NewSprite(pic, pic.Bounds())

	for i := 0; i < 1000; i++ {
		gophers = append(gophers, newGopher(pixel.V(win.Bounds().W()/2, win.Bounds().H()/2)))
	}

	last := time.Now()
	for !win.Closed() && !win.JustPressed(pixelgl.KeyEscape) {
		dt := time.Since(last).Seconds()
		last = time.Now()

		if win.Pressed(pixelgl.MouseButtonLeft) {
			mouse := win.MousePosition()
			for i := 0; i < 10; i++ {
				gophers = append(gophers, newGopher(mouse))
			}
		}

		batch.Clear()
		for i := 0; i < len(gophers); i++ {
			gophers[i].update(dt)
			sprite.Draw(batch, pixel.IM.Moved(gophers[i].pos))
		}

		win.Clear(colornames.Black)
		batch.Draw(win)
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | Gophers: %d | FPS: %d", cfg.Title, len(gophers), frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
