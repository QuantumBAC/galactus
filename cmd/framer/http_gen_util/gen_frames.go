package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	genFrames('0', '9')
	genFrames('A', 'Z')
	genFrames('a', 'z')
}

func genFrames(startChar, endChar byte) {
	for char := startChar; char <= endChar; char++ {
		ascii := int(char) // Convert character to ASCII code
		url := fmt.Sprintf("http://localhost:8080/mask/%d", ascii)

		// Perform an HTTP GET request to the URL.
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error for character %c (ASCII %d): %v\n", char, ascii, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("HTTP GET failed for character %c (ASCII %d): Status %d\n", char, ascii, resp.StatusCode)
			continue
		}

		// Read the response body (image data).
		imageData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed to read response body for character %c (ASCII %d): %v\n", char, ascii, err)
			continue
		}

		// Save the image to the local filesystem.
		imageFileName := fmt.Sprintf("./frames/%d.png", ascii)
		err = ioutil.WriteFile(imageFileName, imageData, 0644)
		if err != nil {
			fmt.Printf("Failed to write image to file for character %c (ASCII %d): %v\n", char, ascii, err)
		} else {
			fmt.Printf("Saved image for character %c (ASCII %d) to %s\n", char, ascii, imageFileName)
		}
	}
}
