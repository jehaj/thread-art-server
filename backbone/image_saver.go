package main

import (
	"errors"
	"log"
)

type ImageSaver interface {
	SaveImage(filename string, imageData []byte) error
}

type ImageSaverStd struct {
}

func (saver *ImageSaverStd) SaveImage(filename string, imageData []byte) error {
	log.Println("ImageSaverStd not implemented")
	return errors.New("ImageSaverStd not implemented")
}
