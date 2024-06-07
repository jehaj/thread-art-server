package main

import (
	"github.com/alexflint/go-arg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"os"
)

var args struct {
	Vips     bool   `default:"false"`
	DataPath string `default:"./data/"`
	Database string `default:"db.sqlite3"`
}

var s Service

func main() {
	arg.MustParse(&args)
	_ = os.Mkdir(args.DataPath, 0750)
	r := chi.NewRouter()
	c := cors.AllowAll()
	r.Use(c.Handler)
	db, err := gorm.Open(sqlite.Open(args.Database), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	s = Service{db, nil}
	s.initialize(workerPool(&s))
	var imageSaver ImageSaver = &ImageSaverStd{}
	if args.Vips {
		imageSaver = &ImageSaverVips{}
	}
	h := Handler{&s, imageSaver}
	InitializeLists()
	r.Get("/", h.GetIndex)
	r.Post("/api/upload", h.UploadImage)
	r.Get("/api/user/{id}", h.GetUser)
	r.Route("/api/{id}", func(r chi.Router) {
		r.Use(ImageIdMustBeValid)
		r.Get("/in.png", h.GetImageOrig)
		r.Get("/out.png", h.GetImage)
		r.Get("/points", h.GetPoints)
	})
	_ = http.ListenAndServe("localhost:8080", r)
}

func ImageIdMustBeValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		imageID := chi.URLParam(r, "id")
		isValid := s.ValidImageId(imageID)
		if !isValid {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
