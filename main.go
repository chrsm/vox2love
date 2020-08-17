package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	voxDir  string // path/to/src
	outDir  string // path/to/dest
	extrude bool   // whether or not to pad images
)

func main() {
	log.Println("vox2love")

	flag.StringVar(&voxDir, "p", "", "path to directory containin vox files")
	flag.StringVar(&outDir, "o", "", "output directory")
	flag.BoolVar(&extrude, "e", false, "extrude image frames by 1px each side")
	flag.Parse()

	if len(outDir) == 0 {
		log.Fatalf("-o not specified")
	}

	if len(voxDir) == 0 {
		log.Fatalf("-p not specified")
	}

	if extrude {
		log.Printf("extruding all images by 1px")
	}

	voxes, err := enumer(voxDir)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	for i := range voxes {
		writeVox(voxes[i], fmt.Sprintf("%s/vox_%d", outDir, i), extrude)
	}
}

func enumer(dir string) ([]*vox, error) {
	var uf []*vox

	log.Printf("looking for files in %s", dir)
	err := filepath.Walk(dir, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fi.IsDir() && filepath.Ext(p) == ".vox" {
			log.Printf("\t%s", p)

			fp, err := os.Open(p)
			if err != nil {
				return fmt.Errorf("err opening %s: %s", p, err)
			}

			v, err := newVox(fp)
			if err != nil {
				return fmt.Errorf("err reading vox %s: %s", p, err)
			}

			uf = append(uf, v)
		}

		return nil
	})

	return uf, err
}
