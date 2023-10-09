package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Start an HTTP server to generate masks",
	Run: func(cmd *cobra.Command, args []string) {

		http.HandleFunc("/mask/", handleImageRequest)
		port := 8080
		fmt.Printf("Server is listening on :%d\n", port)
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	},
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
