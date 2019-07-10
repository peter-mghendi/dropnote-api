package models

import (
	u "dropnote-backend/utils"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Token is a JWT claims struct
type Token struct {
	UserID uuid.UUID
	jwt.StandardClaims
}

// User is a struct to rep user
type User struct {
	Base
	Name string `json:"user"`
	Mail string `json:"mail"`
	Pass string `json:"pass"`
	Auth string `sql:"-" json:"auth"`
}

// Validate checks incoming user details
func (user *User) Validate(db *gorm.DB) (map[string]interface{}, bool) {
	if !strings.Contains(user.Mail, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(user.Pass) < 6 {
		return u.Message(false, "Password is required"), false
	}

	temp := &User{}
	err := db.Where(User{Mail: user.Mail}).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Mail != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

// Create adds the referenced user to the database
func (user *User) Create(db *gorm.DB) map[string]interface{} {
	if resp, ok := user.Validate(db); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	user.Pass = string(hashedPassword)

	if db.Create(&user).Error != nil {
		return u.Message(false, "Failed to create user, connection error.")
	}

	token := &Token{UserID: user.ID}
	auth := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), token)
	authString, _ := auth.SignedString([]byte(os.Getenv("token_password")))
	user.Auth = authString
	user.Pass = ""

	response := u.Message(true, "User has been created")
	response["user"] = user
	return response
}

// Login authorizes a user and assigns JWT token
func Login(db *gorm.DB, mail, pass string) map[string]interface{} {
	user := &User{}
	err := db.Where(User{Mail: mail}).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(pass))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	user.Pass = ""

	token := &Token{UserID: user.ID}
	auth := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), token)
	authString, _ := auth.SignedString([]byte(os.Getenv("token_password")))
	user.Auth = authString

	resp := u.Message(true, "Logged In")
	resp["user"] = user
	return resp
}

// GetUser fetches the user from db
func GetUser(db *gorm.DB, u uuid.UUID) (user *User) {
	user = &User{}
	db.Where(User{Base: Base{ID: u}}).First(user)
	if user.Mail == "" {
		return nil
	}
	user.Mail = ""
	return
}
