package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	count := make(map[Pixel]uint32)
	// order := make(map[uint32]Pixel)

	var pixelCountReducer = func(p Pixel) {
		if val, ok := count[p]; ok {
			count[p] = val + 1
		} else {
			count[p] = 1
		}
	}

	rd, err := ioutil.ReadDir("./images")
	if err != nil {
		fmt.Println("read dir fail:", err)
	}

	for _, pic := range rd {
		name := "./images/" + pic.Name()
		fmt.Println(name)
		file, _ := os.Open(name)

		iterateImage(file, pixelCountReducer)
	}

	fmt.Println(count)

}

func iterateImage(file io.Reader, callback func(p Pixel)) error {
	img, _, err := image.Decode(file)

	if err != nil {
		return err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			callback(rgbaToPixel(img.At(x, y).RGBA()))
		}
	}

	return nil
}

// Get the bi-dimensional pixel array
func getPixels(file io.Reader) ([][]Pixel, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := bounds.Min.Y; y < height; y++ {
		var row []Pixel
		for x := bounds.Min.X; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

// Pixel struct example
type Pixel struct {
	R int
	G int
	B int
	A int
}
