package models

type Subject struct {
	ID          int    `gorm:"primarykey"`
	SubjecID    string `gorm:"type:varchar(10)"`
	SubjectName string `gorm:"type:varchar(50)"`
	PersonnelID int
	Personnel   []Personnel `gorm:"many2many:instructors"`
	Appeal      []Appeal    `gorm:"foreignkey:SubjectID;references:ID;"`
}
