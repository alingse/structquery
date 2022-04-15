package tests_test

import (
	"encoding/json"
	"net/http"

	"github.com/alingse/structquery"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

// DB -------------------------------------------------------------------------
var GormDB *gorm.DB

type UserModel struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u UserModel) TableName() string {
	return "users"
}

// HTTP -----------------------------------------------------------------------
var queryDecoder = schema.NewDecoder()

func init() {
	queryDecoder.IgnoreUnknownKeys(true)
	queryDecoder.SetAliasTag("json")
}

// Query ----------------------------------------------------------------------

type UserQuery struct {
	Name  string `json:"name"  sq:"like"`
	Email string `json:"email" sq:"eq"`
}

var structQueryer = structquery.NewQueryer()

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var query UserQuery

	// decode from URL
	err := queryDecoder.Decode(&query, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := GormDB.Model(&UserModel{})
	// bind queryer
	db, err = structQueryer.Where(db, query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var users []UserModel
	db.Find(&users)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}
