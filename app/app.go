package app

import (
	"fmt"

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

func (a *App) initDB(u URI) {
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", u.Host, u.User, u.Name, u.Pass)
	fmt.Sprintf(dbURI)
}

func (a *App) initRoutes() {
	a.Router = mux.NewRouter()
	fmt.Println("Initializing routes")
}

func (a *App) initVars() {
	fmt.Println("Exporting variables")
}

func (a *App) Init(u URI) {
	a.initDB(u)
	a.initRoutes()
	a.initVars()
}

func (a *App) Run() {
	fmt.Printf("Running on port %s", a.Port)
}
