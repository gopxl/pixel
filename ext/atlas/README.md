# The Atlas Extension

This extension provides an implementation of a texture atlas. A texture atlas groups individual textures together into a single texture that is sent to the GPU. This lessens the time spent switching GPU textures which can be expensive.

## Usage

You create an atlas by just declaring a variable:

```go
import "github.com/gopxl/pixel/v2/ext/atlas"

var textures atlas.Atlas
```

No initialization needed! You're ready to start adding textures into the atlas.

### Adding Textures

There are a few supported ways to add textures to the atlas:
1. Image Data
2. Loading from a file
3. Loading from an embedded file

Each of these methods returns an `atlas.TextureId`, you'll need to keep track of this value as it's how you draw that individual texture.
The `atlas` package doesn't not enforce how you store them. So you could store them statically in variables:

```go
var (
   textures atlas.Atlas
   texture1 = textures.AddFile("1.png")
   texture2 = textures.AddFile("2.png")
)
```

Or you could dump them into a map with names to lookup:

```go
func run() {
   var textures atlas.Atlas

   textureMap := make(map[string]atlas.TextureId)

   textureMap["1"] = textures.AddFile("1.png")
   textureMap["2"] = textures.AddFile("2.png")
}
```

Or whatever other way you can come up with!

##### Adding Image Data

If you've already loaded an image and want to copy it into the atlas you can use this method.

```go
func run() {
   var textures atlas.Atlas

   f, err := os.Open("1.png")
   defer f.Close()

   i, err := png.Decode(f)

   texture1 := textures.AddImage(i)
}
```

##### Adding an Image File

The atlas also has convience methods to load a file directly from the path.

```go
func run() {
   var textures atlas.Atlas

   texture1 := textures.AddFile("1.png")
}
```

##### Adding an Embedded Image

The atlas also supports Go's embedded file system.

```go
// go:embed 1.png
var embedded embed.FS

func run() {
   var textures atlas.Atlas

   texture1 := textures.AddEmbed(embedded, "1.png")
}
```

#### Sliced Textures

Sometimes, you already have a texture that contains multiple sprites of equal size in it; atlas can load these directly, cutting the single texture into multiple textures to easily use.

Instead of an `atlas.TextureId`, creating a sliced texture returns an `atlas.SliceId`. This allows you to access the frames of the sliced texture; we'll show this in more detail later.

Each of the `atlas.Slice` methods take an additional `pixel.Vec`. This is the sub-image size and it tells the atlas how to slice up the texture.

Once you have a `atlas.SliceId`, you can use that to get all of the textures that were added to the atlas

```go
var textures atlas.Atlas

sliced1 := textures.SliceFile("sheet.png", pixel.V(8, 8))
walk0 := sliced1.Frame(0)
walk1 := sliced1.Frame(1)
```

**Note:** If you attempt to index a non-existant frame, a panic will be raised.

##### Slicing Image Data

```go
func run() {
   var textures atlas.Atlas

   f, err := os.Open("sheet1.png")
   defer f.Close()

   i, err := png.Decode(f)

   sliced1 := textures.SliceImage(f, pixel.V(8, 8))
}
```

##### Slicing Image File

```go
func run() {
   var textures atlas.Atlas

   sliced1 := textures.SliceFile("sheet1.png", pixel.V(8, 8))
}
```

##### Slicing an Embedded Image

```go
// go:embed sheet1.png
var embedded embed.FS

func run() {
   var textures atlas.Atlas

   sheet1 := textures.SliceEmbed(embedded, "sheet1.png", pixel.V(8, 8))
}
```

### Packing the Atlas

Once you've added all of the textures to the atlas you wish, it needs to be packed.

This sorts the added textures by size to minimize the atlas texture. It then creates as many `pixel.PictureData` textures as it needs to store everything, then copies all of the added images to those textures.

```go
var textures atlas.Atlas

// ... use the Add* and/or Slice* methods to add textures

textures.Pack()
```

### Drawing Atlas Textures

#### Drawing TextureId

When you have an `atlas.TextureId`, you can draw it to a `pixel.Target` normally.

```go
var textures atlas.Atlas

sprite1 := textures.AddFile("1.png")

textures.Pack()

sprite1.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
```

As you can see, it draws just like the `pixel.Sprite` does, but much more efficiently because this limits the amount of texture changing on the GPU.

#### Drawing SliceId

Drawing an `atlas.SliceId` is slightly different because you need to specify which frame of the sliced texture to draw.

```go
var textures atlas.Atlas

sliced1 := textures.SliceFile("1.png", pixel.V(8, 8))

textures.Pack()

sliced1.Draw(win, pixel.IM.Moved(win.Bounds().Center()), 1)
```

This will draw the second (0-indexed) image in the sliced texture. The index is calculated left to right, top to bottom.

--
OR
--

You could expand out the `atlas.SliceId` into the `atlas.TextureId` that make it up.

```go
var textures atlas.Atlas

sliced1 := textures.SliceFile("sheet.png")
walk0 := sliced1.Frame(0)
walk1 := sliced1.Frame(1)

walk0.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
```

### Groups

Groups are a construct that allow logical grouping of textures for a couple of reasons:

1. You want to have some static textures that are always in the atlas.
2. You want to be able to remove level-specific textures without having to re-add the other textures.
3. You want to use another library that wants to adds its own textures to the atlas without having to worry about its textures being removed out from under it.

You can create a group:

```go
var textures atlas.Atlas

group1 := textures.MakeGroup()
```

Or you could use the default group that comes with the atlas (if you've been following along at home, you've been unknowingly been using this: `atlas.Atlas.Add*` and `atlas.Atlas.Slice*` use the default group).

```go
var textures atlas.Atlas

group1 := textures.DefaultGroup()
```

Groups share the same `Add*` and `Slice*` methods as are on the `atlas.Atlas`.

### Clearing Textures

You can remove all of the textures in an atlas with:

```go
var textures atlas.Atlas

// Add some textures to the atlas

// ...

// Actually, we don't want them anymore
textures.Clear()
```

**Note:** You don't need to call `atlas.Atlas.Pack()` after clearing textures, `atlas.Atlas.Clear()` does this automatically.

#### Clearing Groups

The main feature of groups is being able to remove them from the atlas.

```go
var textures atlas.Atlas

keepThisAroundGroup := textures.MakeGroup()
imGoingToDeleteThisSoon := textures.MakeGroup()

// Add some textures to the groups

// ...

// Actually, we don't want them anymore
textures.Clear(imGoingToDeleteThisSoon)

// Atlas still has all textures added to `keepThisAroundGroup`
```

