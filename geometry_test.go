package pixel_test

import (
	"testing"

	pixel "github.com/gopxl/pixel/v2"
)

type sub struct {
	result pixel.Vec
	t      float64
}

func TestBezier(t *testing.T) {
	tests := []struct {
		curve pixel.Bezier

		subTest []sub
		name    string
	}{
		{
			pixel.Constant(pixel.V(1, 0)),
			[]sub{
				{pixel.V(1, 0), 0.0},
				{pixel.V(1, 0), 100.0},
			},
			"constant",
		},
		{
			pixel.Linear(pixel.V(1, 0), pixel.ZV),
			[]sub{
				{pixel.V(1, 0), 0.0},
				{pixel.ZV, 1.0},
			},
			"lenear",
		},
		{
			pixel.B(pixel.V(0, 1), pixel.V(1, 0), pixel.V(-1, 0), pixel.V(1, 0)),
			[]sub{
				{pixel.V(0, 1), 0.0},
				{pixel.V(1, 0), 1.0},
				{pixel.V(.5, .5), 0.5},
			},
			"curved",
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			for _, st := range c.subTest {
				val := c.curve.Point(st.t)
				if val != st.result {
					t.Errorf("inputted: %v expected: %v got: %v", st.t, st.result, val)
				}
			}
		})
	}
}
