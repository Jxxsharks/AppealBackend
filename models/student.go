package models

type Student struct {
	ID          int    `gorm:"primarykey"`
	StudentID   string `gorm:"type:varchar(10);not null"`
	FirstName   string `gorm:"type:varchar(40)"`
	LastName    string `gorm:"type:varchar(40)"`
	CitizenID   string `gorm:"type:varchar(15)"`
	Field       string `gorm:"type:varchar(30)"`
	Faculty     string `gorm:"type:varchar(30)"`
	Email       string `gorm:"type:varchar(50);unique"`
	Phone       string `gorm:"type:varchar(10)"`
	Image       string
	IsChgPass   bool `gorm:"-"`
	Password    string
	PersonnelID int       `gorm:"type:int"`
	Appeal      []Appeal  `gorm:"foreignkey:StudentID;references:ID;"`
	Subject     []Subject `gorm:"many2many:enrolls"`
}
