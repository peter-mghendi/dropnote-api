package app

const (
	get, post, put, delete = "GET", "POST", "PUT", "DELETE"
	createUser             = "/api/users/new"
	authUser               = "/api/users/login"
	createNote             = "/api/notes/new"
	getNote                = "/api/notes/{id}" 
	getNotes               = "/api/notes" //
	getUser                = "/api/me"
	updateUser             = "/api/me/update"
	updatePassword         = "/api/me/change"
	deleteUser             = "/api/me/delete"
	getUserNotes           = "/api/me/notes"
	updateUserNote         = "/api/me/notes/update/{id}"
	toggleUserNote         = "/api/me/notes/toggle/{id}"
	deleteUserNote         = "/api/me/notes/delete/{id}"
	generateCode           = "/api/users/actions/reset"
	executeCode            = "/api/users/{user}/action/{code}"
)
