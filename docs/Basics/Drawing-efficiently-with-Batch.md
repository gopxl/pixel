[Previous Tutorial](./Pressing-keys-and-clicking-mouse.md)

In this part, we'll learn how to increase the performance of drawing using the [Batch](https://godoc.org/github.com/gopxl/pixel/v2#Batch).

## Previous part

In the [previous part](https://github.com/gopxl/pixel/wiki/Pressing-keys-and-clicking-mouse), we've created a pretty nice program for planting trees. But what happens if we plant a lot of trees? As it turns out, our program becomes quite slow, lagging and poorly responding.

So, let's start off with the code we've created and figure out how to make it more efficient.

```go
package main

import (
	"image"
	"math"
	"math/rand"
	"os"
	"time"

	_ "image/png"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"golang.org/x/image/colornames"
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
	cfg := opengl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	spritesheet, err := loadPicture("trees.png")
	if err != nil {
		panic(err)
	}

	var treesFrames []pixel.Rect
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += 32 {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += 32 {
			treesFrames = append(treesFrames, pixel.R(x, y, x+32, y+32))
		}
	}

	var (
		camPos       = pixel.ZV
		camSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
		trees        []*pixel.Sprite
		matrices     []pixel.Matrix
	)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		if win.JustPressed(pixel.MouseButtonLeft) {
			tree := pixel.NewSprite(spritesheet, treesFrames[rand.Intn(len(treesFrames))])
			trees = append(trees, tree)
			mouse := cam.Unproject(win.MousePosition())
			matrices = append(matrices, pixel.IM.Scaled(pixel.ZV, 4).Moved(mouse))
		}
		if win.Pressed(pixel.KeyLeft) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixel.KeyRight) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixel.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixel.KeyUp) {
			camPos.Y += camSpeed * dt
		}
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		win.Clear(colornames.Forestgreen)

		for i, tree := range trees {
			tree.Draw(win, matrices[i])
		}

		win.Update()
	}
}

func main() {
	opengl.Run(run)
}
```

## Measuring FPS

Before we really dive into the art of optimization, we need to be able to somehow measure how fast our program is. Optimizing without that knowledge is not really very reliable. Never optimize if you don't know what's slow.

A common way to measure the performance of a video game is the FPS, which stands for "frames per second". It's basically the number of times per second we manage call `win.Update`. With VSync enabled, this should be at most the refresh rate of your monitor, but if we disable it, we can get much higher numbers.

So first, let's disable VSync.

```go
	cfg := opengl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true, // delete this
	}
	win, err := opengl.NewWindow(cfg)
```

Now, let's go ahead and measure the number of frames per second. What we'll do is that we set up a counter. We're going to increment this counter by one each frame. Additionally, we set up a ticker from the standard `"time"` package to tick every second. This ticker sends a value on a channel every second. We're going to use the `select` statement to figure out when a second passes. When it does, we just display the accumulated number of frames in the window's title and reset the counter. Easy, right? Let's do it!

```go
	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	last := time.Now()
	for !win.Closed() {
```

First, we set up the variables. Variable `frames` is the counter and `second` is the channel, that emits a value every second. Let's use these variables at the end of the main loop.

```go
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
```

Here, we just increment the `frames` counter. Then, if a second passed, we put the FPS number in the title of the window and reset the counter. The `default` clause is very important here. Without it, we'd wait on the `select` statement until the `second` channel sends a value.

Let's run the program to see that this works!

[[images/05_drawing_efficiently_with_batch_fps.png]]

Now that we can measure the performance, let's plant many trees and see how the FPS goes down!

[[images/05_drawing_efficiently_with_batch_low_fps.png]]

At about 500 trees, the FPS drops below 30 on my computer, which is not too bad, but not too good either.

## OpenGL performance

Why is that? Well, PixelGL uses OpenGL to render the trees. In OpenGL, the performance often depends mostly on the number of "draw calls". What does that mean? Every draw operation to the window counts as a draw call. Thus, this cycle

```go
		for i, tree := range trees {
			tree.Draw(win, matrices[i])
		}
```

does more and more draw calls each frame as the number of trees increases. That becomes the bottleneck of our program. To increase the performance, we need to decrease the number of OpenGL draw calls. That's where the Batch comes in.

## Batch

[Batch](https://godoc.org/github.com/gopxl/pixel/v2#Batch) is a type that let's us accumulate many sprites (and other objects, but for now sprites) and draw them all at once using a single draw call. The only limitation is, that we can only use one [Picture](https://godoc.org/github.com/gopxl/pixel/v2#Picture) within one Batch. That's not a problem for us, since we use a spritesheet. All of our trees use the same picture.

Let's create our Batch.

```go
	spritesheet, err := loadPicture("trees.png")
	if err != nil {
		panic(err)
	}

	batch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
```

Now, that's not too simple. Let's break it down. The [NewBatch](https://godoc.org/github.com/gopxl/pixel/v2#NewBatch) constructor takes two arguments.

The first one is a little cryptic, it actually is the container that the Batch will use to accumulate the triangles (a sprite consists of two triangles) it will subsequently draw. The simplest and sufficient container here is [TrianglesData](https://godoc.org/github.com/gopxl/pixel/v2#TrianglesData), but we could use our own type of container if we wanted to. If you don't understand, don't worry too much, the theory behind that is not too important for basic usage.

The second argument is the picture we'll be using with the batch. If we attempted to draw a sprite with a different picture onto the batch, our program would panic.

Now, we need to draw our sprites to the batch. It's actually really simple.

```go
		batch.Clear()
		for i, tree := range trees {
			tree.Draw(batch, matrices[i])
		}
		batch.Draw(win)
```

First, we need to clear the batch. If we didn't do that, the trees would keep accumulating in the batch forever and we would run out of memory soon. Then, instead of drawing the sprites directly to the window, we draw them to the batch. This little magic is made possible thanks to the [Target](https://godoc.org/github.com/gopxl/pixel/v2#Target) interface. Both window and batch are targets, which makes it possible to draw sprites on them. Finally, we draw the accumulated sprites onto the window.

If we run the code now, we notice, that the FPS is not really going down. We can't click fast enough to create the amount of sprites necessary for our program to run out of breath. To test the performance of our program further, we need to be able to draw sprites faster in order to draw a lot more of them. So, let's make a tree brush!

Replace this

```go
		if win.JustPressed(pixel.MouseButtonLeft) {
```

with this

```go
		if win.Pressed(pixel.MouseButtonLeft) {
```

Now we only need to hold the left mouse button and move the mouse around the screen! A tree brush, wonderful!

[[images/05_drawing_efficiently_with_batch_batch.png]]

This time, I needed to draw more than 7000 trees in order to get the FPS down to 30. That's much better!

## Further optimization

We can actually optimize our program even further. See that the trees don't move? Why do we clear the batch and draw them all over again then? All we need to do is keep them in the batch and when the user presses the mouse, we draw another tree into the batch. This way, we don't have to clear the batch. Let's do that!

First, let's get rid of the `trees` slice, we won't need it any more.

```go
	var (
		camPos       = pixel.ZV
		camSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
		trees        []*pixel.Sprite // delete this
		matrices     []pixel.Matrix  // and this
	)
```

Now, instead of appending to the `trees` slice, we just draw to the batch.

```go
		if win.Pressed(pixel.MouseButtonLeft) {
			tree := pixel.NewSprite(spritesheet, treesFrames[rand.Intn(len(treesFrames))])
			mouse := cam.Unproject(win.MousePosition())
			tree.Draw(batch, pixel.IM.Scaled(pixel.ZV, 4).Moved(mouse))
		}
```

And finally, we draw the batch to the window.

```go
		win.Clear(colornames.Forestgreen)
		batch.Draw(win)
		win.Update()
```

Now the program does not seem to be running out of breath at all! Great!

Here's the whole code from this part.

```go
package main

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"os"
	"time"

	_ "image/png"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"golang.org/x/image/colornames"
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
	cfg := opengl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	spritesheet, err := loadPicture("trees.png")
	if err != nil {
		panic(err)
	}

	batch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)

	var treesFrames []pixel.Rect
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += 32 {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += 32 {
			treesFrames = append(treesFrames, pixel.R(x, y, x+32, y+32))
		}
	}

	var (
		camPos       = pixel.ZV
		camSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
	)

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		if win.Pressed(pixel.MouseButtonLeft) {
			tree := pixel.NewSprite(spritesheet, treesFrames[rand.Intn(len(treesFrames))])
			mouse := cam.Unproject(win.MousePosition())
			tree.Draw(batch, pixel.IM.Scaled(pixel.ZV, 4).Moved(mouse))
		}
		if win.Pressed(pixel.KeyLeft) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixel.KeyRight) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixel.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixel.KeyUp) {
			camPos.Y += camSpeed * dt
		}
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		win.Clear(colornames.Forestgreen)
		batch.Draw(win)
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	opengl.Run(run)
}
```

[Next Tutorial](./Drawing-shapes-with-IMDraw.md)
