package atlas

import (
	"fmt"

	"github.com/gopxl/pixel/v2"
)

// Get returns a texture with the given ID.
func (a *Atlas) Get(id uint32) TextureId {
	return TextureId{
		id:    id,
		atlas: a,
	}
}

// TextureId is a reference to a texture in an atlas.
type TextureId struct {
	id     uint32
	atlas  *Atlas
	sprite *pixel.Sprite
}

// ID returns the ID of the texture in the atlas.
func (t TextureId) ID() uint32 {
	return t.id
}

// Frame returns the frame of the texture in the atlas.
func (t TextureId) Frame() pixel.Rect {
	if !t.atlas.clean {
		panic("Atlas is dirty, call atlas.Pack() first")
	}
	s, has := t.atlas.idMap[t.id]
	if !has {
		panic(fmt.Sprintf("id: %v does not exist in atlas", t.id))
	}
	r := image2PixelRect(s.rect)
	c := t.atlas.internal[s.index].Bounds().Center()
	m := pixel.IM.ScaledXY(c, pixel.V(1, -1))
	r.Min = m.Project(r.Min)
	r.Max = m.Project(r.Max)
	return r
}

// Bounds returns the bounds of the texture in the atlas.
func (t TextureId) Bounds() pixel.Rect {
	if !t.atlas.clean {
		panic("Atlas is dirty, call atlas.Pack() first")
	}

	s, has := t.atlas.idMap[t.id]
	if !has {
		panic(fmt.Sprintf("id: %v does not exist in atlas", t.id))
	}
	return pixelRect(0, 0, s.rect.Dx(), s.rect.Dy())
}

// Draw draws the texture in the atlas to the target with the given matrix.
func (t *TextureId) Draw(target pixel.Target, m pixel.Matrix) {
	if !t.atlas.clean {
		panic("Atlas is dirty, call atlas.Pack() first")
	}

	l, has := t.atlas.idMap[t.id]
	if !has {
		panic(fmt.Sprintf("id [%v] does not exist in packer", t.id))
	}

	if t.sprite == nil {
		frame := t.Frame()
		t.sprite = pixel.NewSprite(t.atlas.internal[l.index], frame)
	}
	t.sprite.Draw(target, m)
}
