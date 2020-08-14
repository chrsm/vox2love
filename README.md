## vox2love

`vox2love` converts a MagicaVoxel `.vox` file into images that can be used
in LOVE.

Note that this currently only supports v150 of the `.vox` format, as specified
in `docs/file-format.txt`, which was pulled from [here][1].

### Generated Images

`vox2love` will generate an image that is `W = X * FRAMES, H = Y`, meaning it treats
each 'layer' of a voxel as a single `WxH` image and creates the next layer
as an extension (across the X axis).

For a model that is 32x32x8, it will create an image that is 256 pixels wide
and 32 pixels tall.

### Using inside of LOVE

To use this inside of love, load the texture (`love.graphics.newImage`),
pass it to a `SpriteBatch` and `:add` sprites to it that reference the texture
using quads that match each frame.

IE, for a 32x32x8, you could do something like:

```lua
local sb = love.graphics.newSpriteBatch(texture, 8)
local quads = {}
for i=1,8 do
    local q = love.graphics.newQuad((j - 1) * 32, 0, 32, 32, 256, 32)
    quads[j] = q

    sb:add(q)
end
```

To make the most out of using a voxel, rather than just any regular image,
you want to make sure that you update the sprites with `:set` when rotation
changes to create a psuedo-3d effect - otherwise you're better off using
a single image.

I'm experimenting with this and using `newArrayImage` to create pseudo-3d tile
maps, and hopefully will be able to share something useful in the future. I'm
new to LOVE and graphics in general, though, so let me know if you have thoughts.

### Running

`vox2love` requires two flags:

- `-p`: path to either a single `.vox` file or directory of `.vox` files.
  If multiple `.vox` are found, it will generate a png for each one and each
  animation within it (multiple models, basically).

- `-o`: path to where the `.png` files will be written

There is an optional argument, `-e`, which will extrude generated images.
Each section of the generated image will be padded by 1px, so a 32x32x8 becomes:

    w = 32 * 8 + 8
    h = 33 + 1

The padded pixels will use their neighbor's color, IE:

- `(x: 32, y: 0)` will copy the color from `(x: 31, y: 0)`
- repeat for each section

- `(x: 0, y: 32)` will copy the color from `(x: 0, y: 31)`
- repeat for every X along Y. 

This is designed to sidestep some OpenGL/Love specific (?) issues with black
bars being rendered at certain scales. I am not 100% sure that my
implementation is correct, so please feel free to @ me.

### Building

To build `vox2love` from source, check out this repo, run `go mod download` and
`go build`. This should work on all platforms with no issues; I'll (eventually)
upload some binaries to the release section so that this isn't necessary.


[1]: https://github.com/ephtracy/voxel-model
