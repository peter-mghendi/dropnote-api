package models

import (
	u "dropnote-backend/utils"
	"errors"
	"fmt"

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

// GetNotesFor returns an array of notes created by a specific user
func GetNotesFor(db *gorm.DB, user uuid.UUID) (notes []*Note) {
	notes = make([]*Note, 0)
	err := db.Where(&Note{Creator: user}).Find(&notes).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return
}

// UpdateNoteFor updates a note created by a specific user
func UpdateNoteFor(db *gorm.DB, note *Note) error {
	if err := db.Save(note).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// DeleteNoteFor deletes a note created by a specific user
// TODO move user auth to handler
func DeleteNoteFor(db *gorm.DB, note, user uuid.UUID) error {
	count := db.Where(&Note{Base: Base{ID: note}, Creator: user}).Delete(&Note{}).RowsAffected
	if count == 0 {
		return errors.New("No rows matching criteria")
	}
	return nil
}
