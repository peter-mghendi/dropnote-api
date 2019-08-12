package controllers

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/l3njo/dropnote-api/models"

	"github.com/go-mail/mail"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/matcornic/hermes/v2"
	uuid "github.com/satori/go.uuid"
)

// Application is a copy of the app.App struct
type Application struct {
	Router *mux.Router
	DB     *gorm.DB
}

type key string

// UserKey is a constant context key
const (
	UserKey key = "user"
	api         = "https://dropnote-api.herokuapp.com/api/"
)

// App is an instance of Application
var (
	App *Application
	h   hermes.Hermes
)

func init() {
	App = &Application{}
	h = hermes.Hermes{
		Product: hermes.Product{
			Name:        "DropNote",
			Link:        "https://drop-note.herokuapp.com/",
			Logo:        "https://drop-note.herokuapp.com/static/img/favicon.png",
			Copyright:   "Copyright © 2019 DropNote. All rights reserved.",
			TroubleText: "If the {ACTION} button is not working for you, just copy and paste the URL below into your web browser.",
		},
	}
}

func buildEmail(user *models.User, link string) hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Name: user.Name,
			Intros: []string{
				"You have received this email because a password reset request for your DropNote account was received.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click the button below to reset your password:",
					Button: hermes.Button{
						Color: "#DC4D2F",
						Text:  "Reset Password",
						Link:  link,
					},
				},
			},
			Outros: []string{
				"If you did not request a password reset, no further action is required on your part.",
			},
			Signature: "Thanks",
		},
	}
}

func sendMail(address, content string) error {
	m := mail.NewMessage()
	m.SetHeader("From", os.Getenv("MAIL_USER"))
	m.SetHeader("To", address)
	m.SetHeader("Subject", "RE: Password Reset")
	m.SetHeader("X-Entity-Ref-ID", uuid.NewV4().String())
	m.SetBody("text/html", content)

	host, user, pass := os.Getenv("MAIL_HOST"), os.Getenv("MAIL_USER"), os.Getenv("MAIL_PASS")
	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		return err
	}

	d := mail.NewDialer(host, port, user, pass)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func emailTokenToUser(user *models.User, code uuid.UUID) error {
	link := fmt.Sprintf("%sforms/users/%s/reset/%s\n", api, user.ID.String(), code.String())
	email := buildEmail(user, link)
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		return err
	}

	if err = sendMail(user.Mail, emailBody); err != nil {
		return err
	}

	return nil
}

// Handle deals with top-level errors
func Handle(e error) {
	if e != nil {
		log.Println(e)
	}
}
