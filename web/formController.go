package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

const api = "https://dropnote-api.herokuapp.com/api/"

type content struct {
	Content template.HTML
}

type dataset struct {
	Status, Message string
}

type response struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func (thingOne dataset) equals(thingTwo dataset) bool {
	return thingOne.Status == thingTwo.Status && thingOne.Message == thingTwo.Message
}

// HACK transfer dataset between endpoints
var sets = make(map[uuid.UUID]dataset)

func init() {
	sets[uuid.Nil] = dataset{
		Status:  "Error",
		Message: "You are not allowed to access this page",
	}
}

func getResponse(password string, user, code uuid.UUID) (response, error) {
	reply := response{}
	requestBody, err := json.Marshal(map[string]string{
		"password": password,
	})
	if err != nil {
		return response{}, err
	}

	uri := fmt.Sprintf("%suser/%s/action/%s", api, user.String(), code.String())
	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return response{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &reply)
	if err != nil {
		return response{}, err
	}
	return reply, nil
}

// DoReset recieves, validates and processes reset request
func DoReset(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := uuid.FromStringOrNil(params["user"])
	code := uuid.FromStringOrNil(params["code"])

	if user == uuid.Nil || code == uuid.Nil {
		uri := fmt.Sprintf("/api/forms/result/%s", uuid.Nil.String())
		http.Redirect(w, r, uri, http.StatusFound)
		return
	}

	if r.Method == "GET" {
		bytes, _ := ioutil.ReadFile("templates/auth.html.tmpl")
		body := template.HTML(bytes)
		data := content{
			Content: body,
		}
		t, _ := template.ParseFiles("templates/base.html.tmpl")
		t.Execute(w, data)
		return
	}
	r.ParseForm()
	title, message := "Success", "Password reset successfully"

	if len(r.Form["password"]) != 1 {
		log.Println("Invalid input")
		return
	}
	pass := r.Form["password"][0]

	resp, err := getResponse(pass, user, code)
	handle(err)
	if !resp.Status {
		title, message = "Failed", "Something went wrong"
	}

	// HACK Save struct and redirect
	id := uuid.NewV4()
	data := dataset{Status: title, Message: message}
	sets[id] = data
	uri := fmt.Sprintf("/api/forms/result/%s", id.String())
	http.Redirect(w, r, uri, http.StatusFound)
}

// ShowResult shows results of processing
func ShowResult(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := uuid.FromStringOrNil(params["data"])
	message := sets[id]
	if id != uuid.Nil {
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

func handle(e error) {
	if e != nil {
		log.Println(e)
	}
}
