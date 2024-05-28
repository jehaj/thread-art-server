package main

import (
	"github.com/alexflint/go-arg"
	"github.com/go-chi/chi/v5"
	"net/http"
)

var args struct {
	Vips bool
}

func main() {
	arg.MustParse(&args)
	r := chi.NewRouter()
	s := Service{"db.sqlite3"}
	var imageSaver ImageSaver = &ImageSaverStd{}
	if args.Vips {
		imageSaver = &ImageSaverVips{}
	}
	h := Handler{&s, imageSaver}
	r.Get("/", h.GetIndex)
	r.Post("/api/upload", h.UploadImage)
	http.ListenAndServe("localhost:8080", r)
}
