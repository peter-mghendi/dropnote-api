package app

import (
	"fmt"

	"dropnote-backend/models"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
	Port   string
}

type URI struct {
	Host, User, Name, Pass, Type string
}

var err error

func (a *App) initDB(u URI) {
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", u.Host, u.User, u.Name, u.Pass)
	fmt.Sprintf(dbURI)
	a.DB, err = gorm.Open(u.Type, dbURI)
	if err != nil {
		fmt.Print(err)
	}

	a.DB.Debug().AutoMigrate(&models.User{}, &models.Note{})
}

func (a *App) initRoutes() {
	a.Router = mux.NewRouter()
	fmt.Println("Initializing routes")
}

func (a *App) initVars() {
	fmt.Println("Exporting variables")
}

// Init sets up database and routes
func (a *App) Init(u URI) {
	a.initDB(u)
	a.initRoutes()
	a.initVars()
}

// Run serves the API on a specified port
func (a *App) Run() {
	fmt.Printf("Running on port %s", a.Port)
}
