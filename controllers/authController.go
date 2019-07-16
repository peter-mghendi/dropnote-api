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

	resp := make(map[string]interface{})
	if data, err := user.Create(App.DB); err != nil {
		resp = u.Message(false, err.Error())
	} else {
		resp = data
	}
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

	resp := make(map[string]interface{})
	if data, err := models.Login(App.DB, user.Mail, user.Pass); err != nil {
		resp = u.Message(false, err.Error())
	} else {
		resp = data
	}
	u.Respond(w, resp)
}
