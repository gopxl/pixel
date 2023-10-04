package pixel

// Bezier is cubic Bézier curve used for interpolation. For more info
// see https://en.wikipedia.org/wiki/B%C3%A9zier_curve,
// In case you are looking for visualization see https://www.desmos.com/calculator/d1ofwre0fr
type Bezier struct {
	Start, StartHandle, EndHandle, End Vec
	redundant                          bool
}

// ZB is Zero Bezier Curve that skips calculation and always returns V(1, 0)
// Its mainly because Calculation uses lot of function calls and in case of
// particles, it can make some difference
var ZB = Constant(V(1, 0))

// B returns new curve. if curve is just placeholder use constant. Handles are
// relative to start and end point so:
//
// pixel.B(ZV, ZV, ZV, V(1, 0)) == Bezier{ZV, ZV, V(1, 0), V(1, 0)}
func B(start, startHandle, endHandle, end Vec) Bezier {
	return Bezier{start, startHandle.Add(start), endHandle.Add(end), end, false}
}

// Linear returns linear Bezier curve
func Linear(start, end Vec) Bezier {
	return B(start, ZV, ZV, end)
}

// Constant returns Bezier curve that always return same point,
// This is usefull as placeholder, because it skips calculation
func Constant(constant Vec) Bezier {
	return Bezier{
		Start:     constant,
		redundant: true,
	}
}

// Point returns point along the curve determinate by t (0 - 1)
// You can of course pass any value though its really hard to
// predict what value will it return
func (b Bezier) Point(t float64) Vec {
	if b.redundant || b.Start == b.End {
		b.redundant = true
		return b.Start
	}

	inv := 1.0 - t
	c, d, e, f := inv*inv*inv, inv*inv*t*3.0, inv*t*t*3.0, t*t*t

	return V(
		b.Start.X*c+b.StartHandle.X*d+b.EndHandle.X*e+b.End.X*f,
		b.Start.Y*c+b.StartHandle.Y*d+b.EndHandle.Y*e+b.End.Y*f,
	)
}
