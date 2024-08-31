package atlas

import (
	"embed"
	"image"

	"github.com/gopxl/pixel/v2"
)

type iEntry interface {
	Bounds() image.Rectangle
	Id() uint32
}

type entry struct {
	id     uint32
	bounds image.Rectangle
}

func (e entry) Id() uint32 {
	return e.id
}

func (e entry) Bounds() image.Rectangle {
	return e.bounds
}

type iEmbedEntry interface {
	iFileEntry
	FS() embed.FS
}

type embedEntry struct {
	fileEntry
	fs embed.FS
}

func (e embedEntry) FS() embed.FS {
	return e.fs
}

type iImageEntry interface {
	iEntry
	Data() image.Image
}

type imageEntry struct {
	entry
	data image.Image
}

func (i imageEntry) Data() image.Image {
	return i.data
}

type iFileEntry interface {
	iEntry
	Path() string
	DecoderFunc() pixel.DecoderFunc
}

type fileEntry struct {
	entry
	path        string
	decoderFunc pixel.DecoderFunc
}

func (f fileEntry) Path() string {
	return f.path
}

func (f fileEntry) DecoderFunc() pixel.DecoderFunc {
	return f.decoderFunc
}

type iSliceEntry interface {
	iEntry
	Frame() image.Point
}

type sliceEntry struct {
	frame image.Point
}

func (s sliceEntry) Frame() image.Point {
	return s.frame
}

type sliceImageEntry struct {
	sliceEntry
	imageEntry
}

type sliceFileEntry struct {
	sliceEntry
	fileEntry
}

type sliceEmbedEntry struct {
	sliceEntry
	embedEntry
}
