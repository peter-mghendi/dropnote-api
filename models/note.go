package models

import (
	"fmt"
	u "dropnote-backend/utils"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Note struct represents a note
type Note struct {
	Base
	Subject string    `json:"subject"`
	Content string    `json:"content"`
	Visible bool      `gorm:"default:true" json:"visible"`
	Creator uuid.UUID `gorm:"type:uuid" json:"creator"`
}

// Validate checks the required parameters sent through the http request body
// returns message and true if the requirement is met
func (note *Note) Validate() (map[string]interface{}, bool) {
	if note.Subject == "" {
		return u.Message(false, "Note subject should be on the payload"), false
	}
	if note.Content == "" {
		return u.Message(false, "Note content should be on the payload"), false
	}
	return u.Message(true, "success"), true
}

// Create adds a new note to the database
func (note *Note) Create(db *gorm.DB) map[string]interface{} {
	if resp, ok := note.Validate(); !ok {
		return resp
	}
	db.Create(note)
	resp := u.Message(true, "success")
	resp["note"] = note
	return resp
}

// GetNote returns a single note, if present, that matches provided criteria
func GetNote(db *gorm.DB, id uuid.UUID) (note *Note) {
	note = &Note{}
	err := db.Where(&Note{Base: Base{ID: id}}).First(note).Error
	if err != nil {
		return nil
	}
	return
}

// GetNotes returns an array of all notes
func GetNotes(db *gorm.DB) (notes []*Note) {
	notes = make([]*Note, 0)
	err := db.Find(&notes).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return
}
