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
	"path/filepath"
)

type ImageSaver interface {
	SaveImage(filename string, imageReader io.Reader) error
}

type ImageSaverStd struct {
}

func (saver *ImageSaverStd) SaveImage(imageID string, imageReader io.Reader) error {
	// Go can decode arbitrarily large images, therefore the
	// size needs to be within bounds [400, 1024]. The file limit
	// is already limited by MaxFileSize in handler.go.
	imageData, _ := io.ReadAll(imageReader)
	imageHeight, imageWidth, err := checkBounds(imageReader, imageData)
	if err != nil {
		return err
	}

	grayImage, err := convertImageToGray(imageData, 400, imageHeight, imageWidth)
	if err != nil {
		return err
	}

	err = writeImage(filepath.Join(args.DataPath, imageID, "in.png"), grayImage)
	return err
}

func convertImageToGray(imageData []byte, grayImageSize, imageHeight int, imageWidth int) (*image.Gray, error) {
	// Image is not out of bounds, try to decode it.
	imageReader := bytes.NewReader(imageData)
	decodedImage, _, err := image.Decode(imageReader)
	if err != nil {
		return nil, err
	}
	// Create another image that is grayscale.
	grayImage := image.NewGray(image.Rect(0, 0, grayImageSize, grayImageSize))
	// Only considers last pixel in section to resize the image down to 400x400.
	for py := 0; py < imageHeight; py++ {
		for px := 0; px < imageWidth; px++ {
			pColor := decodedImage.At(px, py)
			ny := grayImageSize * py / imageHeight
			nx := grayImageSize * px / imageWidth
			newColor := color.GrayModel.Convert(pColor).(color.Gray)
			grayImage.SetGray(nx, ny, newColor)
		}
	}
	return grayImage, nil
}

func writeImage(filepath string, grayImage *image.Gray) error {
	// Write the grayscale image to disk
	file, err := os.Create(filepath)
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

func checkBounds(imageReader io.Reader, imageData []byte) (int, int, error) {
	imageReader = bytes.NewReader(imageData)
	config, _, err := image.DecodeConfig(imageReader)
	if err != nil {
		return 0, 0, err
	}
	imageHeight := config.Height
	imageWidth := config.Width
	maxImageSize := 1024
	minImageSize := 400
	if min(imageWidth, imageHeight) > maxImageSize {
		return 0, 0, errors.New("image too large")
	} else if max(imageWidth, imageHeight) < minImageSize {
		return 0, 0, errors.New("image too small")
	}
	return imageHeight, imageWidth, err
}
