package models

import (
	u "github.com/l3njo/dropnote-api/utils"
	"errors"
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
	Name string `json:"name"`
	Mail string `json:"mail"`
	Pass string `json:"pass"`
	Auth string `sql:"-" json:"auth"`
}

// Validate checks incoming user details
func (user *User) Validate(db *gorm.DB) error {
	if !strings.Contains(user.Mail, "@") {
		return errors.New("Email address is required")
	}

	if len(user.Pass) < 6 {
		return errors.New("Password is required")
	}

	temp := &User{}
	err := db.Where(User{Mail: user.Mail}).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.New("Connection error. Please retry")
	}
	if temp.Mail != "" {
		return errors.New("Email address already in use by another user")
	}

	return nil
}

// Create adds the referenced user to the database
func (user *User) Create(db *gorm.DB) (map[string]interface{}, error) {
	if err := user.Validate(db); err != nil {
		return nil, err
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	user.Pass = string(hashedPassword)
	if db.Create(&user).Error != nil {
		return nil, errors.New("Failed to create user, connection error")
	}

	token := &Token{UserID: user.ID}
	auth := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), token)
	authString, _ := auth.SignedString([]byte(os.Getenv("token_password")))
	user.Auth = authString
	user.Pass = ""

	resp := u.Message(true, "User has been created")
	resp["data"] = user
	return resp, nil
}

// Login authorizes a user and assigns JWT token
func Login(db *gorm.DB, mail, pass string) (map[string]interface{}, error) {
	user := &User{}
	err := db.Where(User{Mail: mail}).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Email address not found")
		}
		return nil, errors.New("Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(pass))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, errors.New("Invalid login credentials. Please try again")
	}
	user.Pass = ""

	token := &Token{UserID: user.ID}
	auth := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), token)
	authString, _ := auth.SignedString([]byte(os.Getenv("token_password")))
	user.Auth = authString

	resp := u.Message(true, "Logged In")
	resp["data"] = user
	return resp, nil
}

// GetUser fetches the user from db
func GetUser(db *gorm.DB, u uuid.UUID) (user *User) {
	user = &User{}
	db.Where(User{Base: Base{ID: u}}).First(user)
	if user.Mail == "" {
		return nil
	}
	user.Pass = ""
	return
}

// GetUserByMail fetches the user from db
func GetUserByMail(db *gorm.DB, mail string) (user *User) {
	user = &User{}
	db.Where(User{Mail: mail}).First(user)
	if user.Mail == "" {
		return nil
	}
	user.Pass = ""
	return
}

// UpdateUser updates a user in the db
func UpdateUser(db *gorm.DB, user *User) (map[string]interface{}, error) {
	if err := db.Model(&user).Updates(User{Name: user.Name, Mail: user.Mail}).Error; err != nil {
		return nil, err
	}

	resp := u.Message(true, "success")
	resp["data"] = user
	return resp, nil
}

// UpdatePassword hashes and updates provided password
func UpdatePassword(db *gorm.DB, user *User) (err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	user.Pass = string(hash)
	if err = db.Model(&user).Updates(User{Pass: user.Pass}).Error; err != nil {
		return
	}
	return
}

// DeleteUser removes a user from the database
func DeleteUser(db *gorm.DB, user *User) (err error) {
	err = db.Delete(user).Error
	return
}
