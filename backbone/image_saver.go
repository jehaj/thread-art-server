package main

type ImageSaver interface {
	SaveImage(filename string, imageData []byte) error
}

type ImageSaverStd struct {
}

func (saver *ImageSaverStd) SaveImage(filename string, imageData []byte) error {
	panic("not implemented")
	return nil
}
