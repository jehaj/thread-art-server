package main

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	"math/rand/v2"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

type Handler struct {
	s          *Service
	imageSaver ImageSaver
}

const MaxFileSize int64 = 1024 * 1024 * 5

var colorList []string
var animalList []string

func InitializeLists() {
	colorFile, err1 := os.ReadFile("color-list.txt")
	animalFile, err2 := os.ReadFile("animal-list.txt")
	if err1 != nil || err2 != nil {
		log.Panicln("Error reading color or animal list.")
	}
	colorList = strings.Split(string(colorFile), "\n")
	animalList = strings.Split(string(animalFile), "\n")
}

func (h *Handler) GetIndex(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Hello World"))
}

func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	err, imageReader := getImageFromRequest(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	// get user ID
	authHeader := r.Header.Get("Authorization")
	var userID string
	userAlreadyExists := false
	if len(authHeader) < 4 {
		userID = getRandomUserID()
	} else {
		userID = strings.Split(authHeader, " ")[1]
		if !h.s.ValidUserId(userID) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("UserID is not valid"))
			return
		}
		userAlreadyExists = true
	}
	// get ID for the new image
	imageID := uuid.New().String()
	_ = os.Mkdir(filepath.Join(args.DataPath, imageID), 0750)
	err = h.imageSaver.SaveImage(imageID, imageReader)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if userAlreadyExists {
		err = h.s.AddImage(&Image{ID: imageID, UserID: userID, Finished: false})
	} else {
		err = h.s.AddUserWithImage(&User{ID: userID, Images: []Image{{ID: imageID, UserID: userID, Finished: false}}})
	}
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// tell worker to start
	h.s.AddImageToQueue(imageID)
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(userID))
}

// GetImage gets the image. The id has been checked to be valid.
func (h *Handler) GetImage(w http.ResponseWriter, r *http.Request) {
	imageID := chi.URLParam(r, "id")
	file, err := os.ReadFile(filepath.Join(args.DataPath, imageID, "out.png"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(file)
}

// GetImageOrig gets the image. The id has been checked to be valid.
func (h *Handler) GetImageOrig(w http.ResponseWriter, r *http.Request) {
	imageID := chi.URLParam(r, "id")
	file, err := os.ReadFile(filepath.Join(args.DataPath, imageID, "in.png"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(file)
}

type points struct {
	NumberOfPoints int
	PointIndex     []int
}

// GetPoints gets the points as json. The id has been checked to be valid.
func (h *Handler) GetPoints(w http.ResponseWriter, r *http.Request) {
	imageID := chi.URLParam(r, "id")
	file, err := os.ReadFile(filepath.Join(args.DataPath, imageID, "RESULT.txt"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	lines := strings.Split(string(file), "\n")
	n, _ := strconv.Atoi(lines[0])
	ps := strings.Split(lines[1], ",")
	p := make([]int, len(ps))
	for i, pc := range ps {
		atoi, _ := strconv.Atoi(pc)
		p[i] = atoi
	}
	pointsStruct := points{n, p}
	log.Println(pointsStruct)
	pointsString, err := json.Marshal(pointsStruct)
	log.Println(pointsString)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(pointsString)
}

// GetUser gets the user based on id in url. User has lists of images.
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	user, err := h.s.GetUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonUser)
}

// getRandomUserID returns a random ID with the form
// <number>-<color>-<animal>. It depends on the entries in
// the files animal-list.txt and color-list.txt
func getRandomUserID() string {
	randomNumber := rand.IntN(99) + 1
	randomIndex := rand.IntN(len(colorList))
	randomColor := colorList[randomIndex]
	randomIndex = rand.IntN(len(animalList))
	randomAnimal := animalList[randomIndex]
	combination := []string{strconv.Itoa(randomNumber), randomColor, randomAnimal}
	randomUserID := strings.Join(combination, "-")
	return randomUserID
}

func getImageFromRequest(w http.ResponseWriter, r *http.Request) (error, multipart.File) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize)
	err := r.ParseMultipartForm(MaxFileSize)
	if err != nil {
		return err, nil
	}
	imageReader, imageHeader, err := r.FormFile("image")
	if err != nil {
		return err, nil
	}
	contentType := imageHeader.Header.Get("Content-Type")
	contentTypeEnd := strings.Split(contentType, "/")[1]
	if !slices.Contains([]string{"jpeg", "jpg", "png"}, contentTypeEnd) {
		return errors.New("image format not supported"), nil
	}
	return nil, imageReader
}
