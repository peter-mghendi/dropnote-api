package main

import (
	"os"

	a "github.com/l3njo/dropnote-api/app"
	c "github.com/l3njo/dropnote-api/controllers"

	"github.com/joho/godotenv"

	_ "github.com/jinzhu/gorm/dialects/postgres" // init postgresql drivers
)

var app a.App
var err error
var uri a.URI

func init() {
	e := godotenv.Load()
	c.Handle(e)
	uri = a.URI{
		Host: os.Getenv("db_host"),
		User: os.Getenv("db_user"),
		Name: os.Getenv("db_name"),
		Pass: os.Getenv("db_pass"),
	}

	app = a.App{}
	app.Port = os.Getenv("PORT")
}

func main() {
	app.Init(uri)
	app.Run()
}
