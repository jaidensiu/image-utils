package main

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
)

func convertJPGToPNG(inputPath, outputPath string) error {
	// Open the input JPG file
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the JPG image
	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}

	// Create the output PNG file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Encode the image to the output file as PNG
	err = png.Encode(outFile, img)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <inputFile> <outputFile>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := convertJPGToPNG(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
