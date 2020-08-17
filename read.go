package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime/debug"
)

var (
	errMissingMagic = errors.New("magic VOX string not found")
	errBadVersion   = errors.New("only version 150 supported")
	errMissingMain  = errors.New("missing MAIN chunk")
)

func newVox(r io.ReadSeeker) (*vox, error) {
	br := breader{r, nil}
	buf := make([]byte, 4)

	br.Read(buf)
	if !bytes.Equal(bVOX, buf) {
		return nil, errMissingMagic
	}

	br.Read(buf)
	if !bytes.Equal(bVER, buf) {
		return nil, errBadVersion
	}

	br.Read(buf)
	if !bytes.Equal(bMAIN, buf) {
		return nil, errMissingMain
	}

	// there's an implicit `chunk` quality to main,
	// but it is the 'container' of other chunks,
	// so it isn't represented here inside of a *vox
	mc := chunk{}
	br.ReadV(&mc.sz)
	br.ReadV(&mc.szChildren)

	v := &vox{}
	for br.err == nil {
		c := chunk{}

		br.Read(c.id[:])
		if br.err != nil {
			break
		}

		br.ReadV(&c.sz)
		br.ReadV(&c.szChildren)
		copy(buf, c.id[:])

		log.Printf("\t\tchunk(%s):\tsz=%d\tszc=%d", c.id[:], c.sz, c.szChildren)
		switch {
		case bytes.Equal(bPACK, buf):
			p := &pack{}

			br.ReadV(&p.numModels)

			v.pack = p

		case bytes.Equal(bSIZE, buf):
			sz := size{}

			br.ReadV(&sz.x)
			br.ReadV(&sz.y)
			br.ReadV(&sz.z)

			v.sizes = append(v.sizes, sz)
		case bytes.Equal(bXYZI, buf):
			vx := xyzi{}

			br.ReadV(&vx.szVoxels)

			vx.voxels = make([]voxel, int(vx.szVoxels))

			for i := 0; i < int(vx.szVoxels); i++ {
				br.Read(buf)

				vx.voxels[i].x = buf[0]
				vx.voxels[i].y = buf[1]
				vx.voxels[i].z = buf[2]
				vx.voxels[i].colorIndex = buf[3]
			}

			v.voxels = append(v.voxels, vx)
		case bytes.Equal(bRGBA, buf):
			pal := &rgba{}

			for i := 0; i < 254; i++ {
				br.ReadV(&pal.colors[i+1])
			}

			v.palette = pal

			// skip 4
			br.Seek(16, io.SeekCurrent)

		case bytes.Equal(bMATT, buf):
			// for now, we're skipping this entire section.
			/*
				m := material{}

				br.ReadV(&m.id)
				br.ReadV(&m.typ)
				br.ReadV(&m.weight)
				br.ReadV(&m.properties)
			*/
			br.Seek(int64(c.sz+c.szChildren), io.SeekCurrent)
		}
	}

	if br.err == io.EOF {
		return v, nil
	}

	return v, br.err
}

type breader struct {
	r io.ReadSeeker

	err error
}

func (b *breader) ReadV(v interface{}) {
	if b.err != nil {
		return
	}

	err := binary.Read(b.r, binary.LittleEndian, v)
	b.setErr(err)
}

func (b *breader) Read(dst []byte) {
	if b.err != nil {
		return
	}

	_, err := b.r.Read(dst)
	b.setErr(err)
}

func (b *breader) Seek(pos int64, dir int) {
	if b.err != nil {
		return
	}

	_, err := b.r.Seek(pos, dir)
	b.setErr(err)
}

func (b *breader) setErr(err error) {
	if err != nil {
		if err == io.EOF {
			b.err = err
			return
		}

		debug.PrintStack()
		cur, _ := b.r.Seek(0, io.SeekCurrent)
		b.err = fmt.Errorf("%s (pos: %d)", err, cur)
	}
}
