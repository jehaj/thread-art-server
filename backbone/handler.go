package main

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"math/rand/v2"
	"mime/multipart"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
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
		log.Panicln(err1.Error(), err2.Error())
	}
	colorList = strings.Split(string(colorFile), "\n")
	animalList = strings.Split(string(animalFile), "\n")
}

func (h *Handler) GetIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	err, imageReader := getImageFromRequest(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// get user ID
	userIDCookie, err := r.Cookie("userID")
	var userID string
	if err != nil {
		userID = getRandomUserID()
	} else {
		userID = userIDCookie.Value
		if !h.s.ValidUserId(userID) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("UserID is not valid"))
			return
		}
	}
	// get ID for the new image
	imageID := uuid.New().String()
	err = h.imageSaver.SaveImage("upload.png", imageReader)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = h.s.AddUserWithImage(&User{userID, []Image{{imageID, userID, time.Now(), false}}})
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// tell worker to start
	h.s.AddImageToQueue(imageID)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(userID))
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
	return strings.Join([]string{strconv.Itoa(randomNumber), randomColor, randomAnimal}, "-")
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
