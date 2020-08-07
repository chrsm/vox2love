package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	voxDir string // path/to/src
	outDir string // path/to/dest
)

func main() {
	log.Println("vox2love")

	flag.StringVar(&voxDir, "p", "", "path to directory containin vox files")
	flag.StringVar(&outDir, "o", "", "output directory")
	flag.Parse()

	if len(outDir) == 0 {
		log.Fatalf("-o not specified")
	}

	if len(voxDir) == 0 {
		log.Fatalf("-p not specified")
	}

	voxes, err := enumer(voxDir)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	for i := range voxes {
		writeLovox(voxes[i], fmt.Sprintf("%s/vox_%d", outDir, i))
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
