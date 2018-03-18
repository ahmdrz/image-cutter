package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path"

	"github.com/disintegration/imaging"
)

func sliceImageV(src image.Image, imageSize image.Point, slides int) *image.NRGBA {
	widthSize := imageSize.X / slides

	outputImage := imaging.New(imageSize.X, imageSize.Y, color.NRGBA{0, 0, 0, 0})

	outputImage1 := imaging.New(imageSize.X/2, imageSize.Y, color.NRGBA{0, 0, 0, 0})
	outputImage2 := imaging.New(imageSize.X/2, imageSize.Y, color.NRGBA{0, 0, 0, 0})

	j := 0
	for i := 0; i < slides; i += 2 {
		sliceA := imaging.Crop(src, image.Rect(i*widthSize, 0, (i+1)*widthSize, imageSize.Y))
		sliceB := imaging.Crop(src, image.Rect((i+1)*widthSize, 0, (i+2)*widthSize, imageSize.Y))
		outputImage1 = imaging.Paste(outputImage1, sliceA, image.Pt(j*widthSize, 0))
		outputImage2 = imaging.Paste(outputImage2, sliceB, image.Pt(j*widthSize, 0))
		j += 1
	}

	outputImage = imaging.Paste(outputImage, outputImage1, image.Pt(0, 0))
	outputImage = imaging.Paste(outputImage, outputImage2, image.Pt(imageSize.X/2, 0))
	return outputImage
}

func sliceImageH(src image.Image, imageSize image.Point, slides int) *image.NRGBA {
	heigthSize := imageSize.Y / slides

	outputImage := imaging.New(imageSize.X, imageSize.Y, color.NRGBA{0, 0, 0, 0})

	outputImage1 := imaging.New(imageSize.X, imageSize.Y/2, color.NRGBA{0, 0, 0, 0})
	outputImage2 := imaging.New(imageSize.X, imageSize.Y/2, color.NRGBA{0, 0, 0, 0})

	j := 0
	for i := 0; i < slides; i += 2 {
		sliceA := imaging.Crop(src, image.Rect(0, i*heigthSize, imageSize.X, (i+1)*heigthSize))
		sliceB := imaging.Crop(src, image.Rect(0, (i+1)*heigthSize, imageSize.X, (i+2)*heigthSize))
		outputImage1 = imaging.Paste(outputImage1, sliceA, image.Pt(0, j*heigthSize))
		outputImage2 = imaging.Paste(outputImage2, sliceB, image.Pt(0, j*heigthSize))
		j += 1
	}

	outputImage = imaging.Paste(outputImage, outputImage1, image.Pt(0, 0))
	outputImage = imaging.Paste(outputImage, outputImage2, image.Pt(0, imageSize.Y/2))
	return outputImage
}

func main() {
	var (
		inputPath  string
		slidesX    int
		slidesY    int
		outputPath string
	)
	flag.StringVar(&inputPath, "input", "", "input image")
	flag.StringVar(&outputPath, "output", "output", "output directory")
	flag.IntVar(&slidesX, "x", 20, "number of slides as x (high values => high quality)")
	flag.IntVar(&slidesY, "y", 20, "number of slides as y (high values => high quality)")
	flag.Parse()

	if slidesX*slidesY == 0 {
		log.Fatalln("zero x or y value")
	}

	err := os.MkdirAll(outputPath, 0755)
	if err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	src, err := imaging.Open(inputPath)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	// Resize image to remove unnecessary bounds.
	imageSize := src.Bounds().Size()
	src = imaging.Resize(src, (imageSize.X/slidesX)*slidesX, (imageSize.Y/slidesY)*slidesY, imaging.Lanczos)
	imageSize = src.Bounds().Size()

	outputImage := sliceImageV(sliceImageH(src, imageSize, slidesY), imageSize, slidesX)
	err = imaging.Save(outputImage, path.Join(outputPath, "output.jpg"))
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

	widthSize := imageSize.X / 2
	heigthSize := imageSize.Y / 2

	images := []image.Image{
		imaging.CropAnchor(outputImage, widthSize, heigthSize, imaging.TopLeft),
		imaging.CropAnchor(outputImage, widthSize, heigthSize, imaging.TopRight),
		imaging.CropAnchor(outputImage, widthSize, heigthSize, imaging.BottomLeft),
		imaging.CropAnchor(outputImage, widthSize, heigthSize, imaging.BottomRight),
	}

	for i, image := range images {
		err = imaging.Save(image, path.Join(outputPath, fmt.Sprintf("output-%d.jpg", i)))
		if err != nil {
			log.Fatalf("failed to save image: %v", err)
		}
	}
}
