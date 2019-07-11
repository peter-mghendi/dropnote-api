package controllers

import (
	"dropnote-backend/models"
	u "dropnote-backend/utils"
	"net/http"

	"github.com/gorilla/mux"
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

// GetUserNotes is the handler funcion for getting notes created by current user
func GetUserNotes(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(UserKey).(uuid.UUID)
	user := models.GetNotesFor(App.DB, id)
	resp := u.Message(true, "success")
	resp["data"] = user
	u.Respond(w, resp)
}

// DeleteUserNote is the handler function for removing a note created by current user from the database
func DeleteUserNote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	note, err := uuid.FromString(params["id"])
	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}

	user := r.Context().Value(UserKey).(uuid.UUID)
	err = models.DeleteNoteFor(App.DB, user, note)
	resp := u.Message(true, "success")
	if err != nil {
		// TODO Better error handling
		resp = u.Message(false, "failed")
	}
	u.Respond(w, resp)
}
