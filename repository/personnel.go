package repository

import (
	"gorm.io/gorm"
)

type personnel struct {
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
	Student     []student `gorm:"foreignkey:PersonnelID;references:ID;"`
	Appeal      []Appeal  `gorm:"foreignkey:PersonnelID;references:ID;"`
	PositionID  int       `gorm:"type:int"`
	OwnSubject  []subject `gorm:"foreignKey:PersonnelID"`
	Subject     []subject `gorm:"many2many:instructors"`
}

type personnelRepositoryDB struct {
	db *gorm.DB
}

type PersonnelRepossitory interface {
	GetPersonnel(string) (*personnel, error)
	GetPersonnelName(int) (*string, error)
	UpdatePassword(int, string) error
	UpdateImage(int, string) (*personnel, error)
	UpdatePersonnel(int, string, string) (*personnel, error)
}

func NewPersonnelRepository(db *gorm.DB) PersonnelRepossitory {

	return personnelRepositoryDB{db: db}
}

func (r personnelRepositoryDB) GetPersonnel(username string) (*personnel, error) {
	var err error

	personnel := personnel{}
	err = r.db.Where("personnel_id = ?", username).First(&personnel).Error
	if err != nil {
		return nil, err
	}

	return &personnel, nil
}

func (r personnelRepositoryDB) UpdatePassword(id int, password string) error {
	personnel := personnel{}
	personnel.ID = id
	err := r.db.Model(&personnel).Update("password", password).Error
	if err != nil {
		return err
	}
	return nil
}

func (r personnelRepositoryDB) UpdatePersonnel(id int, email string, phone string) (*personnel, error) {
	personnel := personnel{}
	personnel.ID = id
	personnel.Email = email
	personnel.Phone = phone

	err := r.db.Model(&personnel).Updates(&personnel).Error
	if err != nil {
		return nil, err
	}
	return &personnel, nil
}

func (r personnelRepositoryDB) UpdateImage(id int, image string) (*personnel, error) {
	personnel := personnel{}
	personnel.ID = id
	personnel.Image = image

	err := r.db.Model(&personnel).Update("image", personnel.Image).Error
	if err != nil {
		return nil, err
	}
	return &personnel, nil
}

func (r personnelRepositoryDB) GetPersonnelName(id int) (*string, error) {
	personnel := personnel{}
	personnel.ID = id
	err := r.db.Select("first_name", "last_name").Find(&personnel).Error
	if err != nil {
		return nil, err
	}
	name := personnel.FirstName + " " + personnel.LastName
	return &name, nil
}
