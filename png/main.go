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

func makeBackgroundTransparent(inputPath, outputPath string) error {
	// Open the input image file
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the image
	img, err := png.Decode(file)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	// Define the background color (white in this case) and tolerance
	bgColor := color.RGBA{255, 255, 255, 255}
	tolerance := uint8(89) // Adjust the tolerance as needed

	// Process each row of pixels
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		firstNonWhite := -1
		lastNonWhite := -1

		// Identify the first and last non-white pixels in the row
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()

			if !withinTolerance(uint8(r>>8), bgColor.R, tolerance) ||
				!withinTolerance(uint8(g>>8), bgColor.G, tolerance) ||
				!withinTolerance(uint8(b>>8), bgColor.B, tolerance) ||
				!withinTolerance(uint8(a>>8), bgColor.A, tolerance) {
				if firstNonWhite == -1 {
					firstNonWhite = x
				}
				lastNonWhite = x
			}
		}

		// Update the row based on the identified bounds
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if x < firstNonWhite || x > lastNonWhite {
				newImg.Set(x, y, color.Transparent)
			} else {
				newImg.Set(x, y, img.At(x, y))
			}
		}
	}

	// Create the output image file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Encode the new image to the output file
	err = png.Encode(outFile, newImg)
	if err != nil {
		return err
	}

	return nil
}

func withinTolerance(value, target, tolerance uint8) bool {
	return value >= target-tolerance || value <= target+tolerance
}

func makeBackgroundWhite(inputPath, outputPath string) error {
	// Open the input image file
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the image
	img, err := png.Decode(file)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	// Define transparency tolerance
	alphaTolerance := uint8(50) // Pixels with alpha below this are considered transparent

	// Process each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			_, _, _, a := c.RGBA()
			alpha := uint8(a >> 8)

			// If the pixel is transparent, make it white
			if alpha < alphaTolerance {
				newImg.Set(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				newImg.Set(x, y, c)
			}
		}
	}

	// Create the output image file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Encode the new image to the output file
	err = png.Encode(outFile, newImg)
	if err != nil {
		return err
	}

	return nil
}

func makeNonTransparentPixelColored(inputPath, outputPath string) error {
	// Open the input image file
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the image
	img, err := png.Decode(file)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	// Define transparency threshold
	alphaThreshold := uint8(50) // Pixels with alpha above this are considered non-transparent

	// Kotlin gradient colors: purple to magenta/pink
	purpleR, purpleG, purpleB := float64(91), float64(75), float64(138)
	magentaR, magentaG, magentaB := float64(245), float64(50), float64(180)

	width := float64(bounds.Dx())
	height := float64(bounds.Dy())

	// Process each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			_, _, _, a := c.RGBA()
			alpha := uint8(a >> 8)

			// If the pixel is non-transparent, apply gradient
			if alpha >= alphaThreshold {
				// Calculate gradient position based on x and y (diagonal from bottom-left to top-right)
				tx := float64(x-bounds.Min.X) / width
				ty := 1.0 - (float64(y-bounds.Min.Y) / height) // Invert y so top is 1
				t := (tx + ty) / 2.0                           // Average for diagonal effect

				if t > 1 {
					t = 1
				}
				if t < 0 {
					t = 0
				}

				// Interpolate between purple and red
				r := uint8(purpleR + (magentaR-purpleR)*t)
				g := uint8(purpleG + (magentaG-purpleG)*t)
				b := uint8(purpleB + (magentaB-purpleB)*t)

				newImg.Set(x, y, color.RGBA{r, g, b, 255})
			} else {
				newImg.Set(x, y, c)
			}
		}
	}

	// Create the output image file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Encode the new image to the output file
	err = png.Encode(outFile, newImg)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// ===== SECTION: add margins =====

	// if len(os.Args) < 4 {
	// 	fmt.Println("Usage: go run main.go <inputFile> <outputFile> <margin>")
	// 	return
	// }

	// inputFile := os.Args[1]
	// outputFile := os.Args[2]
	// margin := parseDimension(os.Args[3])

	// if err := addMarginsToPNG(inputFile, outputFile, margin); err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// }

	// ===== SECTION: make background transparent =====

	// if len(os.Args) < 3 {
	// 	fmt.Println("Usage: go run main.go <inputFile> <outputFile>")
	// 	return
	// }

	// inputFile := os.Args[1]
	// outputFile := os.Args[2]

	// if err := makeBackgroundTransparent(inputFile, outputFile); err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// }

	// ===== SECTION: make background white =====

	// if len(os.Args) < 3 {
	// 	fmt.Println("Usage: go run main.go <inputFile> <outputFile>")
	// 	return
	// }

	// inputFile := os.Args[1]
	// outputFile := os.Args[2]

	// if err := makeBackgroundWhite(inputFile, outputFile); err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// }

	// ===== SECTION: make non-transparent pixels colored =====

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <inputFile> <outputFile>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := makeNonTransparentPixelColored(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
