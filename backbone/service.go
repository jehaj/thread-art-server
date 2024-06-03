package main

import (
	"gorm.io/gorm"
	"log"
)

type Service struct {
	*gorm.DB
	jobs chan string
}

func (s *Service) initialize() {
	err := s.DB.AutoMigrate(&User{})
	if err != nil {
		log.Println(err.Error())
	}
	err = s.DB.AutoMigrate(&Image{})
	if err != nil {
		log.Println(err.Error())
	}
}

func (s *Service) AddUserWithImage(user *User) error {
	s.DB.Create(user)
	return nil
}

func (s *Service) ValidUserId(userID string) bool {
	result := s.DB.First(&User{}, userID)
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
