//go:build windows

package main

import "io"

type ImageSaverVips struct {
}

func (saver *ImageSaverVips) SaveImage(filename string, imageReader io.Reader) error {
	panic("vips is not supported on windows!")
}
