//go:build js && wasm

package opengl

import (
	"image"

	"github.com/gopxl/pixel/v2"
)

// StandardCursor matches the desktop-side typed constant. Under WASM we have
// no real cursor API wired up, so the constants exist only so game code can
// reference them unchanged.
type StandardCursor int

const (
	ArrowCursor StandardCursor = iota
	IBeamCursor
	CrosshairCursor
	HandCursor
	HResizeCursor
	VResizeCursor
)

// Cursor is an opaque handle. We keep the struct empty; the browser manages
// the mouse cursor itself.
type Cursor struct{}

func CreateStandardCursor(StandardCursor) *Cursor       { return &Cursor{} }
func CreateCursorImage(image.Image, pixel.Vec) *Cursor  { return &Cursor{} }
func (w *Window) SetCursor(*Cursor)                     {}
