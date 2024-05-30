package main

import (
	"gorm.io/gorm"
	"log"
)

type Service struct {
	*gorm.DB
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
