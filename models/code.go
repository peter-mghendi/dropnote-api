package models

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Actions holds constants for different token types
var Actions = map[string]int{
	"reset": 0,
}

// Code represents a one time voucher provided to user in certain circumstances
type Code struct {
	Ticket uuid.UUID `json:"ticket"`
	UserID uuid.UUID `json:"userid"`
	Expiry time.Time `json:"expiry"`
	Action int       `json:"action"`
}

func (c *Code) exists(db *gorm.DB) bool {
	code := Code{}
	var count int
	db.Where(&Code{Ticket: c.Ticket}).First(&code).Count(&count)
	return count != 0
}

func getCode(db *gorm.DB, ticket uuid.UUID) *Code {
	code := Code{}
	db.Where(&Code{Ticket: ticket}).First(&code)
	return &code
}

// isValid checks the validity of a Code
func (c *Code) isValid(db *gorm.DB, user uuid.UUID) bool {
	if !c.exists(db) {
		return false
	}
	if time.Now().After(c.Expiry) {
		return false
	}
	if !uuid.Equal(c.UserID, user) {
		return false
	}
	return true
}

func (c *Code) perform(db *gorm.DB, body io.ReadCloser) error {
	switch c.Action {
	case Actions["reset"]:
		password := struct {
			Password string `json:"password"`
		}{}
		if err := json.NewDecoder(body).Decode(&password); err != nil {
			return errors.New("Error while decoding request body")
		}
		if password.Password == "" {
			return errors.New("Empty password string")
		}
		user := GetUser(db, c.UserID)
		user.Pass = password.Password
		if err := UpdatePassword(db, user); err != nil {
			return err
		}
	default:
		return errors.New("Invalid action")
	}
	return nil
}

// New returns a pointer to a new Code variable
func New(db *gorm.DB, action int, userID uuid.UUID) (uuid.UUID, error) {
	ticket, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}
	expiry := time.Now().Add(time.Hour * 24)
	code := &Code{ticket, userID, expiry, action}
	db.Create(code)
	return code.Ticket, nil
}

// Execute runs c.Action against user
func Execute(db *gorm.DB, body io.ReadCloser, code, user uuid.UUID) error {
	c := getCode(db, code)
	if !c.isValid(db, user) {
		return errors.New("Code is invalid")
	}
	if err := c.perform(db, body); err != nil {
		return err
	}
	db.Unscoped().Delete(c)
	return nil
}
