package controllers

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Application is a copy of the app.App struct
type Application struct {
	Router *mux.Router
	DB     *gorm.DB
}

type key string

const UserKey key = "user"

// App is an instance of Application
var App *Application

func init() {
	App = &Application{}
}
