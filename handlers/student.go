package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jxxsharks/petitionbackend/services"
)

type supdatePassword struct {
	Password    string `json:"password"`
	ConfirmPass string `json:"confirm_password"`
}

type supdateInfo struct {
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required"`
}

type studentHandler struct {
	studentSrv services.StudentService
	jwtSrv     services.JWTServices
}
type StudentHandler interface {
	GetPersonnel(*fiber.Ctx) error
	GetSubjects(*fiber.Ctx) error
	GetSubject(*fiber.Ctx) error
	UpdatePassword(*fiber.Ctx) error
	UpdateImage(*fiber.Ctx) error
	UpdateStudent(*fiber.Ctx) error
}

func NewStudentHandler(studentSrv services.StudentService, jwtSrv services.JWTServices) StudentHandler {
	return studentHandler{studentSrv, jwtSrv}
}

func (h studentHandler) GetPersonnel(c *fiber.Ctx) error {
	GetId, err := h.jwtSrv.GetUserID(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Broken Server",
		})
	}

	getPersonnel, err := h.studentSrv.GetPersonnel(*GetId)
	if err != nil {
		return c.JSON("failed")
	}

	return c.JSON(getPersonnel)

}

func (h studentHandler) GetSubjects(c *fiber.Ctx) error {
	GetId, err := h.jwtSrv.GetUserID(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Broken Server",
		})
	}

	getPersonnel, err := h.studentSrv.GetSubjects(*GetId)
	if err != nil {
		return c.JSON("failed")
	}

	return c.JSON(getPersonnel)

}

func (h studentHandler) GetSubject(c *fiber.Ctx) error {
	_, err := h.jwtSrv.GetUserID(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Broken Server",
		})
	}

	subjectID, _ := strconv.Atoi(c.Params("subject_id"))

	getSubjectInfo, err := h.studentSrv.GetSubject(subjectID)
	if err != nil {
		return c.JSON("failed")
	}

	return c.JSON(getSubjectInfo)

}

func (h studentHandler) UpdatePassword(c *fiber.Ctx) error {
	update := supdatePassword{}
	if err := c.BodyParser(&update); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "please send data",
		})
	}

	if update.Password != update.ConfirmPass {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "data mismatch",
		})
	}

	Getid, _ := h.jwtSrv.GetUserID(c)
	id := *Getid
	err := h.studentSrv.UpdatePassword(id, update.Password)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "can not process",
		})
	}
	h.jwtSrv.DeleteCookie(c)
	return c.JSON(fiber.Map{
		"message": "success",
		"status":  fiber.StatusOK,
	})
}

func (h studentHandler) UpdateStudent(c *fiber.Ctx) error {
	update := supdateInfo{}
	if err := c.BodyParser(&update); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	errors := ValidateStruct(&update)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Wrong format",
			"Fields":  errors,
		})

	}

	Getid, err := h.jwtSrv.GetUserID(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Broken Server",
		})
	}

	updatestudent, err := h.studentSrv.UpdateInfo(*Getid, update.Email, update.Phone)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Broken Server",
		})
	}
	return c.JSON(updatestudent)
}

func (h studentHandler) UpdateImage(c *fiber.Ctx) error {
	Getid, err := h.jwtSrv.GetUserID(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Broken Server",
		})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"meassge": "please upload image",
		})
	}

	buffer, err := file.Open()
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Broken Server",
		})
	}

	updateImage, err := h.studentSrv.UpdateImage(*Getid, buffer)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Broken Server",
		})
	}
	return c.JSON(updateImage)
}
