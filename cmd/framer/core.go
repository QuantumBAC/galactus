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
		imageSize      = 256
		lineWidth      = 2.0
		paddingPercent = 0.1
	)
	bitSize := imageSize / gridSize
	padding := paddingPercent * bitSize
	radius := bitSize / 2

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

		geoY := (y * bitSize)
		geoX := (x * bitSize)

		fmt.Printf("drawing bit #%f at (%f, %f)\n", i, x, y)
		dc.SetRGB(1, 1, 1)

		bit := structuredString[int(i)]
		if bit == '1' {
			// fmt.Println("geoX,geoY,radius", geoX, geoY, radius)
			dc.DrawCircle(geoX+radius, geoY+radius, radius-padding)
			// dc.DrawRectangle(geoX, geoY, bitSize, bitSize)
			dc.Fill()
		} else {
			dc.DrawEllipse(geoX+radius, geoY+radius, radius-padding, radius-padding)
			dc.Stroke()
		}

	}

	return dc, nil
}
