package main

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"sync"

	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/theonlyjohnny/galactus.git/pkg/krea"
)

const (
	imageURLBase          = "https://galactus-frames.s3.us-east-1.amazonaws.com/"
	defaultPrompt         = "planets, astrophotography, an array of alien planets"
	defaultNegativePrompt = "bad quality, low resolution, image with artifacts, words, characters"
)

func convertToBinary(ascii int) string {
	return strconv.FormatInt(int64(ascii), 2)
}

func downloadAndSaveImage(url, outputDir, fileName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP GET failed for %s: Status %d", url, resp.StatusCode)
	}

	fullPath := filepath.Join(outputDir, fileName)
	os.MkdirAll(outputDir, os.ModePerm)
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func generateGalaxy(inputString string) error {
	client := krea.New()
	outputDir := "./output/" + inputString
	if err := os.MkdirAll(outputDir, 0777); err != nil {
		return err
	}

	fmt.Printf("Generating galaxy for string: %s\n", inputString)

	var wg sync.WaitGroup
	wg.Add(len(inputString))

	// Iterate through each character in the input string.
	for i, char := range inputString {
		go func(client krea.Client, i int, char rune, outputDir string) {
			if err := generateImage(client, i, char, outputDir); err != nil {
				fmt.Printf("Error generating image: %s\n", err.Error())
			}
			wg.Add(-1)
		}(client, i, char, outputDir)
	}

	wg.Wait()

	return nil
}

func generateImage(client krea.Client, i int, char rune, outputDir string) error {
	ascii := int(char)
	binaryString := convertToBinary(ascii)

	// Call the GenImage function with the binary representation.
	// Replace {binary representation} with the actual binaryString and set up your client and krea.GenImageOpts.
	imageURL, err := client.GenImage(krea.GenImageOpts{
		ImageResolution:             "1:1",
		PromptInfluence:             0,
		ControlnetConditioningScale: 1.5,
		DiffusionSteps:              30,
		PatternURL:                  1.5,
		ImageURL:                    fmt.Sprintf("%s%d.png", imageURLBase, ascii),
		Prompt:                      defaultPrompt,         // TODO make this configurable
		NegativePrompt:              defaultNegativePrompt, // TODO make this configurable
	})

	fmt.Printf("Generating image for char #%d: %c, (%s)\n", i, char, binaryString)

	if err != nil {
		return errors.New(fmt.Sprintf("Error for character '%c' (ASCII %d): %v\n", char, ascii, err))
	} else {
		// Download the image and save it in the output folder.
		err := downloadAndSaveImage(imageURL, outputDir, fmt.Sprintf("%d.png", i))
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to download and save image for character '%c' (ASCII %d): %v\n", char, ascii, err))
		} else {
			fmt.Printf("Saved image for character '%c' (ASCII %d) to %s\n", char, ascii, outputDir)
		}
	}

	return nil
}

func main() {

	rootCmd := &cobra.Command{
		Use:   "galactus",
		Short: "Galactus - text2img encoder",
		Long:  "A simple experiment on top of Krea.ai to encode a string into a planetary array.",
	}
	rootCmd.PersistentFlags().StringP("string", "s", "", "Input string")
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		inputString, _ := cmd.Flags().GetString("string")
		if inputString == "" {
			return errors.New("Usage: myapp --string 'your string here'")
		} else {
			return generateGalaxy(inputString)
		}
	}

	// Add the root command to the command-line interface.
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
