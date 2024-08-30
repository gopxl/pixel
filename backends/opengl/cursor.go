package opengl

import (
	"image"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/gopxl/mainthread/v2"
	"github.com/gopxl/pixel/v2"
)

type StandardCursor = glfw.StandardCursor

const (
	ArrowCursor     = glfw.ArrowCursor
	IBeamCursor     = glfw.IBeamCursor
	CrosshairCursor = glfw.CrosshairCursor
	HandCursor      = glfw.HandCursor
	HResizeCursor   = glfw.HResizeCursor
	VResizeCursor   = glfw.VResizeCursor
)

type Cursor = glfw.Cursor

// CreateStandardCursor creates a new standard cursor.
func CreateStandardCursor(cursorId StandardCursor) *Cursor {
	c := mainthread.CallVal(func() *Cursor {
		return glfw.CreateStandardCursor(cursorId)
	})
	runtime.SetFinalizer(c, (*Cursor).Destroy)
	return c
}

// CreateCursorImage creates a new cursor from an image with the specified hotspot (where the click is registered).
func CreateCursorImage(img image.Image, hot pixel.Vec) *Cursor {
	c := mainthread.CallVal(func() *Cursor {
		return glfw.CreateCursor(img, int(hot.X), int(hot.Y))
	})
	runtime.SetFinalizer(c, (*Cursor).Destroy)
	return c
}

// SetCursor sets the cursor for the window.
func (w *Window) SetCursor(cursor *Cursor) {
	mainthread.Call(func() {
		w.window.SetCursor(cursor)
		w.cursor = cursor
	})
}
