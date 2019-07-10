package controllers

import (
	"dropnote-backend/models"
	u "dropnote-backend/utils"
	"encoding/json"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// CreateNote is the handler function for adding a note to the database
func CreateNote(w http.ResponseWriter, r *http.Request) {
	var id string
	if r.Context().Value("user") == nil {
		id = ""
	} else {
		id = (r.Context().Value("user").(uuid.UUID)).String()
	}
	user := uuid.FromStringOrNil(id)
	note := &models.Note{}
	err := json.NewDecoder(r.Body).Decode(note)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
	note.Creator = user
	resp := note.Create(App.DB)
	u.Respond(w, resp)
}
