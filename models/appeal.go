package models

import (
	"time"

	"gorm.io/gorm"
)

type Appeal struct {
	ID              int    `gorm:"primarykey"`
	PetitionSubject string `gorm:"type:varchar(10);not null"`
	Detail          string `gorm:"not null"`
	Scoretype       string `gorm:"type:varchar(20)"`
	Semester        string `gorm:"type:varchar(10)"`
	Year            string `gorm:"type:varchar(10)"`
	CreatedAt       time.Time
	Updated1        time.Time
	Updated2        time.Time
	Request1        string
	Request2        string
	Request3        string
	Request4        string
	Request5        string
	File_1          string
	File_2          string
	AppealFile      string
	Status          string `gorm:"type:varchar(20);not null"`
	StudentID       int    `gorm:"not null"`
	PersonnelID     int
	SubjectID       int
}

type test struct {
	gorm.Model
}
