package text_test

import (
	"image/color"
	"testing"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/ext/text"
	"golang.org/x/image/font/inconsolata"
)

func TestAtlas7x13(t *testing.T) {
	if text.Atlas7x13 == nil {
		t.Fatalf("Atlas7x13 is nil")
	}

	for _, tt := range []struct {
		runes []rune
		want  bool
	}{{text.ASCII, true}, {[]rune("ÅÄÖ"), false}} {
		for _, r := range tt.runes {
			if got := text.Atlas7x13.Contains(r); got != tt.want {
				t.Fatalf("Atlas7x13.Contains('%s') = %v, want %v", string(r), got, tt.want)
			}
		}
	}
}

func TestAtlasInconsolata(t *testing.T) {
	text.NewAtlas(inconsolata.Regular8x16, text.ASCII)
}

func TestAtlasPictureDataCopy(t *testing.T) {
	a := text.NewAtlas(inconsolata.Regular8x16, text.ASCII)
	orig := a.Picture().(*pixel.PictureData)

	cp := a.PictureDataCopy()
	if cp == orig {
		t.Fatal("PictureDataCopy returned the same pointer as the original")
	}
	if cp.Stride != orig.Stride {
		t.Errorf("Stride mismatch: got %d, want %d", cp.Stride, orig.Stride)
	}
	if cp.Rect != orig.Rect {
		t.Errorf("Rect mismatch: got %v, want %v", cp.Rect, orig.Rect)
	}
	if len(cp.Pix) != len(orig.Pix) {
		t.Fatalf("Pix length mismatch: got %d, want %d", len(cp.Pix), len(orig.Pix))
	}
	// Verify deep copy: mutating the copy does not affect the original.
	if len(cp.Pix) > 0 {
		origPix := orig.Pix[0]
		cp.Pix[0] = color.RGBA{R: ^origPix.R, G: ^origPix.G, B: ^origPix.B, A: ^origPix.A}
		if orig.Pix[0] != origPix {
			t.Error("PictureDataCopy is not a deep copy: mutating the copy affected the original")
		}
	}
}

func TestAtlasCloneWithPictureData(t *testing.T) {
	a := text.NewAtlas(inconsolata.Regular8x16, text.ASCII)
	tests := []struct {
		name         string
		offset       pixel.Vec
		shouldPanic  bool
		panicMessage string
	}{
		{
			name:   "identity (frame == original bounds)",
			offset: pixel.ZV,
		},
		{
			name:   "translated into larger shared picture",
			offset: pixel.V(50, 30),
		},
		{
			name:         "wrong frame size panics",
			shouldPanic:  true,
			panicMessage: "atlas: new frame dimensions do not match prior picture",
		},
	}

	origBounds := a.Picture().Bounds()
	origW, origH := origBounds.W(), origBounds.H()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					r := recover()
					if r == nil {
						t.Fatal("expected panic, got none")
					}
					if msg, ok := r.(string); !ok || msg != tt.panicMessage {
						t.Errorf("wrong panic: %v", r)
					}
				}()
				// Supply a frame with wrong dimensions.
				wrongPic := pixel.MakePictureData(pixel.R(0, 0, origW+1, origH))
				wrongFrame := pixel.R(0, 0, origW+1, origH)
				a.CloneWithPictureData(wrongPic, wrongFrame)
				return
			}

			// Build a shared picture large enough to hold the atlas at the given offset.
			sharedW := origW + tt.offset.X
			sharedH := origH + tt.offset.Y
			sharedPic := pixel.MakePictureData(pixel.R(0, 0, sharedW, sharedH))
			frame := pixel.R(tt.offset.X, tt.offset.Y, tt.offset.X+origW, tt.offset.Y+origH)

			clone := a.CloneWithPictureData(sharedPic, frame)

			// The clone must contain every rune the original does.
			for _, r := range text.ASCII {
				if !clone.Contains(r) {
					t.Errorf("clone does not contain rune %q", r)
				}
			}

			// Glyph coordinates must be shifted by the frame offset relative to
			// the original atlas origin.
			delta := frame.Min.Sub(origBounds.Min)
			for _, r := range text.ASCII {
				if !a.Contains(r) {
					continue
				}
				origGlyph := a.Glyph(r)
				cloneGlyph := clone.Glyph(r)
				wantDot := origGlyph.Dot.Add(delta)
				if cloneGlyph.Dot != wantDot {
					t.Errorf("rune %q: Dot got %v, want %v", r, cloneGlyph.Dot, wantDot)
				}
				wantFrame := origGlyph.Frame.Moved(delta)
				if cloneGlyph.Frame != wantFrame {
					t.Errorf("rune %q: Frame got %v, want %v", r, cloneGlyph.Frame, wantFrame)
				}
				if cloneGlyph.Advance != origGlyph.Advance {
					t.Errorf("rune %q: Advance got %v, want %v", r, cloneGlyph.Advance, origGlyph.Advance)
				}
			}
		})
	}
}
