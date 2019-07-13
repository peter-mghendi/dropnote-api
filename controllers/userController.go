package controllers

import (
	"dropnote-backend/models"
	u "dropnote-backend/utils"
	"encoding/json"
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

// UpdateUser is the handler function for editing a user in the database
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(UserKey).(uuid.UUID)
	user, temp := models.GetUser(App.DB, id), &models.User{}

	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request payload"))
		return
	}

	if tempUser := models.GetUserByMail(App.DB, temp.Mail); tempUser != nil {
		u.Respond(w, u.Message(false, "That email address is already in use"))
		return
	}

	user.Name = temp.Name
	user.Mail = temp.Mail

	err = models.UpdateUser(App.DB, user)
	resp := u.Message(true, "success")
	if err != nil {
		resp = u.Message(false, "failed")
		resp["error"] = err
	}
	u.Respond(w, resp)
}

// UpdatePassword is the handler function for changing a user's password
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(UserKey).(uuid.UUID)
	user := models.GetUser(App.DB, id)
	payload := struct {
		Current string `json:"current"`
		Updated string `json:"updated"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request payload"))
		return
	}
	if payload.Current == "" || payload.Updated == "" {
		u.Respond(w, u.Message(false, "One of the fields is empty"))
		return
	}

	mail, pass := user.Mail, payload.Current

	result := models.Login(App.DB, mail, pass)
	ok := result["status"].(bool)
	if !ok {
		message := result["message"].(string)
		u.Respond(w, u.Message(false, message))
		return
	}

	user.Pass = payload.Updated
	err = models.UpdatePassword(App.DB, user)
	resp := u.Message(true, "success")
	if err != nil {
		resp = u.Message(false, "failed")
		resp["error"] = err
	}
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

// UpdateUserNote is the handler function for editing a note created by current user in the database
func UpdateUserNote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := r.Context().Value(UserKey).(uuid.UUID)
	id, err := uuid.FromString(params["id"])
	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}

	note := models.GetNote(App.DB, id)
	if note == nil {
		u.Respond(w, u.Message(false, "That note does not exist"))
		return
	}

	if !uuid.Equal(note.Creator, user) {
		u.Respond(w, u.Message(false, "You are not authorized to modify that record"))
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&note); err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request payload"))
		return
	}

	err = models.UpdateNote(App.DB, note)
	resp := u.Message(true, "success")
	if err != nil {
		// TODO Better error handling
		resp = u.Message(false, "failed")
	}
	u.Respond(w, resp)
}

// DeleteUserNote is the handler function for removing a note created by current user from the database
// TODO Auth
func DeleteUserNote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	note, err := uuid.FromString(params["id"])
	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}

	user := r.Context().Value(UserKey).(uuid.UUID)
	err = models.DeleteNoteFor(App.DB, note, user)
	resp := u.Message(true, "success")
	if err != nil {
		// TODO Better error handling
		resp = u.Message(false, "failed")
	}
	u.Respond(w, resp)
}
