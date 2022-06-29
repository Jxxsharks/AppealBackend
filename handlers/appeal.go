package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jxxsharks/petitionbackend/services"
	"github.com/ledongthuc/goterators"
)

type ScorePetitionResponse struct {
	ID           int       `json:"id"`
	Date         time.Time `json:"date"`
	Sname        string    `json:"subject"`
	Aname        string    `json:"student_name"`
	Pname        string    `json:"professor"`
	PID          int       `json:"-"`
	Stype        string    `json:"score_type"`
	Status       string    `json:"status"`
	PetitionType string    `json:"-"`
}

type appealHandler struct {
	appealSrv    services.AppealService
	personnelSrv services.PersonnelService
	studentSrv   services.StudentService
	jwtSrv       services.JWTServices
}

type AppealHandler interface {
	CreateAppeal(*fiber.Ctx) error
	GetScorePetition(*fiber.Ctx) error
	GetPetitionOfStudent(*fiber.Ctx) error
	GetPetitionOfPersonnel(*fiber.Ctx) error
	UpdatePersonnelPetition(*fiber.Ctx) error
	UpdateScorePetition(*fiber.Ctx) error
}

func NewAppealHandler(appealSrv services.AppealService, personnelSrv services.PersonnelService, studentSrv services.StudentService, jwtSrv services.JWTServices) AppealHandler {
	return appealHandler{appealSrv, personnelSrv, studentSrv, jwtSrv}
}

func (h appealHandler) CreateAppeal(c *fiber.Ctx) error {
	GetID, err := h.jwtSrv.GetUserID(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Broken Server",
		})
	}
	appealdata := services.AppealRequest{}
	err = c.BodyParser(&appealdata)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Request data failed",
		})
	}

	err = h.appealSrv.NewPetitionOfPersonnel(*GetID, appealdata)
	if err != nil {
		c.Status(fiber.StatusNotImplemented)
		return c.JSON(fiber.Map{
			"message": "Can't create petition",
		})
	}

	return c.JSON("success")
}

func (h appealHandler) UpdatePersonnelPetition(c *fiber.Ctx) error {
	GetID, err := h.jwtSrv.GetUserID(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Broken Server",
		})
	}
	appealdata := services.AppealRequest{}
	err = c.BodyParser(&appealdata)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Send some data",
		})
	}

	err = h.appealSrv.UpdatePersonnelPetition(*GetID, appealdata)
	if err != nil {
		c.Status(fiber.StatusNotImplemented)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON("success")
}

func (h appealHandler) UpdateScorePetition(c *fiber.Ctx) error {

	appealdata := services.AppealRequest{}
	err := c.BodyParser(&appealdata)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = h.appealSrv.UpdateScorePetition(appealdata.ID, appealdata)
	if err != nil {
		c.Status(fiber.StatusNotImplemented)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON("success")
}

func (h appealHandler) GetPetitionOfStudent(c *fiber.Ctx) error {
	GetID, err := h.jwtSrv.GetUserID(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Broken Server",
		})
	}
	types := c.Params("type")

	petitions, err := h.appealSrv.GetSpetitionOfStudent(*GetID, types)
	if err != nil {
		c.Status(fiber.StatusNotImplemented)
		return c.JSON(fiber.Map{
			"message": "Can not get data",
		})
	}

	scorePetitions := []ScorePetitionResponse{}

	for _, petition := range petitions {
		subject, _ := h.studentSrv.GetSubject(petition.SID)
		personnel, _ := h.personnelSrv.GetName(subject.PersonnelId)
		student, _ := h.studentSrv.GetStudentName(*GetID)
		scorePetitions = append(scorePetitions, ScorePetitionResponse{
			ID:     petition.ID,
			Date:   petition.Date,
			Aname:  *student,
			Sname:  subject.SubjectName,
			Pname:  *personnel,
			Stype:  petition.Stype,
			Status: petition.Status,
		})
	}
	return c.JSON(scorePetitions)
}

func (h appealHandler) GetPetitionOfPersonnel(c *fiber.Ctx) error {

	GetID, err := h.jwtSrv.GetUserID(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Broken Server",
		})
	}
	types := c.Params("type")
	position, _ := strconv.Atoi(c.Params("position"))
	fmt.Println(position)
	petitions, err := h.appealSrv.GetSpetitionOfPersonnel(types, position)
	if err != nil {
		c.Status(fiber.StatusNotImplemented)
		return c.JSON(fiber.Map{
			"message": "Can not get data",
		})
	}

	scorePetitions := []ScorePetitionResponse{}
	for _, petition := range petitions {
		subject, _ := h.studentSrv.GetSubject(petition.SID)
		personnel, _ := h.personnelSrv.GetName(subject.PersonnelId)
		student, _ := h.studentSrv.GetStudentName(petition.STDID)
		scorePetitions = append(scorePetitions, ScorePetitionResponse{
			ID:           petition.ID,
			Date:         petition.Date,
			Aname:        *student,
			Sname:        subject.SubjectName,
			Pname:        *personnel,
			PID:          subject.PersonnelId,
			Stype:        petition.Stype,
			Status:       petition.Status,
			PetitionType: "score",
		})
	}

	switch position {
	case 2:
		scorePetitions = goterators.Filter(scorePetitions, func(item ScorePetitionResponse) bool {

			if item.Status == "แจ้งอาจารย์ประจำวิชา" || item.Status == "แจ้งอาจารย์ประจำวิชาครั้งที่2" {
				fmt.Println(item)
				return item.PID == *GetID
			} else {
				return item.PetitionType == "score"
			}

		})
	case 3:
		scorePetitions = goterators.Filter(scorePetitions, func(item ScorePetitionResponse) bool {
			fmt.Println(item)
			return item.PID == *GetID
		})
	}

	return c.JSON(scorePetitions)

}

func (h appealHandler) GetScorePetition(c *fiber.Ctx) error {
	petition_id, _ := strconv.Atoi(c.Params("petitionid"))
	scoreAppeal, err := h.appealSrv.GetScorePetition(petition_id)
	if err != nil {
		c.Status(fiber.StatusNotImplemented)
		return c.JSON(fiber.Map{
			"message": "Can not get data",
		})
	}
	fmt.Println(scoreAppeal.Updated3)
	subject, _ := h.studentSrv.GetSubject(scoreAppeal.SubjectID)
	personnel, _ := h.personnelSrv.GetName(subject.PersonnelId)
	student, _ := h.studentSrv.GetStudentInfo(scoreAppeal.StudentID)
	advisor, _ := h.personnelSrv.GetName(student.AdvisorID)

	scoreAppeal.SName = subject.SubjectName
	scoreAppeal.SID = subject.SubjectID
	scoreAppeal.Pname = *personnel
	scoreAppeal.AName = student.SName
	scoreAppeal.IdenNumber = student.Identify
	scoreAppeal.Advisor = *advisor
	scoreAppeal.Faculty = student.Faculty
	scoreAppeal.Field = student.Field
	scoreAppeal.Email = student.Email
	scoreAppeal.Phone = student.Phone
	branchH, err := h.personnelSrv.GetName(1)
	dean, err := h.personnelSrv.GetName(2)
	scoreAppeal.BranchH = *branchH
	scoreAppeal.Dean = *dean

	return c.JSON(scoreAppeal)
}
