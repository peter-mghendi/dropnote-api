package utils

import (
	"encoding/json"
	"net/http"
)

// Message builds a response in the form map[string]interface{}
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// Respond adds data in the form map[string]interface{} to http header
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
