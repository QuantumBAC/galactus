package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var ranges []string
var defaultRanges = []string{
	"a-z",
	"A-Z",
	"0-9",
}

var charRangeCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate frames within specified ranges",
	Run: func(cmd *cobra.Command, args []string) {
		for _, r := range ranges {
			min, max := parseRange(r)
			if min > max {
				fmt.Println("Invalid range:", r)
			} else {
				if err := genFrames(min, max); err != nil {
					fmt.Printf("Error generating range %c-%c: %s \n", min, max, err)
				}
			}
		}
	},
}

// charRangeCmd.Flags().StringSliceVarP(&ranges, "range", "r", nil, "Character ranges (e.g., 'a-z' or 'A-F')")

func genFrames(startChar, endChar rune) error {
	fmt.Printf("Generating frames for range %c - %c\n", startChar, endChar)
	for char := startChar; char <= endChar; char++ {
		ascii := int(char) // Convert character to ASCII code

		binaryString := convertToBinary(ascii)
		dc, err := generateImage(binaryString)
		if err != nil {
			return err
		}

		imageFileName := fmt.Sprintf("./frames/%d.png", ascii)
		outputFile, err := os.Create(imageFileName)
		if err != nil {
			return err
		}
		defer outputFile.Close()

		// Encode the context as a PNG and save it to the output file.
		err = dc.EncodePNG(outputFile)

		if err != nil {
			return errors.New(fmt.Sprintf("Failed to write image to file for character %c (ASCII %d): %v\n", char, ascii, err))
		} else {
			fmt.Printf("Saved image for character %c (ASCII %d) to %s\n", char, ascii, imageFileName)
		}
	}
	return nil
}

func parseRange(r string) (min, max rune) {
	parts := strings.Split(r, "-")

	if len(parts) == 2 {
		min = []rune(parts[0])[0]
		max = []rune(parts[1])[0]
	}
	return
}
