package atlas

import (
	"embed"
	"fmt"
	"image"

	"github.com/gopxl/pixel/v2"
	"golang.org/x/exp/maps"
)

type Group struct {
	atlas    *Atlas
	textures []TextureId
	slices   []SliceId
}

// MakeGroup creates a new group of textures.
func (a *Atlas) MakeGroup() Group {
	return Group{
		atlas: a,
	}
}

// DefaultGroup returns the default group of the atlas.
func (a *Atlas) DefaultGroup() *Group {
	if a.defaultGroup.atlas == nil {
		a.defaultGroup.atlas = a
	}
	return &a.defaultGroup
}

// Clear removes the given texture groups from the atlas.
// If no groups are given, all textures are removed.
func (a *Atlas) Clear(groups ...Group) {
	if len(groups) == 0 {
		maps.Clear(a.idMap)
	}

	for _, group := range groups {
		for _, texture := range group.textures {
			delete(a.idMap, texture.id)
		}
		for _, slice := range group.slices {
			for i := uint32(0); i < slice.len; i++ {
				delete(a.idMap, slice.start.id+i)
			}
		}
	}

	a.clean = false

	a.Pack()
}

func (g *Group) addEntry(entry iEntry) (id TextureId) {
	if bw, bh := entry.Bounds().Dx(), entry.Bounds().Dy(); bw > MaxTextureSize || bh > MaxTextureSize {
		panic(fmt.Errorf("Texture is larger (%v, %v) than the maximum allowed texture (%v, %v)", bw, bh, MaxTextureSize, MaxTextureSize))
	}

	id = TextureId{id: g.atlas.id, atlas: g.atlas}
	g.textures = append(g.textures, id)
	switch entry := entry.(type) {
	case iSliceEntry:
		g.atlas.id += uint32((entry.Bounds().Dx() / entry.Frame().X) * (entry.Bounds().Dy() / entry.Frame().Y))
	default:
		g.atlas.id++
	}
	g.atlas.adding = append(g.atlas.adding, entry)
	g.atlas.clean = false
	return
}

// AddImage loads an image to the atlas.
func (g *Group) AddImage(img image.Image) (id TextureId) {
	e := imageEntry{
		entry: entry{
			id:     g.atlas.id,
			bounds: img.Bounds(),
		},
		data: img,
	}
	return g.addEntry(e)
}

// AddEmbed loads an embed.FS image to the atlas.
func (g *Group) AddEmbed(fs embed.FS, path string) (id TextureId) {
	img, err := loadEmbedSprite(fs, path)
	if err != nil {
		panic(err)
	}
	e := embedEntry{
		fileEntry: fileEntry{
			entry: entry{
				id:     g.atlas.id,
				bounds: img.Bounds(),
			},
			path: path,
		},
		fs: fs,
	}
	return g.addEntry(e)
}

// AddFile loads an image file to the atlas.
func (g *Group) AddFile(path string) (id TextureId) {
	img, err := loadSprite(path)
	if err != nil {
		panic(err)
	}
	e := fileEntry{
		entry: entry{
			id:     g.atlas.id,
			bounds: img.Bounds(),
		},
		path: path,
	}
	return g.addEntry(e)
}

// SliceImage evenly divides the given image into cells of the given size.
func (g *Group) SliceImage(img image.Image, cellSize pixel.Vec) (id SliceId) {
	frame := image.Pt(int(cellSize.X), int(cellSize.Y))
	bounds := img.Bounds()
	if bounds.Dx()%frame.X != 0 || bounds.Dy()%frame.Y != 0 {
		panic(fmt.Sprintf("Texture size (%v,%v) must be multiple of cellSize (%v,%v)", bounds.Dx(), bounds.Dy(), cellSize.X, cellSize.Y))
	}

	e := sliceImageEntry{
		imageEntry: imageEntry{
			entry: entry{
				id:     g.atlas.id,
				bounds: bounds,
			},
			data: img,
		},
		sliceEntry: sliceEntry{
			frame: frame,
		},
	}
	return SliceId{
		start: g.addEntry(e),
		len:   uint32((bounds.Dx() / frame.X) * (bounds.Dy() / frame.Y)),
	}
}

// SliceFile loads an image and evenly divides it into cells of the given size.
func (g *Group) SliceFile(path string, cellSize pixel.Vec) (id SliceId) {
	frame := image.Pt(int(cellSize.X), int(cellSize.Y))
	img, err := loadSprite(path)
	if err != nil {
		panic(err)
	}
	bounds := img.Bounds()
	if bounds.Dx()%frame.X != 0 || bounds.Dy()%frame.Y != 0 {
		panic(fmt.Sprintf("Texture size (%v,%v) must be multiple of cellSize (%v,%v)", bounds.Dx(), bounds.Dy(), cellSize.X, cellSize.Y))
	}

	e := sliceFileEntry{
		fileEntry: fileEntry{
			entry: entry{
				id:     g.atlas.id,
				bounds: bounds,
			},
			path: path,
		},
		sliceEntry: sliceEntry{
			frame: frame,
		},
	}

	return SliceId{
		start: g.addEntry(e),
		len:   uint32((bounds.Dx() / frame.X) * (bounds.Dy() / frame.Y)),
	}
}

// SliceEmbed loads an embeded image and evenly divides it into cells of the given size.
func (g *Group) SliceEmbed(fs embed.FS, path string, cellSize pixel.Vec) (id SliceId) {
	img, err := loadEmbedSprite(fs, path)
	if err != nil {
		panic(err)
	}
	frame := image.Pt(int(cellSize.X), int(cellSize.Y))
	bounds := img.Bounds()
	if bounds.Dx()%frame.X != 0 || bounds.Dy()%frame.Y != 0 {
		panic(fmt.Sprintf("Texture size (%v,%v) must be multiple of cellSize (%v,%v)", bounds.Dx(), bounds.Dy(), cellSize.X, cellSize.Y))
	}

	e := sliceEmbedEntry{
		embedEntry: embedEntry{
			fileEntry: fileEntry{
				entry: entry{
					id:     g.atlas.id,
					bounds: bounds,
				},
				path: path,
			},
			fs: fs,
		},
		sliceEntry: sliceEntry{
			frame: frame,
		},
	}

	return SliceId{
		start: g.addEntry(e),
		len:   uint32((bounds.Dx() / frame.X) * (bounds.Dy() / frame.Y)),
	}
}
