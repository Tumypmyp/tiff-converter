package main

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"

	"github.com/chai2010/tiff"
	"github.com/disintegration/imaging"
	"github.com/signintech/gopdf"
)

func main() {
	var data []byte
	var err error

	var files = []string{
		"./testdata/test.tiff",
	}
	var images []image.Image
	for _, filename := range files {
		// Load file data
		if data, err = ioutil.ReadFile(filename); err != nil {
			log.Fatal(err)
		}

		// Decode tiff
		m, errors, err := tiff.DecodeAll(bytes.NewReader(data))
		if err != nil {
			log.Println(err)
		}

		// Encode tiff
		for i := 0; i < len(m); i++ {
			for j := 0; j < len(m[i]); j++ {

				img, ok := m[i][j].(*image.RGBA)
				//fmt.Println(img0.Stride, img0.Rect)
				//img := &image.RGBA64{Pix: img0.Pix, Stride: img0.Stride, Rect: img0.Rect}
				if !ok {
					log.Fatal("cant convert to RGBA")
				}
				fmt.Printf("%v %v\n", reflect.TypeOf(img), img.Bounds())
				fmt.Printf("%v\n", img.Opaque())

				if img.Bounds().Dx() <= 200 && img.Bounds().Dy() <= 120 {
					continue
				}

				newname := fmt.Sprintf("%s-%02d-%02d.png", filepath.Base(filename), i, j)
				if errors[i][j] != nil {
					log.Printf("%s: %v\n", newname, err)
					continue
				}

				img_rev := imaging.FlipV(img)
				images = append(images, img_rev)
			}
		}
	}

	fmt.Printf("%v layers\n", len(images))
	imagesToPdf(images, "result.pdf")
}

func imagesToPdf(images []image.Image, name string) {
	pdf := gopdf.GoPdf{}
	conf := gopdf.Config{Unit: gopdf.Unit_PT, PageSize: gopdf.Rect{W: 1920, H: 1080}}
	pdf.Start(conf)
	for _, img := range images {
		pdf.AddPage()

		x := (conf.PageSize.W - float64(img.Bounds().Dx())) / 2
		y := (conf.PageSize.H - float64(img.Bounds().Dy())) / 2
		if err := pdf.ImageFrom(img, x, y, nil); err != nil {
			log.Fatal(err)
		}

	}
	pdf.WritePdf(name)
}
