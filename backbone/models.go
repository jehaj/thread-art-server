package main

import "time"

type User struct {
	UserID string  `gorm:"primaryKey"`
	Images []Image `gorm:"foreignKey:UserID"`
}

type Image struct {
	ImageID   string `gorm:"primaryKey"`
	UserID    string
	Timestamp time.Time
	Finished  bool
}
