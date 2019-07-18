package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	rand "github.com/l3njo/play/misc"
	uuid "github.com/satori/go.uuid"
)

type content struct {
	Content template.HTML
}

type dataset struct {
	Status, Message, Error string
}

// HACK for testing only
var sets = make(map[uuid.UUID]dataset)

func DoReset(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := uuid.FromStringOrNil(params["user"])
	code := uuid.FromStringOrNil(params["code"])

	if user == uuid.Nil || code == uuid.Nil {
		uri := fmt.Sprintf("/api/forms/result/%s", uuid.Nil.String())
		http.Redirect(w, r, uri, http.StatusFound)
		return
	}
	// TODO additional validation, check if user and code exist?

	fmt.Println("method", r.Method)

	if r.Method == "GET" {
		bytes, _ := ioutil.ReadFile("templates/auth.html.tmpl")
		body := template.HTML(bytes)
		data := content{
			Content: body,
		}
		t, _ := template.ParseFiles("templates/base.html.tmpl")
		t.Execute(w, data)
	} else {
		// TODO Update password, get error if any
		r.ParseForm()
		fmt.Println(r.Form["password"])
		title := "Success"
		ok, str, err := genRandomStatus()
		if !ok {
			title = "Failed"
		}

		// HACK Save struct
		id, _ := uuid.NewV4()
		message := dataset{
			Status:  title,
			Message: str,
			Error:   err.Error(),
		}
		sets[id] = message

		// HACK Redirect
		uri := fmt.Sprintf("/api/forms/result/%s", id.String())
		http.Redirect(w, r, uri, http.StatusFound)
	}
}

func ShowResult(w http.ResponseWriter, r *http.Request) {
	// HACK Get Struct
	params := mux.Vars(r)
	id, _ := uuid.FromString(params["data"])
	message := sets[id]
	if message.equals(dataset{}) {
		message = dataset{
			Status:  "Error",
			Message: "You are not allowed to access this page",
			Error:   "",
		}
	} else {
		delete(sets, id)
	}

	t, _ := template.ParseFiles("templates/done.html.tmpl")
	var tpl bytes.Buffer
	t.Execute(&tpl, message)
	body := template.HTML(tpl.String())

	data := content{
		Content: body,
	}
	t, _ = template.ParseFiles("templates/base.html.tmpl")
	t.Execute(w, data)
}

func (thingOne dataset) equals(thingTwo dataset) bool {
	return thingOne.Status == thingTwo.Status && thingOne.Message == thingTwo.Message && thingOne.Error == thingTwo.Error
}

func genRandomStatus() (status bool, str string, err error) {
	status = rand.Bool()
	str = "Password reset successfully"
	if !status {
		str = "Something went wrong"
	}
	err = errors.New(getRandStr(status))
	return
}

func getRandStr(ok bool) string {
	failed := []string{"Something went wrong", "Something happened", "All your base are belong to us"}
	if ok {
		return ""
	}
	return failed[rand.IntInRange(0, len(failed))]
}
