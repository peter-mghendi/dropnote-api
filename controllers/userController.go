package controllers

import (
	"dropnote-backend/models"
	u "dropnote-backend/utils"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// GetUser is the handler funcion for getting a user fom the database.
func GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(UserKey).(uuid.UUID)
	user := models.GetUser(App.DB, id)
	resp := u.Message(true, "success")
	resp["data"] = user
	u.Respond(w, resp)
}

// DeleteUser is the handler function for removing a user from the database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(UserKey).(uuid.UUID)
	user := models.GetUser(App.DB, id)
	err := models.DeleteUser(App.DB, user)
	resp := u.Message(true, "success")
	if err != nil {
		// TODO Better error handling
		resp = u.Message(false, "failed")
	}
	u.Respond(w, resp)
}

// GetUserNotes is the handler funcion for getting a user fom the database.
func GetUserNotes(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(UserKey).(uuid.UUID)
	user := models.GetNotesFor(App.DB, id)
	resp := u.Message(true, "success")
	resp["data"] = user
	u.Respond(w, resp)
}
