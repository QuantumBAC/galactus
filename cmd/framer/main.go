package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

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
	const imageSize = 256
	dc := gg.NewContext(imageSize, imageSize)
	dc.SetRGB(0, 0, 0)
	dc.Clear()

	structuredString := [bitCount]rune{}
	for i, b := range binaryString {
		structuredString[i] = b
	}

	for i := float64(0); i < bitCount; i++ {

		y := math.Floor(float64(i) / gridSize)
		x := i - (y * gridSize)

		bit := structuredString[int(i)]
		if bit == '1' {
			fmt.Printf("drawing bit #%f at (%f, %f)\n", i, x, y)
			dc.SetRGB(1, 1, 1)
			dc.DrawRectangle(x*imageSize/gridSize, y*imageSize/gridSize, imageSize/gridSize, imageSize/gridSize)
			dc.Fill()
		}

	}

	return dc, nil
}

func handleImageRequest(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	input, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	binaryString := convertToBinary(input)
	dc, err := generateImage(binaryString)
	if err != nil {
		http.Error(w, "Image generation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	err = dc.EncodePNG(w)
	if err != nil {
		http.Error(w, "Image encoding failed", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/mask/", handleImageRequest)
	port := 8080
	fmt.Printf("Server is listening on :%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
