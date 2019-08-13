package models

import (
	"errors"
	"log"

	u "github.com/l3njo/dropnote-api/utils"

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
	resp["data"] = note
	return resp
}

// GetNote returns a single note, if present, that matches provided criteria
// and has visible == true
func GetNote(db *gorm.DB, id uuid.UUID) (note *Note) {
	note = &Note{}
	err := db.Where(&Note{Base: Base{ID: id}, Visible: true}).First(note).Error
	if err != nil {
		return nil
	}
	return
}

// GetNoteUnscoped returns a single note, if present, that matches provided criteria
func GetNoteUnscoped(db *gorm.DB, id uuid.UUID) (note *Note) {
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
		log.Println(err)
		return nil
	}
	return
}

// GetNotesFor returns an array of notes created by a specific user
func GetNotesFor(db *gorm.DB, user uuid.UUID) (notes []*Note) {
	notes = make([]*Note, 0)
	err := db.Where(&Note{Creator: user}).Find(&notes).Error
	if err != nil {
		log.Println(err)
		return nil
	}
	return
}

// UpdateNote updates a note
func UpdateNote(db *gorm.DB, note *Note) (map[string]interface{}, error) {
	if err := db.Model(&note).Updates(Note{Subject: note.Subject, Content: note.Content}).Error; err != nil {
		return nil, err
	}
	resp := u.Message(true, "success")
	resp["data"] = note
	return resp, nil
}

// ToggleNote updates a note
func ToggleNote(db *gorm.DB, note *Note) (map[string]interface{}, error) {
	if err := db.Save(note).Error; err != nil {
		return nil, err
	}
	resp := u.Message(true, "success")
	resp["data"] = note
	return resp, nil
}

// DeleteNote deletes a note
func DeleteNote(db *gorm.DB, note *Note) error {
	count := db.Delete(note).RowsAffected
	if count == 0 {
		return errors.New("No rows matching criteria")
	}
	return nil
}
