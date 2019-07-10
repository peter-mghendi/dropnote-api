package main

import (
	"fmt"
	"os"

	a "dropnote-backend/app"

	"github.com/joho/godotenv"

	_ "github.com/jinzhu/gorm/dialects/postgres" // init postgresql drivers
)

var app a.App
var uri a.URI

var err error
var port string

func init() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	uri = a.URI{}
	uri.Host = os.Getenv("db_host")
	uri.User = os.Getenv("db_user")
	uri.Name = os.Getenv("db_name")
	uri.Pass = os.Getenv("db_pass")
	uri.Type = os.Getenv("db_type")

	app = a.App{}
	app.Port = os.Getenv("PORT")
	if app.Port == "" {
		app.Port = "8000"
	}
}

func main() {
	app.Init(uri)
	app.Run()
}