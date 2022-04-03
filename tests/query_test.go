package tests

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

var DB *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db
}

func TestDB(t *testing.T) {
	if DB == nil {
		t.Fail()
	}
}
