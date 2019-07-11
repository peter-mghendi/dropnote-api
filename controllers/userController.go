package controllers

import (
	"dropnote-backend/models"
	u "dropnote-backend/utils"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// GetUser is the handler funcion for getting a user fom the database.
func GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value(UserKey) == nil {
		u.Respond(w, u.Message(false, "Please log in"))
		return
	}
	id := r.Context().Value(UserKey).(uuid.UUID)
	user := models.GetUser(App.DB, id)
	resp := u.Message(true, "success")
	resp["data"] = user
	u.Respond(w, resp)
}
