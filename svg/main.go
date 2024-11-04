package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type SVG struct {
	XMLName xml.Name `xml:"svg"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	Content []byte   `xml:",innerxml"`
}

func cropSVG(inputFile, outputFile string) error {
	// Read the input SVG file
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	// Parse the SVG
	var svg SVG
	if err := xml.Unmarshal(data, &svg); err != nil {
		return fmt.Errorf("failed to unmarshal SVG: %v", err)
	}

	// Get the dimensions
	width := parseDimension(svg.Width)
	height := parseDimension(svg.Height)

	// Determine the size of the square
	size := width
	if height < width {
		size = height
	}

	// Create the new SVG content
	newSVG := fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d">%s</svg>`, size, size, size, size, svg.Content)

	// Write the new SVG to the output file
	if err := ioutil.WriteFile(outputFile, []byte(newSVG), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	return nil
}

func parseDimension(dim string) int {
	var value int
	fmt.Sscanf(dim, "%d", &value)
	return value
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <inputFile> <outputFile>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cropSVG(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
