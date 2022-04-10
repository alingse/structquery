package tests_test

import (
	"encoding/json"
	"net/http"

	"github.com/alingse/structquery"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

var GormDB *gorm.DB
var queryDecoder = schema.NewDecoder()
var structQueryer = structquery.NewQueryer()

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

	err := queryDecoder.Decode(&query, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	db := GormDB.Model(&UserModel{})
	db, err = structQueryer.Where(db, query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var users []UserModel
	db.Find(&users)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
