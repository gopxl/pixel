package pixel

import "errors"

// ComposeTarget is a BasicTarget capable of Porter-Duff composition.
type ComposeTarget interface {
	BasicTarget

	// SetComposeMethod sets a Porter-Duff composition method to be used.
	SetComposeMethod(ComposeMethod)
}

// ComposeMethod is a Porter-Duff composition method.
type ComposeMethod int

// Here's the list of all available Porter-Duff composition methods. Use ComposeOver for the basic
// alpha blending.
const (
	ComposeOver ComposeMethod = iota
	ComposeIn
	ComposeOut
	ComposeAtop
	ComposeRover
	ComposeRin
	ComposeRout
	ComposeRatop
	ComposeXor
	ComposePlus
	ComposeCopy
	ComposeMultiply
	ComposeScreen
)

// Compose composes two colors together according to the ComposeMethod. A is the foreground, B is
// the background.
func (cm ComposeMethod) Compose(a, b RGBA) RGBA {
	var fa, fb float64

	switch cm {
	case ComposeOver:
		fa, fb = 1, 1-a.A
	case ComposeIn:
		fa, fb = b.A, 0
	case ComposeOut:
		fa, fb = 1-b.A, 0
	case ComposeAtop:
		fa, fb = b.A, 1-a.A
	case ComposeRover:
		fa, fb = 1-b.A, 1
	case ComposeRin:
		fa, fb = 0, a.A
	case ComposeRout:
		fa, fb = 0, 1-a.A
	case ComposeRatop:
		fa, fb = 1-b.A, a.A
	case ComposeXor:
		fa, fb = 1-b.A, 1-a.A
	case ComposePlus:
		fa, fb = 1, 1
	case ComposeCopy:
		fa, fb = 1, 0
	case ComposeMultiply:
		sa, da := a.A, b.A

		// term1: backdrop where source is transparent
		cr := b.R * (1 - sa)
		cg := b.G * (1 - sa)
		cb := b.B * (1 - sa)

		// term2: source where backdrop is transparent
		cr += a.R * (1 - da)
		cg += a.G * (1 - da)
		cb += a.B * (1 - da)

		// blended term: multiply of unpremultiplied colors, then re-premultiply by Sa*Da
		if sa > 0 && da > 0 {
			// unpremultiply
			asr, asg, asb := a.R/sa, a.G/sa, a.B/sa
			bsr, bsg, bsb := b.R/da, b.G/da, b.B/da
			cr += (asr * bsr) * sa * da
			cg += (asg * bsg) * sa * da
			cb += (asb * bsb) * sa * da
		}

		ao := sa + da - sa*da
		return RGBA{R: cr, G: cg, B: cb, A: ao}
	case ComposeScreen:
		sa, da := a.A, b.A

		// term1: backdrop where source is transparent
		cr := b.R * (1 - sa)
		cg := b.G * (1 - sa)
		cb := b.B * (1 - sa)

		// term2: source where backdrop is transparent
		cr += a.R * (1 - da)
		cg += a.G * (1 - da)
		cb += a.B * (1 - da)

		// blended term: screen of unpremultiplied colors, then re-premultiply by Sa*Da
		if sa > 0 && da > 0 {
			// unpremultiply
			asr, asg, asb := a.R/sa, a.G/sa, a.B/sa
			bsr, bsg, bsb := b.R/da, b.G/da, b.B/da
			// screen formula: 1 - (1-as)*(1-bs)
			cr += (1 - (1-asr)*(1-bsr)) * sa * da
			cg += (1 - (1-asg)*(1-bsg)) * sa * da
			cb += (1 - (1-asb)*(1-bsb)) * sa * da
		}

		ao := sa + da - sa*da
		return RGBA{R: cr, G: cg, B: cb, A: ao}
	default:
		panic(errors.New("Compose: invalid ComposeMethod"))
	}

	return a.Mul(Alpha(fa)).Add(b.Mul(Alpha(fb)))
}
