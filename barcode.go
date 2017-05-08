package goboleto

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/twooffive"
	"github.com/disintegration/imaging"
	"github.com/llgcode/draw2d/draw2dimg"
)

func GerarBarcode2of5(code string) (io.Reader, error) {
	// create interleaved barcode, set to true
	// see https://godoc.org/github.com/boombuler/barcode/twooffive#Encode
	bcode, err := twooffive.Encode(code, true)

	if err != nil {
		fmt.Printf("String %s cannot be encoded: %s", code, err.Error())
		return nil, err
	}

	bcode, err = barcode.Scale(bcode, 1000, 50)

	if err != nil {
		fmt.Println("Two by five barcode scaling error!", err.Error())
		return nil, err
	}

	// Initialize the graphic context on an RGBA image
	img := image.NewRGBA(image.Rect(0, 0, 1000, 50))

	// set background to white
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)

	gc := draw2dimg.NewGraphicContext(img)

	gc.FillStroke()

	gc.SetFillColor(image.Black)
	newImg := imaging.New(1000, 50, color.NRGBA{255, 255, 255, 255})

	//paste the codabar to new blank image
	newImg = imaging.Paste(newImg, bcode, image.Pt(0, 0))

	buf := new(bytes.Buffer)

	// Write the image into the buffer
	err = png.Encode(buf, newImg)
	if err != nil {
		return nil, err
	}


	// Para salvar o código de barras para análise, descomentar as linhas abaixo

	//f, _ := os.Create("barcode.png")
	//defer f.Close()
	//f.Write(buf.Bytes())



	return buf, nil
}

