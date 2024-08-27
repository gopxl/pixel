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
			Name:        "imdraw-static-batched",
			Description: "Stationary RGB triangles in a grid with batched draw",
			New:         newStaticTrianglesBatched,
			Duration:    30 * time.Second,
		},
		Config{
			Name:        "imdraw-moving",
			Description: "Columns of RGB triangles moving in opposite directions",
			New:         newMovingTriangles,
			Duration:    30 * time.Second,
		},
		Config{
			Name:        "imdraw-moving-batched",
			Description: "Columns of RGB triangles moving in opposite directions with batched draw",
			New:         newMovingTrianglesBatched,
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

func newStaticTrianglesBatched(win *opengl.Window) (Benchmark, error) {
	benchmark, err := newStaticTriangles(win)
	if err != nil {
		return nil, err
	}
	st := benchmark.(*staticTriangles)
	st.target = pixel.NewBatch(&pixel.TrianglesData{}, nil)
	return st, nil
}

type staticTriangles struct {
	imd        *imdraw.IMDraw
	batch      *pixel.Batch
	target     pixel.BasicTarget
	rows, cols int
	cell       pixel.Vec
}

func (st *staticTriangles) Step(win *opengl.Window, delta float64) {
	win.Clear(backgroundColor)

	var target pixel.BasicTarget
	if st.batch != nil {
		st.batch.Clear()
		target = st.batch
	} else {
		target = win
	}

	for i := 0; i < st.cols; i++ {
		for j := 0; j < st.rows; j++ {
			pos := pixel.V(float64(i)*st.cell.X, float64(j)*st.cell.Y)
			target.SetMatrix(pixel.IM.Moved(pos))
			st.imd.Draw(target)
		}
	}

	if st.batch != nil {
		st.batch.Draw(win)
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

func newMovingTrianglesBatched(win *opengl.Window) (Benchmark, error) {
	benchmark, err := newMovingTriangles(win)
	if err != nil {
		return nil, err
	}

	mt := benchmark.(*movingTriangles)
	mt.batch = pixel.NewBatch(&pixel.TrianglesData{}, nil)
	return mt, nil
}

type movingTriangles struct {
	imd        *imdraw.IMDraw
	batch      *pixel.Batch
	rows, cols int
	cell       pixel.Vec
	yOffset    float64
}

func (mt *movingTriangles) Step(win *opengl.Window, delta float64) {
	win.Clear(backgroundColor)

	var target pixel.BasicTarget
	if mt.batch != nil {
		mt.batch.Clear()
		target = mt.batch
	} else {
		target = win
	}

	mt.yOffset += mt.cell.Y * delta * 3
	if mt.yOffset >= mt.cell.Y {
		mt.yOffset = 0
	}

	for i := 0; i < mt.cols; i++ {
		columnOffset := -mt.cell.Y
		if i%2 == 0 {
			columnOffset += mt.yOffset
		} else {
			columnOffset -= mt.yOffset
		}

		for j := 0; j < mt.rows+2; j++ {
			pos := pixel.V(float64(i)*mt.cell.X, (float64(j)*mt.cell.Y)+columnOffset)
			matrix := pixel.IM.Moved(pos)
			if i%2 == 1 {
				matrix = matrix.Rotated(pos.Add(pixel.V(mt.cell.X/2, mt.cell.Y/2)), math.Pi)
			}
			target.SetMatrix(matrix)
			mt.imd.Draw(target)
		}
	}

	if mt.batch != nil {
		mt.batch.Draw(win)
	}
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
