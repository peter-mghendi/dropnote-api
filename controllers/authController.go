package controllers

import (
	"encoding/json"
	"net/http"
	"dropnote-backend/models"
	u "dropnote-backend/utils"
)

// CreateUser is the handler function for adding a new account into the databsse
func CreateUser (w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := user.Create(App.DB)
	u.Respond(w, resp)
}