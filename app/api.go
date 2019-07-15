package app

const (
	get, post, put, delete = "GET", "POST", "PUT", "DELETE"
	createUser             = "/api/user/new"
	authUser               = "/api/user/login"
	createNote             = "/api/note/new"
	getNote                = "/api/note/{id}"
	getNotes               = "/api/notes"
	getUser                = "/api/me"
	updateUser             = "/api/me/update"
	updatePassword         = "/api/me/change"
	deleteUser             = "/api/me/delete"
	getUserNotes           = "/api/me/notes"
	updateUserNote         = "/api/me/notes/update/{id}"
	deleteUserNote         = "/api/me/notes/delete/{id}"
	generateCode           = "/api/user/{user}/action/reset"
	executeCode            = "/api/user/{user}/action/{code}"
)
