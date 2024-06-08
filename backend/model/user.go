package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Jwt      string `json:"jwt"`
}

type Register struct {
}
