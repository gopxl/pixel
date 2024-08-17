package atlas

import (
	"fmt"

	"github.com/gopxl/pixel/v2"
)

// TextureId is a reference to a texture in an atlas.
type TextureId struct {
	id     uint32
	atlas  *Atlas
	sprite *pixel.Sprite
}

// Bounds returns the bounds of the texture in the atlas.
func (t TextureId) Bounds() pixel.Rect {
	s, has := t.atlas.idMap[t.id]
	if !has {
		panic(fmt.Sprintf("id: %v does not exit in atlas", t.id))
	}
	return pixelRect(0, 0, s.rect.Dx(), s.rect.Dy())
}

// Draw draws the texture in the atlas to the target with the given matrix.
func (t *TextureId) Draw(target pixel.Target, m pixel.Matrix) {
	if !t.atlas.clean {
		panic("Packer is dirty, call atlas.Pack() first")
	}

	l, has := t.atlas.idMap[t.id]
	if !has {
		panic(fmt.Sprintf("id [%v] does not exist in packer", t.id))
	}

	if t.sprite == nil {
		r := image2PixelRect(l.rect)
		c := t.atlas.internal[l.index].Bounds().Center()
		m := pixel.IM.ScaledXY(c, pixel.V(1, -1))
		r.Min = m.Project(r.Min)
		r.Max = m.Project(r.Max)
		t.sprite = pixel.NewSprite(t.atlas.internal[l.index], r)
	}
	t.sprite.Draw(target, m)
}
