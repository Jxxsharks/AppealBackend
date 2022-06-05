package models

type Position struct {
	ID        int         `gorm:"primarykey"`
	Position  string      `gorm:"type:varchar(30)"`
	Personnel []Personnel `gorm:"foreignkey:PositionID;references:ID;"`
}
