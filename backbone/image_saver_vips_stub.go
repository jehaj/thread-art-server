//go:build windows

package main

type ImageSaverVips struct {
}

func (saver *ImageSaverVips) SaveImage(filename string, imageData []byte) error {
	panic("vips is not supported on windows!")
}
