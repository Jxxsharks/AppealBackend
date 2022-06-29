package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jxxsharks/petitionbackend/errs"
	"github.com/jxxsharks/petitionbackend/repository"
	"github.com/ledongthuc/goterators"
	"golang.org/x/crypto/bcrypt"
)

type studentResponse struct {
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Image string `json:"image,omitempty"`
}

type studentEnroll struct {
	StudentID int         `json:"student_id"`
	Personnel []personnel `json:"personnel,omitempty"`
	Subject   []subject   `json:"subject,omitempty"`
}

type personnel struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Image       string `json:"image"`
	PersonnelID int    `json:"personnel_id"`
}

type GetStudentInfo struct {
	SName     string `json:"student_name"`
	Identify  string `json:"identify"`
	Field     string `json:"field"`
	Faculty   string `json:"faculty"`
	AdvisorID int    `json:"advisor_id"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type subject struct {
	ID            int       `json:"id"`
	SubjectID     string    `json:"subject_id"`
	SubjectName   string    `json:"subject_name"`
	PersonnelId   int       `json:"personnel_id,omitempty"`
	PersonnelName string    `json:"personnel_name,omitempty"`
	Personnel     personnel `json:"personnel"`
}

type studentService struct {
	studentRepo repository.StudentRepository
}

type StudentService interface {
	GetPersonnel(int) (*studentEnroll, error)
	GetSubjects(int) (interface{}, error)
	GetStudentInfo(int) (*GetStudentInfo, error)
	GetStudentName(int) (*string, error)
	GetSubject(int) (*subject, error)
	UpdatePassword(int, string, string, string) error
	UpdateImage(int, multipart.File) (*studentResponse, error)
	UpdateInfo(int, string, string) (*studentResponse, error)
}

func NewStudentService(studentRepo repository.StudentRepository) StudentService {
	return studentService{studentRepo}
}

func (s studentService) GetPersonnel(id int) (*studentEnroll, error) {
	getSubjects, err := s.studentRepo.GetSubjects(id)
	if err != nil {
		return nil, err
	}

	subjectsID := []int{}
	for _, subject := range getSubjects.Subject {
		subjectsID = append(subjectsID, subject.ID)
	}
	getPersonnels, err := s.studentRepo.GetPersonnelOfSubjects(subjectsID)
	if err != nil {
		return nil, err
	}

	personnels := []personnel{}
	index := 0
	for _, subject := range getPersonnels {

		for _, person := range subject.Personnel {

			exist := goterators.Exist(personnels, personnel{PersonnelID: person.ID,
				FirstName: person.FirstName,
				LastName:  person.LastName})
			if !exist {
				personnels = append(personnels, personnel{
					PersonnelID: person.ID,
					FirstName:   person.FirstName,
					LastName:    person.LastName,
					Image:       person.Image,
				})
				index++
			}

		}
	}

	studentEnroll := studentEnroll{
		StudentID: id,
		Personnel: personnels,
	}

	return &studentEnroll, nil

}

func (s studentService) GetStudentInfo(id int) (*GetStudentInfo, error) {
	studentDB, err := s.studentRepo.GetStudentInfo(id)
	if err != nil {
		return nil, err
	}

	studentInfo := GetStudentInfo{}

	if studentDB.Gender == "Female" {
		studentInfo.SName = fmt.Sprintf("นางสาว %s %s", studentDB.FirstName, studentDB.LastName)
	} else {
		studentInfo.SName = fmt.Sprintf("นาย %s %s", studentDB.FirstName, studentDB.LastName)
	}

	studentInfo.Identify = studentDB.StudentID
	studentInfo.Field = studentDB.Field
	studentInfo.Faculty = studentDB.Faculty
	studentInfo.AdvisorID = studentDB.PersonnelID
	studentInfo.Email = studentDB.Email
	studentInfo.Phone = studentDB.Phone

	return &studentInfo, nil
}

func (s studentService) GetSubjects(id int) (interface{}, error) {
	getSubjects, err := s.studentRepo.GetSubjects(id)
	if err != nil {
		return nil, err
	}

	subjectsID := []int{}
	for _, subject := range getSubjects.Subject {
		subjectsID = append(subjectsID, subject.ID)
	}

	allSubjects, err := s.studentRepo.GetPersonnelOfSubjects(subjectsID)
	if err != nil {
		return err, nil
	}

	ownSubjects := []subject{}
	personnel := personnel{}
	for _, subjects := range allSubjects {

		for _, teach := range subjects.Personnel {

			if teach.ID == subjects.PersonnelID {
				fmt.Println(teach.ID)
				personnel.PersonnelID = teach.ID
				personnel.FirstName = teach.FirstName
				personnel.LastName = teach.LastName
			}
		}

		ownSubjects = append(ownSubjects, subject{
			ID:          subjects.ID,
			SubjectID:   subjects.SubjecID,
			SubjectName: subjects.SubjectName,
			Personnel:   personnel,
		})
	}

	studentEnroll := studentEnroll{
		StudentID: id,
		Subject:   ownSubjects,
	}

	return studentEnroll, nil
}

func (s studentService) GetSubject(id int) (*subject, error) {
	subjectID := []int{}
	subject := subject{}
	subjectID = append(subjectID, id)
	pOfSubject, err := s.studentRepo.GetPersonnelOfSubjects(subjectID)
	if err != nil {
		return nil, err
	}
	for _, p := range pOfSubject {
		subject.ID = p.ID
		subject.SubjectID = p.SubjecID
		subject.SubjectName = p.SubjectName
		subject.PersonnelId = p.PersonnelID
		for _, teach := range p.Personnel {
			subject.PersonnelName = teach.FirstName + " " + teach.LastName
		}
	}
	return &subject, nil
}

func (s studentService) UpdatePassword(id int, username string, password string, oldPassword string) error {
	passcrypt, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		return err
	}

	student, _ := s.studentRepo.LoginStudent(username)

	if student.Password == "" {
		match, _ := regexp.MatchString(oldPassword+"$", student.CitizenID)
		if match {
			err = s.studentRepo.UpdatePassword(id, string(passcrypt))
			if err != nil {
				return err
			}
		} else {
			return errs.NewNotImplement("can not update password")
		}
	} else if student.Password != "" {
		err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(oldPassword))
		if err != nil {
			return errs.NewNotImplement("can not update password")
		} else {
			err := s.studentRepo.UpdatePassword(id, string(passcrypt))
			if err != nil {
				return errs.NewNotImplement("can not update password")
			}
		}
	}
	return nil
}

func (s studentService) UpdateInfo(id int, email string, phone string) (*studentResponse, error) {
	student, err := s.studentRepo.UpdateStudent(id, email, phone)
	if err != nil {
		return nil, err
	}
	studentres := studentResponse{}
	studentres.Email = student.Email
	studentres.Phone = student.Phone

	return &studentres, nil
}

func (s studentService) UpdateImage(id int, img multipart.File) (*studentResponse, error) {
	var filename string
	var contentType string

	fileHeader := make([]byte, 512)
	img.Read(fileHeader)

	switch http.DetectContentType(fileHeader) {
	case "image/jpeg":
		filename = fmt.Sprintf("%s/%s%d%s%s", "img", "S", id, time.Now().String(), ".jpg")
		contentType = "image/jpeg"
	case "image/png":
		filename = fmt.Sprintf("%s/%s%d%s%s", "img", "S", id, time.Now().String(), ".png")
		contentType = "image/png"
	}
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("ap-southeast-1")},
		Profile: "default",
	})
	if err != nil {
		return nil, err
	}
	uploader := s3manager.NewUploader(sess)
	check, err := uploader.Upload(&s3manager.UploadInput{
		ACL:         aws.String("public-read"),
		Bucket:      aws.String("petitionplease"),
		Key:         aws.String(filename),
		Body:        img,
		ContentType: aws.String(contentType),
	})

	studentRes := studentResponse{}
	student, err := s.studentRepo.UpdateImage(id, check.Location)
	if err != nil {
		return nil, err
	}

	studentRes.Image = student.Image

	return &studentRes, nil

}

func (s studentService) GetStudentName(id int) (*string, error) {
	name, err := s.studentRepo.GetStudentName(id)
	if err != nil {
		return nil, errors.New("cannot get name")
	}
	return name, nil
}
