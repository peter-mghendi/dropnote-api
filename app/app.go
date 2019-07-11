package app

import (
	"fmt"
	"log"
	"net/http"

	"dropnote-backend/controllers"
	"dropnote-backend/models"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// App holds details about router, database and port
type App struct {
	Router *mux.Router
	DB     *gorm.DB
	Port   string
}

// URI holds database connection credentials
type URI struct {
	Host, User, Name, Pass, Type string
}

var err error

func (a *App) initDB(u URI) {
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", u.Host, u.User, u.Name, u.Pass)
	_ = fmt.Sprintf(dbURI)
	a.DB, err = gorm.Open(u.Type, dbURI)
	if err != nil {
		fmt.Print(err)
	}

	a.DB.Debug().AutoMigrate(&models.User{}, &models.Note{})
}

func (a *App) initRoutes() {
	a.Router = mux.NewRouter()
	a.Router.Use(JwtAuthentication)
	a.Router.HandleFunc(createUser, controllers.CreateUser).Methods(post)
	a.Router.HandleFunc(authUser, controllers.AuthUser).Methods(post)
	a.Router.HandleFunc(createNote, controllers.CreateNote).Methods(post)
	a.Router.HandleFunc(getNote, controllers.GetNote).Methods(get)
	a.Router.HandleFunc(getNotes, controllers.GetNotes).Methods(get)
	a.Router.HandleFunc(getUser, controllers.GetUser).Methods(get)
	a.Router.HandleFunc(deleteUser, controllers.DeleteUser).Methods(delete)
	a.Router.HandleFunc(getUserNotes, controllers.GetUserNotes).Methods(get)
}

func (a *App) initVars() {
	controllers.App.Router, controllers.App.DB = a.Router, a.DB
}

// Init sets up database and routes
func (a *App) Init(u URI) {
	a.initDB(u)
	a.initRoutes()
	a.initVars()
}

// Run serves the API on a specified port
func (a *App) Run() {
	fmt.Printf("Serving on localhost:%v\n", a.Port)
	log.Fatal(http.ListenAndServe(":"+a.Port, a.Router))
}
