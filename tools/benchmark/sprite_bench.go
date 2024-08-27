package benchmark

import (
	"image"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

var (
	basepath  string
	logoPath  = "logo/LOGOTYPE-HORIZONTAL-BLUE2.png"
	logoFrame = pixel.R(98, 44, 234, 180)
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath = filepath.ToSlash(filepath.Dir(filepath.Dir(filepath.Dir(b))))
	logoPath = path.Join(basepath, logoPath)

	Benchmarks.Add(
		Config{
			Name:        "sprite-moving",
			Description: "Columns of sprites moving in opposite directions",
			New:         newSpriteMoving,
			Duration:    30 * time.Second,
		},
		Config{
			Name:        "sprite-moving-batched",
			Description: "Columns of sprites moving in opposite directions with batched draw",
			New:         newSpriteMovingBatched,
			Duration:    30 * time.Second,
		},
		Config{
			Name:        "sprite-static",
			Description: "Draw a sprite to the window in a grid",
			New:         newSpriteStatic,
			Duration:    30 * time.Second,
		},
		Config{
			Name:        "sprite-static-batched",
			Description: "Draw a sprite to the window in a grid with batched draw",
			New:         newSpriteStaticBatched,
			Duration:    30 * time.Second,
		},
	)
}

func newSpriteStatic(win *opengl.Window) (Benchmark, error) {
	sprite, err := loadSprite(logoPath, logoFrame)
	if err != nil {
		return nil, err
	}

	bounds := win.Bounds()
	width := bounds.W()
	height := bounds.H()
	rows, cols := 32, 32

	benchmark := &spriteStatic{
		sprite: sprite,
		rows:   rows,
		cols:   rows,
		cell:   gridCell(width, height, rows, cols),
	}
	return benchmark, nil
}

func newSpriteStaticBatched(win *opengl.Window) (Benchmark, error) {
	benchmark, err := newSpriteStatic(win)
	if err != nil {
		return nil, err
	}
	ss := benchmark.(*spriteStatic)
	ss.batch = pixel.NewBatch(&pixel.TrianglesData{}, ss.sprite.Picture())
	return ss, nil
}

type spriteStatic struct {
	sprite     *pixel.Sprite
	rows, cols int
	cell       pixel.Vec
	batch      *pixel.Batch
}

func (ss *spriteStatic) Step(win *opengl.Window, delta float64) {
	win.Clear(backgroundColor)
	var target pixel.Target
	if ss.batch != nil {
		ss.batch.Clear()
		target = ss.batch
	} else {
		target = win
	}
	spriteGrid(ss.sprite, target, ss.rows, ss.cols, ss.cell)
	if ss.batch != nil {
		ss.batch.Draw(win)
	}
}

func newSpriteMoving(win *opengl.Window) (Benchmark, error) {
	sprite, err := loadSprite(logoPath, logoFrame)
	if err != nil {
		return nil, err
	}
	bounds := win.Bounds()
	width := bounds.W()
	height := bounds.H()
	rows, cols := 32, 32
	benchmark := &spriteMoving{
		sprite: sprite,
		rows:   rows,
		cols:   cols,
		cell:   gridCell(width, height, rows, cols),
	}
	return benchmark, nil
}

func newSpriteMovingBatched(win *opengl.Window) (Benchmark, error) {
	benchmark, err := newSpriteMoving(win)
	if err != nil {
		return nil, err
	}
	sm := benchmark.(*spriteMoving)
	sm.batch = pixel.NewBatch(&pixel.TrianglesData{}, sm.sprite.Picture())
	return sm, nil
}

type spriteMoving struct {
	sprite     *pixel.Sprite
	batch      *pixel.Batch
	rows, cols int
	cell       pixel.Vec
	yOffset    float64
}

func (sm *spriteMoving) Step(win *opengl.Window, delta float64) {
	win.Clear(backgroundColor)
	var target pixel.Target
	if sm.batch != nil {
		sm.batch.Clear()
		target = sm.batch
	} else {
		target = win
	}

	sm.yOffset += sm.cell.Y * delta * 3
	if sm.yOffset >= sm.cell.Y {
		sm.yOffset = 0
	}

	spriteGridMoving(sm.sprite, target, sm.rows, sm.cols, sm.cell, sm.yOffset)
	if sm.batch != nil {
		sm.batch.Draw(win)
	}
}

func spriteGrid(sprite *pixel.Sprite, target pixel.Target, rows, cols int, cell pixel.Vec) {
	spriteBounds := sprite.Frame().Bounds()
	spriteWidth := spriteBounds.W()
	spriteHeight := spriteBounds.H()
	matrix := pixel.IM.ScaledXY(pixel.ZV, pixel.V(cell.X/spriteWidth, cell.Y/spriteHeight))
	offset := pixel.V(cell.X/2, cell.Y/2)
	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			pos := pixel.V(float64(i)*cell.X, float64(j)*cell.Y).Add(offset)
			sprite.Draw(target, matrix.Moved(pos))
		}
	}
}

func spriteGridMoving(sprite *pixel.Sprite, target pixel.Target, rows, cols int, cell pixel.Vec, yOffset float64) {
	spriteBounds := sprite.Frame().Bounds()
	spriteWidth := spriteBounds.W()
	spriteHeight := spriteBounds.H()
	matrix := pixel.IM.ScaledXY(pixel.ZV, pixel.V(cell.X/spriteWidth, cell.Y/spriteHeight))
	offset := pixel.V(cell.X/2, cell.Y/2)
	for i := 0; i < cols; i++ {
		columnOffset := -cell.Y
		if i%2 == 0 {
			columnOffset += yOffset
		} else {
			columnOffset -= yOffset
		}

		for j := 0; j < rows+2; j++ {
			pos := pixel.V(float64(i)*cell.X, (float64(j)*cell.Y)+columnOffset).Add(offset)
			sprite.Draw(target, matrix.Moved(pos))
		}
	}
}

func loadSprite(file string, frame pixel.Rect) (sprite *pixel.Sprite, err error) {
	image, err := loadPng(file)
	if err != nil {
		return nil, err
	}

	pic := pixel.PictureDataFromImage(image)
	if frame.Empty() {
		frame = pic.Bounds()
	}
	sprite = pixel.NewSprite(pic, frame)
	return sprite, nil
}

func loadPng(file string) (i image.Image, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	i, err = png.Decode(f)
	return
}
