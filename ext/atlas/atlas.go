package atlas

import (
	"embed"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path"
	"sort"

	"github.com/gopxl/pixel/v2"
	"github.com/pkg/errors"
)

const (
	MaxTextureSize = 8192
)

type loc struct {
	index int
	rect  image.Rectangle
}

type spaces []image.Rectangle

type sheet struct {
	size   image.Rectangle
	spaces spaces
}

type Atlas struct {
	adding       []iEntry
	internal     []*pixel.PictureData
	clean        bool
	idMap        map[uint32]loc
	id           uint32
	defaultGroup Group
}

// Dump writes out the internal textures to disk as PNG files.
func (a *Atlas) Dump(dir string) {
	if !a.clean {
		panic("Atlas is dirty, call atlas.Pack() first")
	}

	// TODO improve dump so that the resulting file(s) can be directly loaded into the atlas from disk.
	//		Could be something like a zip file with the images and a file with the locations of the individual frames.
	for i, t := range a.internal {
		f, err := os.Create(path.Join(dir, fmt.Sprintf("%v.png", i)))
		if err != nil {
			log.Println(i, err)
			continue
		}
		defer f.Close()

		if err := png.Encode(f, t.Image()); err != nil {
			log.Println(i, err)
			continue
		}
	}
}

// Textures returns a copy of all of the internal packed textures.
func (a *Atlas) Textures() []*pixel.PictureData {
	if !a.clean {
		panic("Atlas is dirty, call atlas.Pack() first")
	}

	data := make([]*pixel.PictureData, len(a.internal))

	for i := range a.internal {
		data[i] = pixel.PictureDataFromPicture(a.internal[i])
	}

	return data
}

// Images returns a copy of all of the internal packed textures as image.Image.
func (a *Atlas) Images() []image.Image {
	if !a.clean {
		panic("Atlas is dirty, call atlas.Pack() first")
	}

	images := make([]image.Image, len(a.internal))

	for i := range a.internal {
		images[i] = a.internal[i].Image()
	}

	return images
}

// AddImage loads an image to the atlas.
func (a *Atlas) AddImage(img image.Image) (id TextureId) {
	return a.DefaultGroup().AddImage(img)
}

// AddEmbed loads an embed.FS image to the atlas.
func (a *Atlas) AddEmbed(fs embed.FS, path string) (id TextureId) {
	return a.DefaultGroup().AddEmbed(fs, path)
}

// AddFile loads an image file to the atlas.
func (a *Atlas) AddFile(path string) (id TextureId) {
	return a.DefaultGroup().AddFile(path)
}

// SliceImage evenly divides the given image into cells of the given size.
func (a *Atlas) SliceImage(img image.Image, cellSize pixel.Vec) (id SliceId) {
	return a.DefaultGroup().SliceImage(img, cellSize)
}

// Slice loads an image and evenly divides it into cells of the given size.
func (a *Atlas) SliceFile(path string, cellSize pixel.Vec) (id SliceId) {
	return a.DefaultGroup().SliceFile(path, cellSize)
}

// SliceEmbed loads an embeded image and evenly divides it into cells of the given size.
func (a *Atlas) SliceEmbed(fs embed.FS, path string, cellSize pixel.Vec) (id SliceId) {
	return a.DefaultGroup().SliceEmbed(fs, path, cellSize)
}

// Pack takes all of the added textures and adds them to the atlas largest to smallest,
// trying to waste as little space as possible. After this call, the textures added
// to the atlas can be used.
func (a *Atlas) Pack() {
	// If there's nothing to do, don't do anything
	if a.clean || len(a.adding) == 0 {
		return
	}

	// If we've already packed the textures, we need to convert them back to images to repack them
	if a.internal != nil && len(a.internal) > 0 {
		images := make([]*image.RGBA, len(a.internal))
		for i, data := range a.internal {
			images[i] = data.Image()
		}

		for id, loc := range a.idMap {
			bounds := image.Rect(0, 0, loc.rect.Dx(), loc.rect.Dy())
			rgba := image.NewRGBA(bounds)
			i := images[loc.index]

			for y := 0; y < bounds.Dy(); y++ {
				for x := 0; x < bounds.Dx(); x++ {
					rgba.Set(x, y, i.At(loc.rect.Min.X+x, loc.rect.Min.Y+y))
				}
			}

			entry := imageEntry{data: rgba}
			entry.id = id
			entry.bounds = bounds

			a.adding = append(a.adding, entry)
		}
	}

	// reset internal stuff
	a.internal = a.internal[:0]
	if a.idMap == nil {
		a.idMap = make(map[uint32]loc)
	} else {
		clear(a.idMap)
	}

	sort.Slice(a.adding, func(i, j int) bool {
		return area(a.adding[i].Bounds()) >= area(a.adding[j].Bounds())
	})

	sheets := make([]sheet, 1)
	for i := range sheets {
		sheets[i] = sheet{
			spaces: []image.Rectangle{image.Rect(0, 0, MaxTextureSize, MaxTextureSize)},
		}
	}

	for _, add := range a.adding {
		bw, bh := add.Bounds().Dx(), add.Bounds().Dy()

		found := image.Rectangle{}
		foundI := -1

	Loop:
		for i := range sheets {
			for j := range sheets[i].spaces {
				found, sheets[i].spaces = split(sheets[i].spaces, j, bw, bh)
				if found.Empty() {
					continue
				}
				sort.Slice(sheets[i].spaces, func(a, b int) bool {
					return area(sheets[i].spaces[a]) < area(sheets[i].spaces[b])
				})
				foundI = i
				break Loop
			}
		}

		if foundI == -1 {
			foundI = len(sheets)
			sheets = append(sheets, sheet{})
			found, sheets[foundI].spaces = split([]image.Rectangle{image.Rect(0, 0, MaxTextureSize, MaxTextureSize)}, 0, bw, bh)
		}

		// Increase the size of the Atlas so we can allocate the minimum-sized
		// 	texture later.
		if found.Min.X == 0 {
			sheets[foundI].size.Max.Y += found.Dy()
		}
		if found.Min.Y == 0 {
			sheets[foundI].size.Max.X += found.Dx()
		}

		switch add := add.(type) {
		case iSliceEntry:
			// If we have a frame, that means we just added a sprite sheet to the sprite sheet
			// 	so we need to add id entries for each of the sprites
			id := add.Id()
			for y := 0; y < add.Bounds().Dy(); y += add.Frame().Y {
				for x := 0; x < add.Bounds().Dx(); x += add.Frame().X {
					a.idMap[id] = loc{
						index: foundI,
						rect:  rect(found.Min.X+x, found.Min.Y+y, add.Frame().X, add.Frame().Y),
					}
					id++
				}
			}
		default:
			// Found a spot, add it to the map
			a.idMap[add.Id()] = loc{
				index: foundI,
				rect:  found,
			}
		}
	}

	// Create internal textures
	sprites := make([]*image.RGBA, len(sheets))
	for i := range sheets {
		if !sheets[i].size.Empty() {
			sprites[i] = image.NewRGBA(sheets[i].size)
		}
	}

	// Copy individual sprite data into internal textures
	for _, add := range a.adding {
		var (
			err    error
			sprite image.Image
			s      = a.idMap[add.Id()]
		)

		switch add := add.(type) {
		case iImageEntry:
			sprite = add.Data()
		case iEmbedEntry:
			sprite, err = loadEmbedSprite(add.FS(), add.Path())
			err = errors.Wrapf(err, "failed to load embed sprite: %v", add.Path())
		case iFileEntry:
			sprite, err = loadSprite(add.Path())
			err = errors.Wrapf(err, "failed to load sprite file: %v", add.Path())
		}
		if err != nil {
			panic(err)
		}
		draw.Draw(sprites[s.index], rect(s.rect.Min.X, s.rect.Min.Y, add.Bounds().Dx(), add.Bounds().Dy()), sprite, image.Point{}, draw.Src)
	}

	// Make the internal Textures
	a.internal = make([]*pixel.PictureData, len(sprites))
	for i, sprite := range sprites {
		data := pixel.PictureDataFromImage(sprite)
		a.internal[i] = data
	}

	a.adding = nil
	a.clean = true

	return
}
