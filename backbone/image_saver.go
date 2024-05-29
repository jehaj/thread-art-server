package main

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io"
	"os"
)

type ImageSaver interface {
	SaveImage(filename string, imageReader io.Reader) error
}

type ImageSaverStd struct {
}

func (saver *ImageSaverStd) SaveImage(filename string, imageReader io.Reader) error {
	// Go can decode arbitrarily large images, therefore the
	// size needs to be within bounds [400, 1024]. The file limit
	// is already limited by MaxFileSize in handler.go.
	imageData, _ := io.ReadAll(imageReader)
	imageReader = bytes.NewReader(imageData)
	config, _, err := image.DecodeConfig(imageReader)
	if err != nil {
		return err
	}
	imageHeight := config.Height
	imageWidth := config.Width
	maxImageSize := 1024
	minImageSize := 400
	if min(imageWidth, imageHeight) > maxImageSize {
		return errors.New("image too large")
	} else if max(imageWidth, imageHeight) < minImageSize {
		return errors.New("image too small")
	}

	// Image is not out of bounds, try to decode it.
	imageReader = bytes.NewReader(imageData)
	decodedImage, _, err := image.Decode(imageReader)
	if err != nil {
		return err
	}
	// Create another image that is grayscale.
	grayImageSize := 400
	grayImage := image.NewGray(image.Rect(0, 0, grayImageSize, grayImageSize))
	// Only considers last pixel in section to resize the image down to 400x400.
	for py := 0; py < imageHeight; py++ {
		for px := 0; px < imageWidth; px++ {
			pColor := decodedImage.At(px, py)
			ny := grayImageSize * py / imageHeight
			nx := grayImageSize * px / imageWidth
			var _ uint32 = uint32(imageHeight / grayImageSize * imageWidth / grayImageSize)
			newColor := color.GrayModel.Convert(pColor).(color.Gray)
			grayImage.SetGray(nx, ny, newColor)
		}
	}

	// Write the grayscale image to disk
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	err = png.Encode(file, grayImage)
	if err != nil {
		return err
	}
	return nil
}
