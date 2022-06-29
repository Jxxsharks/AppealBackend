package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type student struct {
	ID          int    `gorm:"primarykey"`
	StudentID   string `gorm:"type:varchar(10);not null"`
	FirstName   string `gorm:"type:varchar(40)"`
	LastName    string `gorm:"type:varchar(40)"`
	CitizenID   string `gorm:"type:varchar(15)"`
	Field       string `gorm:"type:varchar(30)"`
	Faculty     string `gorm:"type:varchar(30)"`
	Email       string `gorm:"type:varchar(50);unique"`
	Phone       string `gorm:"type:varchar(10)"`
	Gender      string `gorm:"type:varchar(10)"`
	Image       string
	IsChgPass   bool `gorm:"-"`
	Password    string
	PersonnelID int       `gorm:"type:int"`
	Appeal      []Appeal  `gorm:"foreignkey:StudentID;references:ID;"`
	Subject     []subject `gorm:"many2many:enrolls"`
}

type subject struct {
	ID          int    `gorm:"primarykey"`
	SubjecID    string `gorm:"type:varchar(10)"`
	SubjectName string `gorm:"type:varchar(50)"`
	PersonnelID int
	Personnel   []personnel `gorm:"many2many:instructors"`
	Appeal      []Appeal    `gorm:"foreignkey:SubjectID;references:ID;"`
}

type studentRepositoryDB struct {
	db *gorm.DB
}

type StudentRepository interface {
	LoginStudent(string) (*student, error)
	GetStudentName(int) (*string, error)
	GetStudentInfo(int) (*student, error)
	GetSubjects(int) (*student, error)
	GetPersonnelOfSubjects([]int) ([]subject, error)
	UpdatePassword(int, string) error
	UpdateImage(int, string) (*student, error)
	UpdateStudent(int, string, string) (*student, error)
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	db.AutoMigrate(&student{})
	return studentRepositoryDB{db: db}
}

func (r studentRepositoryDB) GetSubjects(id int) (*student, error) {
	getSubjects := student{}

	err := r.db.Preload("Subject").Find(&getSubjects, id).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &getSubjects, nil
}

func (r studentRepositoryDB) GetStudentInfo(id int) (*student, error) {
	student := student{}
	student.ID = id
	err := r.db.Find(&student).Error
	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (r studentRepositoryDB) GetPersonnelOfSubjects(subjectsID []int) ([]subject, error) {
	getPersonnels := []subject{}
	err := r.db.Preload("Personnel").Find(&getPersonnels, "id IN ?", subjectsID).Error
	if err != nil {
		return nil, err
	}
	return getPersonnels, nil
}

func (r studentRepositoryDB) LoginStudent(username string) (*student, error) {
	var err error
	student := student{}

	err = r.db.Where("student_id = ?", username).First(&student).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &student, nil
}

func (r studentRepositoryDB) UpdatePassword(id int, password string) error {
	student := student{}
	student.ID = id
	err := r.db.Model(&student).Update("password", password).Error
	if err != nil {
		return err
	}

	return nil
}

func (r studentRepositoryDB) UpdateStudent(id int, email string, phone string) (*student, error) {
	student := student{}
	student.ID = id
	student.Email = email
	student.Phone = phone

	err := r.db.Model(&student).Updates(&student).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &student, nil
}

func (r studentRepositoryDB) UpdateImage(id int, image string) (*student, error) {
	student := student{}
	student.ID = id
	student.Image = image

	err := r.db.Model(&student).Update("image", student.Image).Error
	if err != nil {
		return nil, err
	}

	return &student, err
}

func (r studentRepositoryDB) GetStudentName(id int) (*string, error) {
	student := student{}
	student.ID = id
	err := r.db.Select("first_name", "last_name").Find(&student).Error
	if err != nil {
		return nil, err
	}
	name := student.FirstName + " " + student.LastName
	return &name, nil

}
