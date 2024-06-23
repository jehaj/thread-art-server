//go:build windows || linux

package main

import "io"

type ImageSaverVips struct {
}

func (saver *ImageSaverVips) SaveImage(filename string, imageReader io.Reader) error {
	panic("vips is not yet supported on windows!")
}
