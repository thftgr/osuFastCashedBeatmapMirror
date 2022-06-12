package main

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"time"

	// Package image/jpeg is not used explicitly in the code below,
	// but is imported for its initialization side-effect, which allows
	// image.Decode to understand JPEG formatted images. Uncomment these
	// two lines to also understand GIF and PNG images:
	// _ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func main() {
	// Decode the JPEG data. If reading from file, create a reader with
	//
	reader, err := os.Open("./screenshot380.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	st := time.Now().UnixMilli()
	//reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	m, err := png.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	//mm := image.NewRGBA(m.Bounds())
	//0.0625
	dst := image.NewRGBA(image.Rectangle{image.Pt(0, 0), image.Pt(m.Bounds().Size().X, int(math.Round(float64(m.Bounds().Size().Y)*0.11111)))})
	draw.Draw(dst, m.Bounds(), m, image.Point{}, draw.Src)

	//target := image.Rect(0, 120, dst.Rect.Size().X, dst.Rect.Size().Y)
	//draw.Draw(dst, target, &image.Uniform{C: color.RGBA{
	//	R: 0,
	//	G: 0,
	//	B: 0,
	//	A: 255,
	//}}, image.Point{}, draw.Src)
	et := time.Now().UnixMilli()
	log.Println(et-st, "ms")
	//draw.Draw
	f, err := os.Create("./screenshot380-" + time.Now().Format("150405") + ".png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err = png.Encode(f, dst); err != nil {
		log.Printf("failed to encode: %v", err)
	}

}
