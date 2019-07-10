package models

import (
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
