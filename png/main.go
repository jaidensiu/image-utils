package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
)

func addMarginsToPNG(inputFile, outputFile string, margin int) error {
	// Open the input file
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	// Decode the PNG image
	img, err := png.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode PNG image: %v", err)
	}

	// Get the original image dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create a new image with margins
	newWidth := width + 2*margin
	newHeight := height + 2*margin
	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Fill the new image with color
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	// Draw the original image onto the new image with margins
	draw.Draw(newImg, image.Rect(margin, margin, margin+width, margin+height), img, bounds.Min, draw.Over)

	// Create the output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Encode the new image to the output file
	if err := png.Encode(outFile, newImg); err != nil {
		return fmt.Errorf("failed to encode PNG image: %v", err)
	}

	return nil
}

func parseDimension(dim string) int {
	value, err := strconv.Atoi(dim)
	if err != nil {
		fmt.Printf("Invalid dimension: %v\n", err)
		return 0
	}
	return value
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <inputFile> <outputFile> <margin>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	margin := parseDimension(os.Args[3])

	if err := addMarginsToPNG(inputFile, outputFile, margin); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
