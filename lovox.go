package main

import (
	"fmt"
	"log"
	"os"
	"unsafe"

	"github.com/fogleman/gg"
)

func writeLovox(v *vox, out string) {
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
		dc := gg.NewContext(w, h)
		vv := v.voxels[i]

		log.Printf("%d points to draw", len(vv.voxels))
		for j := range vv.voxels {
			vvv := vv.voxels[j]

			x := float64(
				int(vvv.z)*int(sz.x)+
					int(vvv.x)) + 0.5
			y := float64(vvv.y) + 0.5

			idx := int(vvv.colorIndex)
			pAt := (*[4]byte)(unsafe.Pointer(&palette[idx]))[:]

			dc.SetRGBA255(
				int(pAt[0]),
				int(pAt[1]),
				int(pAt[2]),
				int(pAt[3]),
			)

			dc.SetPixel(int(x), int(y))
		}

		dc.SavePNG(fmt.Sprintf("%s_%d.png", out, i))

		// save flags.lua
		fp, err := os.OpenFile(fmt.Sprintf("%s_%d.lua", out, i), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
		if err != nil {
			panic(err)
		}

		fp.WriteString(
			fmt.Sprintf(
				flagsTpl,
				w/int(sz.x),
				sz.x,
				sz.y,
			),
		)

		fp.Close()
	}
}

const flagsTpl = `return {
    frames = %d,
    width  = %d,
    height = %d,
}`
