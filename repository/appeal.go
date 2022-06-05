package repository

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Appeal struct {
	ID              int     `gorm:"primarykey"`
	PetitionType    string  `gorm:"type:varchar(20);not null"`
	PetitionSubject string  `gorm:"type:varchar(30)"`
	Detail          string  `gorm:"not null"`
	Gpax            float64 `gorm:"decimal(4,2)"`
	Scoretype       string  `gorm:"type:varchar(20)"`
	Semester        string  `gorm:"type:varchar(10)"`
	Year            string  `gorm:"type:varchar(10)"`
	CreatedAt       time.Time
	Updated1        time.Time
	Updated2        time.Time
	Request1        string
	Request2        string
	Request3        string
	Request4        string
	Request5        string
	File_1          string
	Status          string `gorm:"type:varchar(40);not null"`
	StudentID       int    `gorm:"not null"`
	PersonnelID     *int
	SubjectID       *int
}

type appealRepository struct {
	db *gorm.DB
}

type AppealRepository interface {
	CreateAppeal(Appeal) error
	GetScorePetition(int) (*Appeal, error)
	GetPetitionForStudents(int, string) ([]Appeal, error)
	GetPetitionForPersonnel(string) ([]Appeal, error)
	UpdatePersonnelPetition(int, string, Appeal) error
}

func NewAppealRepository(db *gorm.DB) AppealRepository {
	db.AutoMigrate(&Appeal{})
	return appealRepository{db: db}
}

func (r appealRepository) GetScorePetition(id int) (*Appeal, error) {
	appeal := Appeal{}

	err := r.db.Where("id = ? and petition_type = ?", id, "score").Find(&appeal).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &appeal, nil
}

func (r appealRepository) CreateAppeal(appealdata Appeal) error {
	result := r.db.Create(&appealdata)
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}
	if result.RowsAffected <= 0 {
		return errors.New("Can't create appeal")
	}

	appealdata.ID = int(result.RowsAffected)

	return nil
}

func (r appealRepository) UpdatePersonnelPetition(id int, prevStatus string, appealdata Appeal) error {
	fmt.Println(appealdata)
	appeal := Appeal{}
	result := r.db.Model(&appeal).Where("id = ? and status = ?", id, prevStatus).Updates(&appealdata)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected <= 0 {
		return errors.New("Can't update appeal ")
	}

	return nil
}

func (r appealRepository) GetPetitionForStudents(studentID int, types string) ([]Appeal, error) {
	appeal := []Appeal{}

	err := r.db.Where("student_id = ? and petition_type = ?", studentID, types).Order("created_at desc").Find(&appeal).Error
	if err != nil {
		return nil, err
	}

	return appeal, nil
}

func (r appealRepository) GetPetitionForPersonnel(types string) ([]Appeal, error) {
	appeal := []Appeal{}

	err := r.db.Where("petition_type = ?", types).Order("created_at desc").Find(&appeal).Error
	if err != nil {
		return nil, err
	}

	return appeal, nil
}
