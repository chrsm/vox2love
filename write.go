package main

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/fogleman/gg"
)

func writeVox(v *vox, out string, extrude bool) {
	sz := v.sizes[0]

	w := int(sz.z * sz.x)
	h := int(sz.y)

	var palette []uint32
	if v.palette == nil {
		palette = defPalette
		log.Printf("\t\tusing def palette")
	} else {
		palette = v.palette.colors[:]
		log.Printf("\t\tusing custom palette")
	}

	for i := range v.voxels {
		// if we want to extrude, we need to pad EACH FRAME.
		// that means on the boundaries, for example:
		//
		// frame = 0; position is from x[0,32] and y[0,32].
		// we want this to be a 33x33.
		//
		// so each context is now:
		// w = w + n-frames
		// h = h + 1
		// as we write frames X->.

		vv := v.voxels[i]

		w, h := w, h
		if extrude {
			w, h = w+int(sz.z), h+1
		}

		dc := gg.NewContext(w, h)

		log.Printf("%d points to draw", len(vv.voxels))
		for j := range vv.voxels { // this is _not_ a frame
			vvv := vv.voxels[j]

			x := int(vvv.z)*int(sz.x) +
				int(vvv.x)
			y := int(vvv.y)

			// because we're cloning X at intervals,
			// we need to offset X by something.
			if extrude {
				x += (x / int(sz.z)) % int(sz.z)
			}

			idx := int(vvv.colorIndex)
			pAt := (*[4]byte)(unsafe.Pointer(&palette[idx]))[:]

			dc.SetRGBA255(
				int(pAt[0]),
				int(pAt[1]),
				int(pAt[2]),
				int(pAt[3]),
			)

			dc.SetPixel(x, y)
		}

		if extrude {
			img := dc.Image()

			// clone across interval'd X down to Y
			for x := 0; x < int(sz.x)*int(sz.z); x += int(sz.x) {
				// x is offsetted earlier
				truex := x + ((x / int(sz.x)) % int(sz.x))

				for y := 0; y < int(sz.y); y++ {
					cat := img.At(truex+7, y)
					dc.SetColor(cat)
					dc.SetPixel(truex+8, y)
				}
			}

			// clone across X,MAX(Y)
			for x := 0; x < int(sz.x)*int(sz.z)+int(sz.z); x++ {
				cat := img.At(x, int(sz.y)-1)
				dc.SetColor(cat)
				dc.SetPixel(x, int(sz.y))
			}
		}

		dc.SavePNG(fmt.Sprintf("%s_%d.png", out, i))
	}
}
