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
	"golang.org/x/crypto/bcrypt"
)

type personnelResponse struct {
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Image string `json:"image,omitempty"`
}

type personnelService struct {
	personnelRepo repository.PersonnelRepossitory
}

type PersonnelService interface {
	GetName(int) (*string, error)
	UpdatePassword(int, string, string, string) error
	UpdateImage(int, multipart.File) (*personnelResponse, error)
	UpdateInfo(int, string, string) (*personnelResponse, error)
}

func NewPersonnelService(personnelRepo repository.PersonnelRepossitory) PersonnelService {
	return personnelService{personnelRepo: personnelRepo}
}

func (s personnelService) UpdatePassword(id int, username string, password string, oldPassword string) error {
	var err error
	passcrypt, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		return err
	}

	personnel, _ := s.personnelRepo.GetPersonnel(username)

	if personnel.Password == "" {
		match, _ := regexp.MatchString(oldPassword+"$", personnel.CitizenID)
		if match {
			err := s.personnelRepo.UpdatePassword(id, string(passcrypt))
			if err != nil {
				return err
			}
		} else {
			return errs.NewNotImplement("can not update password")
		}
	} else if personnel.Password != "" {
		err := bcrypt.CompareHashAndPassword([]byte(personnel.Password), []byte(oldPassword))
		if err != nil {
			return errs.NewNotImplement("can not update password")
		} else {
			err := s.personnelRepo.UpdatePassword(id, string(passcrypt))
			if err != nil {
				return errs.NewNotImplement("can not update password")
			}
		}
	}

	return nil
}

func (s personnelService) UpdateInfo(id int, email string, phone string) (*personnelResponse, error) {
	personnel, err := s.personnelRepo.UpdatePersonnel(id, email, phone)
	if err != nil {
		return nil, err
	}
	personnelres := personnelResponse{}
	personnelres.Email = personnel.Email
	personnelres.Phone = personnel.Phone

	return &personnelres, nil
}

func (s personnelService) UpdateImage(id int, img multipart.File) (*personnelResponse, error) {
	var filename string
	var contentType string

	fileHeader := make([]byte, 512)
	img.Read(fileHeader)

	switch http.DetectContentType(fileHeader) {
	case "image/jpeg":
		filename = fmt.Sprintf("%s/%s%d%s%s", "img", "P", id, time.Now().String(), ".jpg")
		contentType = "image/jpeg"
	case "image/png":
		filename = fmt.Sprintf("%s/%s%d%s%s", "img", "P", id, time.Now().String(), ".png")
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

	personnelRes := personnelResponse{}
	personnel, err := s.personnelRepo.UpdateImage(id, check.Location)
	if err != nil {
		return nil, err
	}
	personnelRes.Image = personnel.Image

	return &personnelRes, nil
}

func (s personnelService) GetName(id int) (*string, error) {
	name, err := s.personnelRepo.GetPersonnelName(id)
	if err != nil {
		return nil, errors.New("Cannot get name")
	}

	return name, nil
}
