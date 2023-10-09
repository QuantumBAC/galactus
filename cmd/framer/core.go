package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/fogleman/gg"
)

const (
	bitCount = 9 //NOTE: must be a perfect square. could be moved to query param in future
)

var gridSize = math.Sqrt(bitCount)

func convertToBinary(input int) string {
	return strconv.FormatInt(int64(input), 2)
}

func generateImage(binaryString string) (*gg.Context, error) {
	fmt.Printf("generating image for string: %s\n", binaryString)
	const (
		imageSize = 256
		lineWidth = 2.0
	)

	radius := imageSize / gridSize / 2

	dc := gg.NewContext(imageSize, imageSize)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetLineWidth(lineWidth)

	structuredString := [bitCount]rune{}
	for i, b := range binaryString {
		structuredString[i] = b
	}

	for i := float64(0); i < bitCount; i++ {

		y := math.Floor(float64(i) / gridSize)
		x := i - (y * gridSize)

		geoY := y + radius*imageSize/gridSize
		geoX := x + radius*imageSize/gridSize

		fmt.Printf("drawing bit #%f at (%f, %f)\n", i, x, y)
		dc.SetRGB(1, 1, 1)

		bit := structuredString[int(i)]
		if bit == '1' {
			dc.DrawCircle(geoX, geoY, radius)
			// dc.DrawRectangle(x*imageSize/gridSize, y*imageSize/gridSize, imageSize/gridSize, imageSize/gridSize)
			dc.Fill()
		} else {
			dc.DrawEllipse(geoX, geoY, radius, radius)
			dc.Stroke()
		}

	}

	return dc, nil
}
