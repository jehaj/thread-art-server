package main

import (
	"github.com/alexflint/go-arg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

	demoUserID := os.Getenv("DEMO_USER_ID")
	log.Println("Demo user id:", demoUserID)

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
	h := Handler{&s, imageSaver, demoUserID}
	InitializeLists()

	r.Get("/", h.GetIndex)
	r.Get("/favicon.svg",
		func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "dist/favicon.svg") })
	r.Handle("/assets/*", http.FileServer(http.Dir("dist/")))
	r.Post("/api/upload", h.UploadImage)
	r.Get("/api/user/{id}", h.GetUser)
	r.Route("/api/{id}", func(r chi.Router) {
		r.Use(ImageIdMustBeValid)
		r.Get("/in.png", h.GetImageOrig)
		r.Get("/out.png", h.GetImage)
		r.Get("/points", h.GetPoints)
	})
	r.Get("/*", h.GetIndex)

	checkDemoUser(demoUserID, &s, imageSaver)

	_ = http.ListenAndServe(":8080", r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func checkDemoUser(id string, service *Service, saver ImageSaver) {
	if len(id) == 0 {
		// the demo user should not be made
		return
	}
	var demoUser User
	service.DB.Preload("Images").First(&demoUser, "id = ?", id)
	if len(demoUser.Images) > 0 {
		// demo user has already been created
		return
	}
	// add the images
	service.DB.Create(&User{ID: id})
	imageFilepaths := []string{"demo1.png", "demo2.png", "demo3.png"}
	for _, imageFilepath := range imageFilepaths {
		imageID := uuid.New().String()
		imageReader, err := os.Open(imageFilepath)
		if err != nil {
			log.Println(err.Error())
			return
		}
		_ = os.Mkdir(filepath.Join(args.DataPath, imageID), 0750)
		saver.SaveImage(imageID, imageReader)

		err = service.AddImage(&Image{ID: imageID, UserID: id, Finished: false})
		if err != nil {
			return
		}
		s.AddImageToQueue(imageID)
	}
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
