package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Articles []Article
	Username string
	FullName string
	Email    string
	SocialId string
	Provider string
	Avatar   string
	Role     bool `gorm:"default:0"`
}
