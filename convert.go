package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"

	"github.com/chai2010/tiff"
	"github.com/disintegration/imaging"
	"github.com/jung-kurt/gofpdf"
)

func main() {
	var data []byte
	var err error

	var files = []string{
		"./testdata/test.tiff",
	}
	var layers []string
	var corners []struct {
		x int
		y int
	}
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
				var buf bytes.Buffer
				if err = png.Encode(&buf, img_rev); err != nil {

					log.Fatal(err)
				}
				layers = append(layers, newname)
				corners = append(corners, struct{ x, y int }{x: img.Bounds().Dx(), y: img.Bounds().Dy()})
				if err = ioutil.WriteFile(newname, buf.Bytes(), 0666); err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Save %s ok\n", newname)
			}
		}
	}
	fmt.Println(layers)

	imagesToPdf(layers, corners, "result.pdf")
}

func imagesToPdf(layers []string, corners []struct{ x, y int }, name string) {
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "pt",
		Size:    gofpdf.SizeType{Wd: 1920, Ht: 1080},
	})

	fmt.Println("layers:", len(layers)-1)
	for i := 1; i < len(layers); i++ {
		pdf.AddPage()

		x := float64(1920-corners[i].x) / 2
		y := float64(1080-corners[i].y) / 2

		fmt.Println(x, y)
		pdf.Image(layers[i], x, y, 0, 0, false, "", 0, "")
	}
	err := pdf.OutputFileAndClose(name)
	if err != nil {
		log.Fatal(err)
	}
}
