package main

import (
	"flag"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	input  = flag.String("i", "", "input file (required)")
	output = flag.String("o", "output.jpg", "output file")
	width  = flag.Int("w", 64, "width of the grid")
	height = flag.Int("h", 64, "height of the grid")
)

func copyImage(src image.Image) (c draw.Image) {
	size := src.Bounds().Size()
	r := image.Rect(0, 0, size.X, size.Y)

	switch src.(type) {
	case *image.Alpha:
		c = image.NewAlpha(r)
	case *image.Alpha16:
		c = image.NewAlpha16(r)
	case *image.CMYK:
		c = image.NewCMYK(r)
	case *image.Gray:
		c = image.NewGray(r)
	case *image.Gray16:
		c = image.NewGray16(r)
	case *image.NRGBA:
		c = image.NewNRGBA(r)
	case *image.NRGBA64:
		c = image.NewNRGBA64(r)
	case *image.RGBA:
		c = image.NewRGBA(r)
	case *image.RGBA64:
		c = image.NewRGBA64(r)
	default:
		c = image.NewRGBA(r)
	}
	draw.Draw(c, r, src, src.Bounds().Min, draw.Src)
	return
}

func getGrid(x, y, w, h int) image.Rectangle {
	minPoint := image.Point{x * w, y * h}
	maxPoint := minPoint.Add(image.Point{w, h})
	return image.Rectangle{minPoint, maxPoint}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func openFile(name string) *os.File {
	f, err := os.Open(name)
	if err != nil {
		log.Panic(err)
	}
	return f
}

func ParseFlags() {
	flag.Parse()
	if *input == "" {
		flag.Usage()
		os.Exit(2)
	}
}

func DecodeImage(name string) image.Image {
	f := openFile(name)
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Panic(err)
	}
	return img
}

func GridTranspose(img image.Image, w, h int) image.Image {
	size := img.Bounds().Size()
	log.Printf("The size of the input: %dx%d", size.X, size.Y)

	t := copyImage(img)

	count := min(size.X/w, size.Y/h)
	for y := 0; y < count; y++ {
		for x := 0; x < count; x++ {
			sp := img.Bounds().Min.Add(image.Point{y * w, x * h})
			grid := getGrid(x, y, w, h)
			draw.Draw(t, grid, img, sp, draw.Src)
		}
	}
	return t
}

func SaveImage(name string, img image.Image) (err error) {
	w, err := os.Create(name)
	if err != nil {
		log.Panic(err)
	}
	defer w.Close()

	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".png":
		err = png.Encode(w, img)
	case ".gif":
		err = gif.Encode(w, img, nil)
	default:
		// Uses the jpeg encoder for jpg, jpeg or other extensions
		err = jpeg.Encode(w, img, nil)
	}
	if err == nil {
		log.Printf("Saved as %s", name)
	} else {
		log.Printf("Failed to save img as %s", name)
	}
	return
}

func main() {
	ParseFlags()
	img := DecodeImage(*input)
	t := GridTranspose(img, *width, *height)
	SaveImage(*output, t)
}
