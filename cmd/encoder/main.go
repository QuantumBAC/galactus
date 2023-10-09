package main

import (
	"fmt"

	"github.com/theonlyjohnny/galactus.git/pkg/krea"
)

func main() {

	client := krea.New()
	imageURL, err := client.GenImage(krea.GenImageOpts{
		ImageResolution:             "1:1",
		PromptInfluence:             0,
		ControlnetConditioningScale: 1.5,
		DiffusionSteps:              30,
		PatternURL:                  1.5,
		ImageURL:                    "https://imagedelivery.net/1ddUs0BD8AGrZYOYoauaWw/canvas-demo/2d914d0e-de7c-4071-b194-6668533baa3d/public",
		Prompt:                      "planets, astrophotography, an array of alien planets",
		NegativePrompt:              "bad quality, low resolution, image with artifacts, words, characters",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n final imageURL: %s\n", imageURL)
}
