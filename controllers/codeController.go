package controllers

import (
	"dropnote-backend/models"
	u "dropnote-backend/utils"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// GenerateCode is the handler function for creating a new reset token
func GenerateCode(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user, err := uuid.FromString(params["user"])
	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}
	code, err := models.New(App.DB, models.Actions["reset"], user)
	if err != nil {
		u.Respond(w, u.Message(false, err.Error()))
		return
	}
	if !emailTokenToUser(user, code) {
		u.Respond(w, u.Message(false, "Unable to send email"))
		return
	}
	resp := u.Message(true, "success")
	u.Respond(w, resp)
}
