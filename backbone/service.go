package main

import (
	"gorm.io/gorm"
	"log"
)

type Service struct {
	*gorm.DB
	jobs chan string
}

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

func (s *Service) AddUserWithImage(user *User) error {
	s.DB.Create(user)
	return nil
}

func (s *Service) AddImage(image *Image) error {
	s.DB.Create(image)
	return nil
}

func (s *Service) ValidUserId(userID string) bool {
	result := s.DB.First(&User{}, "id = ?", userID)
	return result.RowsAffected == 1
}

func (s *Service) AddImageToQueue(imageID string) {
	s.jobs <- imageID
}

func (s *Service) ValidImageId(imageID string) bool {
	result := s.DB.First(&Image{}, "id = ?", imageID)
	return result.RowsAffected == 1
}

func (s *Service) GetUser(userID string) (User, error) {
	var user User
	res := s.DB.Preload("Images").First(&user, "id = ?", userID)
	err := res.Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *Service) ImageFinished(imageID string) {
	var image Image
	s.DB.First(&image, Image{ID: imageID})
	image.Finished = true
	s.DB.Save(&image)
}
