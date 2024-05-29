package main

import (
	"io"
	"net/http"
	"slices"
	"strings"
)

type Handler struct {
	s          *Service
	imageSaver ImageSaver
}

func (h *Handler) GetIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

const MaxFileSize int64 = 1024 * 1024 * 5

func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize)
	err := r.ParseMultipartForm(MaxFileSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	imageReader, imageHeader, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	contentType := imageHeader.Header.Get("Content-Type")
	contentTypeEnd := strings.Split(contentType, "/")[1]
	if !slices.Contains([]string{"jpeg", "jpg", "png"}, contentTypeEnd) {
		w.WriteHeader(http.StatusBadRequest)
	}
	imageData, _ := io.ReadAll(imageReader)
	err = h.imageSaver.SaveImage("upload.png", imageData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
