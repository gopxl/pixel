package benchmark

import (
	"math"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
)

var (
	backgroundColor = pixel.RGB(0, 0, 0)
)

func init() {
	Benchmarks.Add(
		Config{
			Name:        "imdraw-static",
			Description: "Stationary RGB triangles in a grid",
			New:         newStaticTriangles,
			Duration:    30 * time.Second,
		},
		Config{
			Name:        "imdraw-moving",
			Description: "Columns of RGB triangles moving in opposite directions",
			New:         newMovingTriangles,
			Duration:    30 * time.Second,
		},
	)
}

func newStaticTriangles(win *opengl.Window) (Benchmark, error) {
	bounds := win.Bounds()
	width := bounds.W()
	height := bounds.H()
	rows, cols := 32, 32
	cell := gridCell(width, height, rows, cols)
	benchmark := &staticTriangles{
		imd:  tri(cell),
		rows: rows,
		cols: cols,
		cell: cell,
	}
	return benchmark, nil
}

type staticTriangles struct {
	imd        *imdraw.IMDraw
	rows, cols int
	cell       pixel.Vec
}

func (st *staticTriangles) Step(win *opengl.Window) {
	win.Clear(backgroundColor)

	for i := 0; i < st.cols; i++ {
		for j := 0; j < st.rows; j++ {
			pos := pixel.V(float64(i)*st.cell.X, float64(j)*st.cell.Y)
			win.SetMatrix(pixel.IM.Moved(pos))
			st.imd.Draw(win)
		}
	}
}

func newMovingTriangles(win *opengl.Window) (Benchmark, error) {
	bounds := win.Bounds()
	width := bounds.W()
	height := bounds.H()
	rows, cols := 32, 32
	cell := gridCell(width, height, rows, cols)
	benchmark := &movingTriangles{
		imd:  tri(cell),
		rows: rows,
		cols: cols,
		cell: cell,
	}
	return benchmark, nil
}

type movingTriangles struct {
	imd        *imdraw.IMDraw
	rows, cols int
	cell       pixel.Vec
	counter    int
}

func (mt *movingTriangles) Step(win *opengl.Window) {
	win.Clear(backgroundColor)

	for i := 0; i < mt.cols; i++ {
		yOffset := -mt.cell.Y
		delta := float64(mt.counter % int(mt.cell.Y))
		if i%2 == 0 {
			yOffset += delta
		} else {
			yOffset -= delta
		}

		for j := 0; j < mt.rows+2; j++ {
			pos := pixel.V(float64(i)*mt.cell.X, (float64(j)*mt.cell.Y)+yOffset)
			matrix := pixel.IM.Moved(pos)
			if i%2 == 1 {
				matrix = matrix.Rotated(pos.Add(pixel.V(mt.cell.X/2, mt.cell.Y/2)), math.Pi)
			}
			win.SetMatrix(matrix)
			mt.imd.Draw(win)
		}
	}

	mt.counter++
}

func tri(cell pixel.Vec) *imdraw.IMDraw {
	imd := imdraw.New(nil)
	imd.Color = pixel.RGB(1, 0, 0)
	imd.Push(pixel.V(0, 0))
	imd.Color = pixel.RGB(0, 1, 0)
	imd.Push(pixel.V(cell.X, 0))
	imd.Color = pixel.RGB(0, 0, 1)
	imd.Push(pixel.V(cell.X/2, cell.Y))
	imd.Polygon(0)
	return imd
}

func gridCell(width, height float64, rows, cols int) (cell pixel.Vec) {
	return pixel.V(width/float64(cols), height/float64(rows))
}
