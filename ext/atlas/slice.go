package atlas

import "github.com/gopxl/pixel/v2"

// A SliceId represents a texture in the atlas added by Atlas.Slice.
// This differs from a TextureId in that it's meant to be drawn with a frame offset (a sub image).
type SliceId struct {
	start TextureId
	len   uint32
}

// Frame returns a TextureId representing the given frame of the slice
func (s SliceId) Frame(frame uint32) TextureId {
	if frame >= s.len {
		panic("slice frame out of bounds")
	}
	return TextureId{id: s.start.id + frame, atlas: s.start.atlas}
}

// Bounds returns the bounds of the slice in the atlas.
func (s SliceId) Bounds(frame uint32) pixel.Rect {
	return s.Frame(frame).Bounds()
}

// Draw draws the slice in the atlas to the target with the given matrix.
func (s SliceId) Draw(t pixel.Target, m pixel.Matrix, frame uint32) {
	f := s.Frame(frame)
	f.Draw(t, m)
}
