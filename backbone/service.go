package main

import (
	"gorm.io/gorm"
	"log"
)

type Service struct {
	*gorm.DB
	jobs chan string
}

// initialize initializes the database tables and will add images to the channel signal given.
func (s *Service) initialize(jobs chan string) {
	err := s.DB.AutoMigrate(&User{})
	if err != nil {
		log.Println(err.Error())
	}
	err = s.DB.AutoMigrate(&Image{})
	if err != nil {
		log.Println(err.Error())
	}
	s.jobs = jobs
}

// AddUserWithImage will add the user to the database. The user should contain an image.
func (s *Service) AddUserWithImage(user *User) error {
	s.DB.Create(user)
	return nil
}

// AddImage Only adds the image to the database.
func (s *Service) AddImage(image *Image) error {
	s.DB.Create(image)
	return nil
}

// ValidUserId Check that the user ID is valid.
func (s *Service) ValidUserId(userID string) bool {
	result := s.DB.First(&User{}, "id = ?", userID)
	return result.RowsAffected == 1
}

// AddImageToQueue add the given imageID to the channel signal. The image must be saved on the disk at
// [datapath]/[id]/[in.png] because the thread can start working immediately and expects it to be there.
func (s *Service) AddImageToQueue(imageID string) {
	s.jobs <- imageID
}

// ValidImageId check that the imageID exists in the database.
func (s *Service) ValidImageId(imageID string) bool {
	result := s.DB.First(&Image{}, "id = ?", imageID)
	return result.RowsAffected == 1
}

// GetUser Get the user with userID if exists.
func (s *Service) GetUser(userID string) (User, error) {
	var user User
	res := s.DB.Preload("Images").First(&user, "id = ?", userID)
	err := res.Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// ImageFinished checks if the image with imageID is marked as finished. The worker thread will mark it when done.
func (s *Service) ImageFinished(imageID string) {
	var image Image
	s.DB.First(&image, Image{ID: imageID})
	image.Finished = true
	s.DB.Save(&image)
}
