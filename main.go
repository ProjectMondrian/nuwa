package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func main() {
	imagePath := "./images/1.jpg"
	file, _ := os.Open(imagePath)

	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	b := img.Bounds()
	width := b.Max.X
	height := b.Max.Y

	fmt.Println("width = ", width)
	fmt.Println("height = ", height)

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			r, g, b, a := img.At(i, j).RGBA()
			println(i, j, r, g, b, a)
		}
	}

}
