package pixel

import (
	"errors"
	"image/color"
)

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
)

// Compose composes two colors together according to the ComposeMethod. A is the foreground, B is
// the background.
func (cm ComposeMethod) Compose(a, b color.Color) RGBA {
	var fa, fb float64

	ac := ToRGBA(a)
	bc := ToRGBA(b)

	aa := float64(ac.A) / 255
	ba := float64(bc.A) / 255

	switch cm {
	case ComposeOver:
		fa, fb = 1, 1-aa
	case ComposeIn:
		fa, fb = ba, 0
	case ComposeOut:
		fa, fb = 1-ba, 0
	case ComposeAtop:
		fa, fb = ba, 1-aa
	case ComposeRover:
		fa, fb = 1-ba, 1
	case ComposeRin:
		fa, fb = 0, aa
	case ComposeRout:
		fa, fb = 0, 1-aa
	case ComposeRatop:
		fa, fb = 1-ba, aa
	case ComposeXor:
		fa, fb = 1-ba, 1-aa
	case ComposePlus:
		fa, fb = 1, 1
	case ComposeCopy:
		fa, fb = 1, 0
	default:
		panic(errors.New("Compose: invalid ComposeMethod"))
	}

	return ac.Mul(Alpha(fa)).Add(bc.Mul(Alpha(fb)))
}
