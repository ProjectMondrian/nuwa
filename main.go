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
	"sort"
)

func main() {
	count := make(map[Pixel]float32)
	order := make(CountedPixels, 10)

	mapRet := make([]map[Pixel]float32, 10)

	// 回调函数
	var pixelCountReducer = func(p Pixel, c map[Pixel]float32) {
		if val, ok := c[p]; ok {
			c[p] = val + 1
		} else {
			c[p] = 1
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

		mapRet = append(mapRet, iterateImage(file, pixelCountReducer))
	}

	// 将每个 image 的 map 结果汇总到一个 map 里
	for _, ret := range mapRet {
		for pixel, c := range ret {
			if val, ok := count[pixel]; ok {
				count[pixel] = val + c
			} else {
				count[pixel] = c
			}
		}
	}

	// 将汇总结果 map 转为 slice，并排序
	for key, value := range count {
		order = append(order, CountedPixel{Count: value, Pixel: key})
	}
	sort.Sort(order)
	fmt.Println(order)
}

func iterateImage(file io.Reader, callback func(p Pixel, c map[Pixel]float32)) map[Pixel]float32 {
	count := make(map[Pixel]float32)
	img, _, err := image.Decode(file)

	if err != nil {
		panic(err)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			callback(RgbaToPixel(img.At(x, y).RGBA()), count)
		}
	}
	for pixel, value := range count {
		count[pixel] = value / float32(height) / float32(width) * 100 * 100
	}
	return count
}

// Pixel struct example
type Pixel struct {
	R int
	G int
	B int
	A int
}

// CountedPixel ...
type CountedPixel struct {
	Count float32
	Pixel Pixel
}

// CountedPixels ...
type CountedPixels []CountedPixel

// Len ...
func (c CountedPixels) Len() int {
	return len(c)
}

// Swap ...
func (c CountedPixels) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Len ...
func (c CountedPixels) Less(i, j int) bool {
	return c[i].Count > c[j].Count
}

// GetPixels  Get the bi-dimensional pixel array
func GetPixels(file io.Reader) ([][]Pixel, error) {
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
			row = append(row, RgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}

// RgbaToPixel  img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func RgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}
