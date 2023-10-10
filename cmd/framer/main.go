package main

import (
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "framer", Short: "Generate masks for use in patterned image generation"}

func init() {
	charRangeCmd.Flags().StringSliceVarP(&ranges, "range", "r", defaultRanges, "Character ranges (e.g., 'a-z' or 'A-F')")

	rootCmd.AddCommand(charRangeCmd)
	rootCmd.AddCommand(httpCmd)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	cobra.CheckErr(rootCmd.Execute())
}
