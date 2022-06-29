package services

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jxxsharks/petitionbackend/errs"
	"github.com/jxxsharks/petitionbackend/repository"
	"github.com/ledongthuc/goterators"
)

type AppealRequest struct {
	ID              int     `json:"id"`
	PetitionType    string  `json:"petition_type"`
	PetitionSubject string  `json:"petition_subject"`
	Detail          string  `json:"detail"`
	Gpax            float64 `json:"gpax"`
	Scoretype       string  `json:"score_type"`
	Semester        string  `json:"semester"`
	Year            string  `json:"year"`
	CreatedAt       time.Time
	Updated1        time.Time
	Updated2        time.Time
	Updated3        time.Time
	Updated4        time.Time
	Request1        string `json:"request_1"`
	Request2        string `json:"request_2"`
	Request3        string `json:"request_3"`
	Request4        string `json:"request_4"`
	Request5        string `json:"request_5"`
	File_1          string `json:"file_1"`
	Status          string `json:"status"`
	PersonnelID     int    `json:"personnel_id"`
	SubjectID       int    `json:"subject_id"`
	Base64          string `json:"base64"`
	Type            string `json:"type"`
	Issuccess       bool   `json:"is_success"`
}

type AppealResponse struct {
	ID              int       `json:"id"`
	PetitionType    string    `json:"petition_type"`
	PetitionSubject string    `json:"petition_subject"`
	Detail          string    `json:"detail"`
	Gpax            float64   `json:"gpax"`
	Scoretype       string    `json:"score_type"`
	Semester        string    `json:"semester"`
	Year            string    `json:"year"`
	CreatedAt       time.Time `json:"created_at"`
	Updated1        time.Time `json:"updated_1"`
	Updated2        time.Time `json:"updated_2"`
	Updated3        time.Time `json:"updated_3"`
	Updated4        time.Time `json:"updated_4"`
	Updated5        time.Time `json:"updated_5"`
	Request1        string    `json:"request_1"`
	Request2        string    `json:"request_2"`
	Request3        string    `json:"request_3"`
	Request4        string    `json:"request_4"`
	Request5        string    `json:"request_5"`
	File_1          string    `json:"file_1"`
	Status          string    `json:"status"`
	StudentID       int       `json:"student_id"`
	PersonnelID     int       `json:"personnel_id"`
	SubjectID       int       `json:"subject_id"`
	SID             string    `json:"subject_ID"`
	SName           string    `json:"subject_name"`
	Pname           string    `json:"professor_name"`
	AName           string    `json:"student_name"`
	IdenNumber      string    `json:"identify_number"`
	Field           string    `json:"field"`
	Faculty         string    `json:"faculty"`
	Advisor         string    `json:"advisor"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	BranchH         string    `json:"branch_head"`
	Dean            string    `json:"dean"`
}

type PersonnelPetitionResponse struct {
	ID              int     `json:"id"`
	Personnel       int     `json:"personnel"`
	PetitionSubject string  `json:"petition_subject"`
	Detail          string  `json:"detail"`
	Gpax            float64 `json:"gpax"`
	Review1         string  `json:"review_1"`
	Review2         string  `json:"review_2"`
	CreatedAt       time.Time
	Status          string `json:"status"`
}

type ScorePetitionResponse struct {
	ID     int       `json:"id"`
	Date   time.Time `json:"date"`
	SID    int       `json:"-"`
	STDID  int       `json:"-"`
	Stype  string    `json:"score_type"`
	Status string    `json:"status"`
}

var statusPersonnel = []string{"รออนุมัติ", "แจ้งหัวหน้าสาขา", "แจ้งอาจารย์ที่เกี่ยวข้อง", "แจ้งผลพิจารณาครั้งที่1", "ขอพิจารณาใหม่", "พิจารณาครั้งที่2", "แจ้งหัวหน้าสาขาครั้งที่2", "แจ้งอาจารย์ที่เกี่ยวข้องครั้งที่2", "แจ้งผลพิจารณาครั้งที่2"}
var statusScore = []string{
	"รออนุมัติ",
	"ไม่อนุมัติ",
	"แจ้งหัวหน้าสาขา",
	"แจ้งอาจารย์ประจำวิชา",
	"แจ้งผลพิจารณาครั้งที่1",
	"พิจารณาใหม่ครั้งที่1",
	"แจ้งผลพิจารณาครั้งที่2",
	"พิจารณาใหม่ครั้งที่2",
	"ระหว่างพิจารณา",
	"แจ้งหัวหน้าสาขาครั้งที่2",
	"แจ้งอาจารย์ประจำวิชาครั้งที่2",
	"แจ้งผลพิจารณาครั้งที่3",
}

type appealService struct {
	appealRepo repository.AppealRepository
}

type AppealService interface {
	NewPetitionOfPersonnel(int, AppealRequest) error
	GetScorePetition(int) (*AppealResponse, error)
	GetSpetitionOfStudent(int, string) ([]ScorePetitionResponse, error)
	GetSpetitionOfPersonnel(string, int) ([]ScorePetitionResponse, error)
	UpdatePersonnelPetition(int, AppealRequest) error
	UpdateScorePetition(int, AppealRequest) error
}

func NewAppealService(appealRepo repository.AppealRepository) AppealService {
	return appealService{appealRepo}
}

func (s appealService) NewPetitionOfPersonnel(id int, appealRequest AppealRequest) error {
	appeal := repository.Appeal{}
	if appealRequest.PetitionType == "personnel" {
		appeal = repository.Appeal{
			PetitionType:    appealRequest.PetitionType,
			PetitionSubject: appealRequest.PetitionSubject,
			Detail:          appealRequest.Detail,
			Gpax:            appealRequest.Gpax,
			StudentID:       id,
			PersonnelID:     &appealRequest.PersonnelID,
			SubjectID:       nil,
			Status:          statusPersonnel[0],
		}
	} else if appealRequest.PetitionType == "score" {
		appeal = repository.Appeal{
			PetitionType:    appealRequest.PetitionType,
			PetitionSubject: appealRequest.PetitionSubject,
			Scoretype:       appealRequest.Scoretype,
			Detail:          appealRequest.Detail,
			Gpax:            appealRequest.Gpax,
			Semester:        appealRequest.Semester,
			Year:            appealRequest.Year,
			StudentID:       id,
			PersonnelID:     nil,
			SubjectID:       &appealRequest.SubjectID,
			Status:          statusScore[0],
		}
	}

	err := s.appealRepo.CreateAppeal(appeal)
	if err != nil {
		return err
	}

	return nil
}

func (s appealService) UpdatePersonnelPetition(id int, appealRequest AppealRequest) error {
	appeal := repository.Appeal{}
	switch appealRequest.Status {
	case "รออนุมัติ":
		appeal.Request1 = appealRequest.Request1
		appeal.Status = statusPersonnel[1]

	case "แจ้งหัวหน้าสาขา":
		appeal.Status = statusPersonnel[2]
		appeal.Request2 = appealRequest.Request2
		appeal.Updated1 = time.Now().Local().UTC()

	case "แจ้งอาจารย์ที่เกี่ยวข้อง":
		appeal.Status = statusPersonnel[3]
	case "แจ้งผลพิจารณาครั้งที่1":
		if appealRequest.Issuccess {
			appeal.Status = "สำเร็จ"
		} else if !appealRequest.Issuccess {
			appeal.Status = statusPersonnel[4]
		}
	case "ขอพิจารณาใหม่":
		appeal.Status = statusPersonnel[5]
		appeal.Request3 = appealRequest.Request3
	case "พิจารณาครั้งที่2":
		appeal.Status = statusPersonnel[6]
		appeal.Request4 = appealRequest.Request4
		appeal.Updated2 = time.Now().Local().UTC()
	case "แจ้งหัวหน้าสาขาครั้งที่2":
		appeal.Status = statusPersonnel[7]

	case "แจ้งอาจารย์ที่เกี่ยวข้องครั้งที่2":
		appeal.Status = statusPersonnel[8]
	case "แจ้งผลพิจารณาครั้งที่2":
		appeal.Status = "สำเร็จ"
	}

	err := s.appealRepo.UpdatePersonnelPetition(appealRequest.ID, appealRequest.Status, appeal)
	if err != nil {
		return errs.NewNotImplement(err.Error())
	}
	return nil
}

func (s appealService) UpdateScorePetition(id int, appealRequest AppealRequest) error {

	appeal := repository.Appeal{}
	fmt.Println(id)
	switch appealRequest.Status {
	case "รออนุมัติ":
		if appealRequest.Issuccess {
			appeal.Status = statusScore[2]
		} else {
			appeal.Status = statusScore[1]
		}
		appeal.Request1 = appealRequest.Request1
	case "แจ้งหัวหน้าสาขา":
		appeal.Request2 = appealRequest.Request2
		appeal.Status = statusScore[3]
	case "แจ้งอาจารย์ประจำวิชา":

		if appealRequest.Base64 != "" {
			var filename string
			file, _ := base64.StdEncoding.DecodeString(appealRequest.Base64)
			read := bytes.NewReader(file)
			switch appealRequest.Type {
			case "application/pdf":
				filename = fmt.Sprintf("%s/%d%s", "petition", id, ".pdf")
			case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
				filename = fmt.Sprintf("%s/%d%s", "petition", id, ".xlsx")
			}

			sess, err := session.NewSessionWithOptions(session.Options{
				Config:  aws.Config{Region: aws.String("ap-southeast-1")},
				Profile: "default",
			})
			if err != nil {
				fmt.Println("here")
				return err
			}
			uploader := s3manager.NewUploader(sess)
			check, err := uploader.Upload(&s3manager.UploadInput{
				ACL:         aws.String("public-read"),
				Bucket:      aws.String("petitionplease"),
				Key:         aws.String(filename),
				Body:        read,
				ContentType: aws.String(appealRequest.Type),
			})
			appeal.File_1 = check.Location
		}

		appeal.Status = statusScore[4]
		appeal.Updated1 = time.Now().Local().UTC()
		appeal.Request3 = appealRequest.Request3
	case "แจ้งผลพิจารณาครั้งที่1":
		if appealRequest.Issuccess {
			appeal.Status = "สำเร็จ"
		} else {
			appeal.Status = statusScore[5]
		}
	case "พิจารณาใหม่ครั้งที่1":
		appeal.Status = statusScore[6]
		appeal.Updated2 = time.Now().Local().UTC()
	case "แจ้งผลพิจารณาครั้งที่2":
		if appealRequest.Issuccess {
			appeal.Status = "สำเร็จ"
			appeal.Updated3 = time.Now().Local().UTC()
		} else {
			appeal.Status = statusScore[7]
			appeal.Updated3 = time.Now().Local().UTC()
		}
	case "พิจารณาใหม่ครั้งที่2":
		appeal.Request4 = appealRequest.Request4
		appeal.Updated4 = time.Now().Local().UTC()
		appeal.Status = statusScore[8]
	case "ระหว่างพิจารณา":
		appeal.Request5 = appealRequest.Request5
		appeal.Status = statusScore[9]
	case "แจ้งหัวหน้าสาขาครั้งที่2":
		appeal.Status = statusScore[10]
		appeal.Updated5 = time.Now().Local().UTC()
	case "แจ้งอาจารย์ประจำวิชาครั้งที่2":
		appeal.Status = statusScore[11]
	case "แจ้งผลพิจารณาครั้งที่3":
		appeal.Status = "สำเร็จ"
	}

	err := s.appealRepo.UpdateScorePetition(appealRequest.ID, appealRequest.Status, appeal)
	if err != nil {
		return errs.NewNotImplement(err.Error())
	}

	return nil
}

func (s appealService) GetScorePetition(id int) (*AppealResponse, error) {

	getAppealResponse, err := s.appealRepo.GetScorePetition(id)
	if err != nil {
		return nil, errs.NewNotFoundError(err.Error())
	}

	appealResponse := AppealResponse{
		ID:              getAppealResponse.ID,
		PetitionType:    getAppealResponse.PetitionType,
		PetitionSubject: getAppealResponse.PetitionSubject,
		Detail:          getAppealResponse.Detail,
		Gpax:            getAppealResponse.Gpax,
		Scoretype:       getAppealResponse.Scoretype,
		Semester:        getAppealResponse.Semester,
		Year:            getAppealResponse.Year,
		CreatedAt:       getAppealResponse.CreatedAt,
		Updated1:        getAppealResponse.Updated1,
		Updated2:        getAppealResponse.Updated2,
		Updated3:        getAppealResponse.Updated3,
		Updated4:        getAppealResponse.Updated4,
		Updated5:        getAppealResponse.Updated5,
		Request1:        getAppealResponse.Request1,
		Request2:        getAppealResponse.Request2,
		Request3:        getAppealResponse.Request3,
		Request4:        getAppealResponse.Request4,
		Request5:        getAppealResponse.Request5,
		File_1:          getAppealResponse.File_1,
		Status:          getAppealResponse.Status,
		StudentID:       getAppealResponse.StudentID,
		SubjectID:       *getAppealResponse.SubjectID,
	}
	return &appealResponse, nil
}

func (s appealService) GetSpetitionOfStudent(id int, types string) ([]ScorePetitionResponse, error) {
	getPetition, err := s.appealRepo.GetPetitionForStudents(id, types)
	if err != nil {
		return nil, errors.New("Cannot Get Petition")
	}
	score := []ScorePetitionResponse{}

	for _, petition := range getPetition {
		score = append(score, ScorePetitionResponse{
			ID:     petition.ID,
			SID:    *petition.SubjectID,
			Date:   petition.CreatedAt,
			Stype:  petition.Scoretype,
			Status: petition.Status,
		})
	}
	return score, nil
}

func (s appealService) GetSpetitionOfPersonnel(types string, position int) ([]ScorePetitionResponse, error) {
	getPetition, err := s.appealRepo.GetPetitionForPersonnel(types)
	if err != nil {
		return nil, errors.New("Cannot Get Petition")
	}
	score := []ScorePetitionResponse{}

	for _, petition := range getPetition {
		score = append(score, ScorePetitionResponse{
			ID:     petition.ID,
			SID:    *petition.SubjectID,
			STDID:  petition.StudentID,
			Date:   petition.CreatedAt,
			Stype:  petition.Scoretype,
			Status: petition.Status,
		})
	}
	switch position {
	case 2:
		score = goterators.Filter(score, func(item ScorePetitionResponse) bool {
			return item.Status != "รออนุมัติ" && item.Status != "แจ้งผลพิจารณาครั้งที่1" && item.Status != "แจ้งผลพิจารณาครั้งที่2" && item.Status != "พิจารณาใหม่ครั้งที่2" &&
				item.Status != "ระหว่างพิจารณา" && item.Status != "แจ้งผลพิจารณาครั้งที่3" && item.Status != "ไม่อนุมัติ" && item.Status != "สำเร็จ"
		})
	case 3:
		score = goterators.Filter(score, func(item ScorePetitionResponse) bool {
			return item.Status == "แจ้งอาจารย์ประจำวิชา" || item.Status == "แจ้งอาจารย์ประจำวิชาครั้งที่2"
		})
	case 4:
		score = goterators.Filter(score, func(item ScorePetitionResponse) bool {
			return item.Status == "พิจารณาใหม่ครั้งที่2" || item.Status == "ระหว่างพิจารณา"
		})
	}
	return score, nil
}
