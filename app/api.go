package app

const (
	get, post, put, delete = "GET", "POST", "PUT", "DELETE"
	createUser             = "/api/user/new"
	authUser               = "/api/user/login"
	createNote             = "/api/note/new"
	getNote                = "/api/note/{id}"
	getNotes               = "/api/notes"
	getUser                = "/api/me"
	deleteUser             = "/api/me/delete"
	getUserNotes           = "/api/me/notes"
)
