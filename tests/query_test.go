package tests_test

import (
	"testing"

	"github.com/alingse/structquery"
	"gorm.io/gorm"
)

type UserQuery struct {
	Name  string `sq:"like"`
	Email string `sq:"eq"`
}

var queryer = structquery.NewQueryer()

func TestUserQuery(t *testing.T) {
	var q = UserQuery{
		Name:  "hello",
		Email: "a@b",
	}
	cond, err := queryer.And(q)
	if err != nil {
		t.Fail()
	}

	db := DB.Model(&User{})
	users := []User{}
	expectSQL := "SELECT * FROM `users` WHERE (`name` LIKE \"%hello%\" AND `email` = \"a@b\")"

	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Where(cond).Find(&users)
	})
	if sql != expectSQL {
		t.Errorf("sql not equal, got %s", sql)
	}
}
