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

// args is the struct that defines the arguments you can give to the binary to change the behavior.
var args struct {
	Vips     bool   `default:"false"`
	DataPath string `default:"./data/"`
	Database string `default:"db.sqlite3"`
}

var s Service

// main is the function that is called when running go run . and is our entry to this program.
// it figures out the arguments that has been used to run this program and dependency injects them
// into the handler and service.
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

// ImageIdMustBeValid is a middleware that checks that the image ID given in the parameter id is valid (i.e. in the
// database).
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
