package controllers

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Application is a copy of the app.App struct
type Application struct {
	Router *mux.Router
	DB     *gorm.DB
}

type key string

// UserKey is a constant context key
const UserKey key = "user"

// App is an instance of Application
var App *Application

func init() {
	App = &Application{}
}

func emailTokenToUser(user, code uuid.UUID) bool {
	// TODO
	fmt.Printf("http://localhost:8000/user/%s/action/%s\n", user.String(), code.String())
	return true
}
