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

func drawSpiral(dc *gg.Context, x, y, maxRadius float64) {
	// Constants for the spiral
	turns := 5.0    // Number of spiral turns
	segments := 100 // Number of line segments to approximate the spiral

	// Calculate the angle increment based on the number of segments and turns
	angleIncrement := 360.0 * turns / float64(segments)

	// Start drawing the spiral
	dc.MoveTo(x, y)
	for i := 0; i < segments; i++ {
		// Calculate the current angle in radians
		angle := float64(i) * angleIncrement * math.Pi / 180

		// Calculate the current maxRadius based on the angle
		currentRadius := (float64(i) / float64(segments-1)) * maxRadius

		if currentRadius > maxRadius {
			currentRadius = maxRadius
		}

		// Calculate the coordinates for the current point on the spiral
		currentX := x + currentRadius*math.Cos(angle)
		currentY := y + currentRadius*math.Sin(angle)

		// fmt.Printf("currentX: %v, currentY: %v, currentRadius: %v, angle: %v \n", currentX, currentY, currentRadius, angle)

		dc.LineTo(currentX, currentY)
	}

	// Stroke the path to actually draw the spiral
	dc.Stroke()
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

		// fmt.Println("margin, margin, noiseX, noiseY", margin, margin, noiseX, noiseY)

		// now we convert that theoretical plane into proper units for use in this canvas context
		geoY := (y * bitSize) + noiseY
		geoX := (x * bitSize) + noiseX

		fmt.Printf("drawing bit #%f at (%f, %f)\n", i, x, y)
		dc.SetRGB(1, 1, 1)

		bit := structuredString[int(i)]
		if bit == '1' {
			// if the bit is truthy, we want to draw a filled in circle

			// fmt.Println("geoX,geoY,radius", geoX, geoY, radius)
			drawSpiral(dc, geoX+radius, geoY+radius, radius-margin)
			// dc.DrawCircle(geoX+radius, geoY+radius, radius-margin)
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
