package main

import (
	"image"
	_ "image/png"
	"log"
	"math"
	"os"
)

type settings struct {
	screenHeight int
	screenWidth  int
	threshold    int
}

var config settings = settings{
	screenHeight: 32,
	screenWidth:  128,
	threshold:    128,
}

func vertical1bit(imgData image.Image, width, height int, invert bool) [512]uint8 {
	outputArray := [512]uint8{}
	outputIndex := 0

	for p := 0; float64(p) < math.Ceil(float64(config.screenHeight)/8); p++ {
		for x := 0; x < config.screenWidth; x++ {
			byteIndex := 7
			var number uint8 = 0

			for y := 7; y >= 0; y-- {
				r, g, b, _ := imgData.At(x, p*8+y).RGBA()

				if invert {
					r = 0xFFFF - r
					g = 0xFFFF - g
					b = 0xFFFF - b
				}

				if int((r+g+b)/3) > config.threshold {
					number++
					// number += uint8(math.Pow(2, float64(byteIndex)))
				}
				if y > 0 {
					number = number << 1
				}
				// fmt.Print(number, " ")

				byteIndex--
			}
			outputArray[outputIndex] = number
			outputIndex++
		}
	}

	return outputArray

}

func readImages(path string, invert bool) [][512]uint8 {
	// Read images from folder at path

	// Count number of images
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	outputArray := make([][512]uint8, len(files))

	// Loop through images
	for i := 0; i < len(files); i++ {
		file := files[i]
		// fmt.Println(file.Name())
		// Read image
		imgFile, err := os.Open(path + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		defer imgFile.Close()
		img, _, err := image.Decode(imgFile)
		if err != nil {
			log.Fatal(err)
		}

		// Convert image to 1 bit
		bounds := img.Bounds()
		values := vertical1bit(img, bounds.Max.X, bounds.Max.Y, invert)

		outputArray[i] = values

		// for i := 0; i < len(values); i++ {
		// 	fmt.Printf("0x%02x, ", values[i])
		// 	if i%16 == 15 {
		// 		fmt.Println()
		// 	}
		// }
		// fmt.Println(values)
	}

	return outputArray
}
