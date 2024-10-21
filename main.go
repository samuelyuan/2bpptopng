package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	inputPtr := flag.String("input", "", "Input filename")
	rowPtr := flag.Int("row", 1, "Number of rows of 8x8 blocks")

	flag.Parse()

	inputFilename := *inputPtr
	row := *rowPtr
	fmt.Println("Input filename: ", inputFilename)

	fileContents, err := ioutil.ReadFile(inputFilename)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(fileContents)
	fmt.Println("Length of data:", len(fileContents))
	numBlocks := len(fileContents) / 16
	remainBlocks := len(fileContents) % 16
	if len(fileContents) % 16 != 0 {
		fmt.Println("Warning: File size bytes isn't a multiple of 16. Will truncate some blocks at the end. Remaining:", remainBlocks)
	}
	fmt.Println("Will create", numBlocks, "8x8 blocks with", row, "rows")

	numRows := row
	totalBlocksRow := (numBlocks / numRows)
	width := 8 * totalBlocksRow
	height := 8 * (numBlocks / totalBlocksRow)

	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

	// color palette only has four colors
	colorArr := [4]color.RGBA{
		color.RGBA{255, 255, 255, 255},
		color.RGBA{127, 127, 127, 255},
		color.RGBA{191, 191, 191, 255},
		color.RGBA{63, 63, 63, 255},
	}


	for i := 0; i < len(fileContents); i += 16 {
		if i + 16 > len(fileContents) {
			fmt.Println("Not a full block, remaining bytes:", fileContents[i : ])
			break
		}

		for j := i; j < i + 16; j += 2 {
			value1 := fileContents[j]
			value2 := fileContents[j + 1]

			blockNum := i / 16
			blockNumX := blockNum % totalBlocksRow
			blockNumY := blockNum / totalBlocksRow
			xOffset := 8 * blockNumX
			yOffset := 8 * blockNumY
			y := (j - i) / 2

			for x := 0; x < 8; x++ {
				bit1 := ((value1 >> (7 - x)) & 1)
				bit2 := ((value2 >> (7 - x)) & 1)

				index := 0
				if bit1 == 0 && bit2 == 0 {
					index = 0
				} else if bit1 == 0 && bit2 == 1 {
					index = 1
				} else if bit1 == 1 && bit2 == 0 {
					index = 2
				} else if bit1 == 1 && bit2 == 1 {
					index = 3
				}

				img.Set(xOffset + x, yOffset + y, colorArr[index])
			}
		}

	}


	// Encode as PNG.
	outputFilename := inputFilename[0 : len(inputFilename) - len(".2bpp")] + ".png"
	f, _ := os.Create(outputFilename)
	png.Encode(f, img)

	fmt.Println("Written output to", outputFilename)
}
