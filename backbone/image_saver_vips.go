//go:build !windows

package main

import (
	"github.com/h2non/bimg"
	"io"
)

type ImageSaverVips struct {
}

func (saver ImageSaverVips) SaveImage(filename string, imageData []byte) error {
	imageData, _ := io.ReadAll(imageReader)
	imageData, _ = bimg.Resize(imageData, bimg.Options{Height: 400, Width: 400})
	bimg.Write("out.jpg", imageData)
	return nil
}
