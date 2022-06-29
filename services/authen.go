package services

import (
	"regexp"

	"github.com/jxxsharks/petitionbackend/errs"
	"github.com/jxxsharks/petitionbackend/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthenRequest struct {
	UserID   string `json:"username"`
	Password string `json:"password"`
}

type authenResponse struct {
	ID            int    `json:"id"`
	Username      string `json:"username"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Field         string `json:"field"`
	Faculty       string `json:"faculty"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Image         string `json:"image,omitempty"`
	IsStudent     bool   `json:"is_student"`
	IsChngPass    bool   `json:"is_change"`
	PositionID    int    `json:"position_id,omitempty"`
	PersonnelID   int    `json:"personnel_id,omitempty"`
	PersonnelName string `json:"personnel_name,omitempty"`
}

type authenService struct {
	personnelRepo repository.PersonnelRepossitory
	studentRepo   repository.StudentRepository
}

func NewAuthenService(personnelRepo repository.PersonnelRepossitory, studentRepo repository.StudentRepository) AuthenService {
	return authenService{personnelRepo: personnelRepo,
		studentRepo: studentRepo}
}

type AuthenService interface {
	Login(AuthenRequest) (*authenResponse, error)
}

func (s authenService) Login(request AuthenRequest) (*authenResponse, error) {

	var IsStudent bool
	role := request.UserID[0:1]
	switch role {
	case "P":

		dataDB, err := s.personnelRepo.GetPersonnel(request.UserID)
		IsStudent = false
		if err != nil {
			return nil, errs.NewNotFoundError("Not Found User")
		}

		if dataDB.Password == "" {

			match, _ := regexp.MatchString(request.Password+"$", dataDB.CitizenID)
			if match {
				authendata := authenResponse{
					ID:         dataDB.ID,
					Username:   dataDB.PersonnelID,
					FirstName:  dataDB.FirstName,
					LastName:   dataDB.LastName,
					Field:      dataDB.Field,
					Faculty:    dataDB.Faculty,
					Email:      dataDB.Email,
					Phone:      dataDB.Phone,
					Image:      dataDB.Image,
					IsChngPass: true,
					IsStudent:  IsStudent,
					PositionID: dataDB.PositionID,
				}
				return &authendata, nil
			}

		}

		if dataDB.Password != "" {

			err := bcrypt.CompareHashAndPassword([]byte(dataDB.Password), []byte(request.Password))
			if err == nil {
				authendata := authenResponse{
					ID:         dataDB.ID,
					Username:   dataDB.PersonnelID,
					FirstName:  dataDB.FirstName,
					LastName:   dataDB.LastName,
					Field:      dataDB.Field,
					Faculty:    dataDB.Faculty,
					Email:      dataDB.Email,
					Phone:      dataDB.Phone,
					Image:      dataDB.Image,
					IsChngPass: false,
					IsStudent:  IsStudent,
					PositionID: dataDB.PositionID,
				}
				return &authendata, nil
			}
		}

	case "S":

		dataDB, err := s.studentRepo.LoginStudent(request.UserID)
		IsStudent = true
		if err != nil {
			return nil, errs.NewNotFoundError("User Not Found")
		}

		if dataDB.Password == "" {
			match, _ := regexp.MatchString(request.Password+"$", dataDB.CitizenID)
			if match {
				authendata := authenResponse{
					ID:          dataDB.ID,
					Username:    dataDB.StudentID,
					FirstName:   dataDB.FirstName,
					LastName:    dataDB.LastName,
					Field:       dataDB.Field,
					Faculty:     dataDB.Faculty,
					Email:       dataDB.Email,
					Phone:       dataDB.Phone,
					Image:       dataDB.Image,
					IsStudent:   IsStudent,
					IsChngPass:  true,
					PersonnelID: dataDB.PersonnelID,
				}
				return &authendata, nil
			}
		}

		if dataDB.Password != "" {
			err := bcrypt.CompareHashAndPassword([]byte(dataDB.Password), []byte(request.Password))
			if err == nil {
				authendata := authenResponse{
					ID:          dataDB.ID,
					Username:    dataDB.StudentID,
					FirstName:   dataDB.FirstName,
					LastName:    dataDB.LastName,
					Field:       dataDB.Field,
					Faculty:     dataDB.Faculty,
					Email:       dataDB.Email,
					Phone:       dataDB.Phone,
					Image:       dataDB.Image,
					IsStudent:   IsStudent,
					IsChngPass:  false,
					PersonnelID: dataDB.PersonnelID,
				}
				return &authendata, nil
			}
		}

	}

	return nil, errs.NewNotFoundError("Wrong user or password")

}
