package main

import (
	"fmt"
	"math"
	"math/rand"
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

	// first we'll define some constants about our frame
	const (
		imageSize     = 256
		lineWidth     = 2.0
		marginPercent = 0.1
	)

	// bitSize is the size each bit is allowed to fill
	bitSize := imageSize / gridSize

	// margin is the amount of cells that are dedicated to keeping each bit separate
	margin := marginPercent * bitSize

	// we initialize the frame context
	dc := gg.NewContext(imageSize, imageSize)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetLineWidth(lineWidth)

	// let's ensure we have the right amount of bits
	structuredString := [bitCount]rune{}
	for i, b := range binaryString {
		structuredString[i] = b
	}

	// now we draw each bit
	for i := float64(0); i < bitCount; i++ {

		// we calculate the desired position for each bit in a theoretical plane
		y := math.Floor(float64(i) / gridSize)
		x := i - (y * gridSize)
		// radius is half of diameter. bitSize is the full width
		radius := bitSize / 2

		noiseX := math.Floor(rand.Float64() * margin)
		noiseY := math.Floor(rand.Float64() * margin)

		// now we convert that theoretical plane into proper units for use in this canvas context
		geoY := (y * bitSize) + noiseY
		geoX := (x * bitSize) + noiseX

		fmt.Printf("drawing bit #%f at (%f, %f)\n", i, x, y)
		dc.SetRGB(1, 1, 1)

		bit := structuredString[int(i)]
		if bit == '1' {
			// if the bit is truthy, we want to draw a filled in circle
			dc.DrawCircle(geoX+radius, geoY+radius, radius-margin)
			// dc.DrawRectangle(geoX, geoY, bitSize, bitSize)
			dc.Fill()
		} else {
			// if the bit is falsey, we want to draw an empty circle

			dc.DrawEllipse(geoX+radius, geoY+radius, radius-margin, radius-margin)
			dc.Stroke()
		}
	}

	return dc, nil
}
