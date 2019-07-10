package controllers

import (
	"dropnote-backend/models"
	u "dropnote-backend/utils"
	"encoding/json"
	"net/http"
)

// CreateUser is the handler function for adding a new account into the databsse
func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := user.Create(App.DB)
	u.Respond(w, resp)
}

// AuthUser is the handler function for authorizing user login
func AuthUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(App.DB, user.Mail, user.Pass)
	u.Respond(w, resp)
}
