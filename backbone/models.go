package main

import (
	"gorm.io/gorm"
)

type User struct {
	ID string `gorm:"primaryKey"`
	gorm.Model
	Images []Image
}

type Image struct {
	ID string `gorm:"primaryKey"`
	gorm.Model
	UserID   string
	Finished bool
}
