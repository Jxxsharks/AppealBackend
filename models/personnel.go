package models

type Personnel struct {
	ID          int    `gorm:"primarykey"`
	PersonnelID string `gorm:"type:varchar(10);not null"`
	FirstName   string `gorm:"type:varchar(40)"`
	LastName    string `gorm:"type:varchar(40)"`
	CitizenID   string `gorm:"type:varchar(15)"`
	Field       string `gorm:"type:varchar(30)"`
	Faculty     string `gorm:"type:varchar(30)"`
	Email       string `gorm:"type:varchar(50);unique"`
	Phone       string `gorm:"type:varchar(10)"`
	Image       string
	Password    string
	IsChgPass   bool      `gorm:"-"`
	Student     []Student `gorm:"foreignkey:PersonnelID;references:ID;"`
	Appeal      []Appeal  `gorm:"foreignkey:PersonnelID;references:ID;"`
	PositionID  int       `gorm:"type:int"`
	OwnSubject  []Subject `gorm:"foreignKey:PersonnelID"`
	Subject     []Subject `gorm:"many2many:instructors"`
}
