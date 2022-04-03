package tests_test

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}
}

type User struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
}

func (User) TableName() string {
	return "users"
}
