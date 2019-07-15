package models

import (
	"errors"
	"fmt"
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
func (c *Code) Execute(db *gorm.DB, user uuid.UUID) error {
	if !c.isValid(db, user) {
		return errors.New("Code is invalid")
	}
	fmt.Printf("Action %v performed on user %v\n", c.Action, user)
	db.Unscoped().Delete(c)
	return nil
}
