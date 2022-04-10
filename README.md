# structquery
gorm query with struct

## Usage

```go
package main

import (
	"github.com/alingse/structquery"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u UserModel) TableName() string {
	return "users"
}

type UserModelQuery struct {
	Name  string `sq:"like"`
	Email string `sq:"eq"`
}

var DB = initGormDB()
var queryer = structquery.NewQueryer()

func main() {
	var q = UserModelQuery{
		Name:  "hello",
		Email: "a@b",
	}

	db := DB.Model(&User{})
	users := []User{}

	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
        tx, _ = queryer.Where(tx, q)
        return tx.Find(&users)
	})
    // SELECT * FROM `users` WHERE (`name` LIKE "%hello%" AND `email` = "a@b")
    fmt.Println(sql)
}
```

## Web Example

```go
package tests_test

import (
	"encoding/json"
	"net/http"

	"github.com/alingse/structquery"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

var structQueryer = structquery.NewQueryer()

var GormDB *gorm.DB
var queryDecoder = schema.NewDecoder()

func init() {
	queryDecoder.IgnoreUnknownKeys(true)
	queryDecoder.SetAliasTag("json")
}

type UserQuery struct {
	Name  string `json:"name"  sq:"like"`
	Email string `json:"email" sq:"eq"`
}

type UserModel struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u UserModel) TableName() string {
	return "users"
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var query UserQuery

    // decode from URL.Query()
	err := queryDecoder.Decode(&query, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	db := GormDB.Model(&UserModel{})

    // auto add Where with sq tag
	db, err = structQueryer.Where(db, query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var users []UserModel
	db.Find(&users)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
```

## TODO

1. support and/or with slice?
2. add structQueryer.Parse return with interface
3. more db test
