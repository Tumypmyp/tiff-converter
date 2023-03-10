package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/chai2010/tiff"
	"github.com/signintech/gopdf"
)

func main() {
	var data []byte
	var err error

	if len(os.Args) <= 1 {
		log.Fatal("no input files\n")
	}
	var files = os.Args[1:]

	for _, filename := range files {
		var images []image.Image
		// Load file data
		if data, err = ioutil.ReadFile(filename); err != nil {
			log.Fatal("readfile:", err)
		}

		// Decode tiff
		m, errors, err := tiff.DecodeAll(bytes.NewReader(data))
		if err != nil {
			log.Println("decode:", err)
		}

		// Get layers
		for i := 0; i < len(m); i++ {
			for j := 1; j < len(m[i]); j++ {
				img, ok := m[i][j].(*image.RGBA)
				if !ok {
					log.Printf("layer %v %v: cant convert to RGBA\n", i, j)
					continue

				}
				if img.Bounds().Dx() <= 200 && img.Bounds().Dy() <= 120 {
					continue
				}
				fmt.Printf("%v %v %v\n", j, reflect.TypeOf(img), img.Bounds())
				fmt.Printf("%v\n", img.Opaque())

				if errors[i][j] != nil {
					log.Printf("%v %v got error: %v\n", i, j, err)
					continue
				}

				images = append(images, img)
			}
		}

		pdfName := filename[:len(filename)-len(filepath.Ext(filename))] + ".pdf"
		fmt.Printf("%v layers\n", len(images))
		encodeToPdf(images, pdfName)
	}
}

func encodeToPdf(images []image.Image, name string) {
	pdf := gopdf.GoPdf{}
	conf := gopdf.Config{Unit: gopdf.Unit_PT, PageSize: gopdf.Rect{W: 1920, H: 1080}}
	pdf.Start(conf)
	for _, img := range images {
		pdf.AddPage()

		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			log.Fatal(err)
		}
		imgH, err := gopdf.ImageHolderByReader(&buf)
		if err != nil {
			log.Fatal(err)
		}
		x := (conf.PageSize.W - float64(img.Bounds().Dx())) / 2
		y := (conf.PageSize.H - float64(img.Bounds().Dy())) / 2
		if err := pdf.ImageByHolderWithOptions(imgH,
			gopdf.ImageOptions{X: x, Y: y, VerticalFlip: true}); err != nil {
			log.Fatal(err)
		}

	}
	pdf.WritePdf(name)
	fmt.Printf("written to %v\n", name)
}
