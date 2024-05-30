package main

import (
	"github.com/alexflint/go-arg"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var args struct {
	Vips bool
}

func main() {
	arg.MustParse(&args)
	r := chi.NewRouter()
	db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	s := Service{db}
	s.initialize()
	db.Create(User{"dorthe", []Image{}})
	db.Create(Image{"dorthe", "dorthe", time.Now(), false})
	var imageSaver ImageSaver = &ImageSaverStd{}
	if args.Vips {
		imageSaver = &ImageSaverVips{}
	}
	h := Handler{&s, imageSaver}
	r.Get("/", h.GetIndex)
	r.Post("/api/upload", h.UploadImage)
	http.ListenAndServe("localhost:8080", r)
}
