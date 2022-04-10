module github.com/alingse/structquery/tests

go 1.18

require (
	github.com/alingse/structquery v0.0.0-20220403113550-3b53d16e153e
	gorm.io/driver/sqlite v1.3.1
	gorm.io/gorm v1.23.4
)

require (
	github.com/gorilla/schema v1.2.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.9 // indirect
)

replace github.com/alingse/structquery => ../
