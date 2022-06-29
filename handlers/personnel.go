package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jxxsharks/petitionbackend/services"
)

type pupdateData struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	ConfirmPass string `json:"confirm_password"`
	OldPassword string `json:"old_password"`
}

type pupdateInfo struct {
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required"`
}

type personnelHandler struct {
	personnelSrv services.PersonnelService
	jwtSrv       services.JWTServices
}

type PersonnelHandler interface {
	GetName(*fiber.Ctx) error
	UpdatePassword(*fiber.Ctx) error
	UpdateImage(*fiber.Ctx) error
	UpdatePersonnel(*fiber.Ctx) error
}

func NewPersonnelHandler(personnelSrv services.PersonnelService, jwtSrv services.JWTServices) PersonnelHandler {
	return personnelHandler{personnelSrv, jwtSrv}
}

func (h personnelHandler) UpdatePassword(c *fiber.Ctx) error {
	update := pupdateData{}
	if err := c.BodyParser(&update); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Bad request",
		})
	}

	if update.Password != update.ConfirmPass {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "data mismatch",
		})
	}
	Getid, err := h.jwtSrv.GetUserID(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Broken Server",
		})
	}

	err = h.personnelSrv.UpdatePassword(*Getid, update.Username, update.Password, update.OldPassword)
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

func (h personnelHandler) UpdatePersonnel(c *fiber.Ctx) error {

	Getid, err := h.jwtSrv.GetUserID(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Broken Server",
		})
	}
	update := pupdateInfo{}
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
	updateperson, err := h.personnelSrv.UpdateInfo(*Getid, update.Email, update.Phone)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Broken Server",
		})
	}
	return c.JSON(updateperson)
}

func (h personnelHandler) UpdateImage(c *fiber.Ctx) error {
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
	updateImage, err := h.personnelSrv.UpdateImage(*Getid, buffer)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Broken Server",
		})
	}
	return c.JSON(updateImage)
}

func (h personnelHandler) GetName(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	name, err := h.personnelSrv.GetName(id)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Broken Server",
		})
	}
	return c.JSON(name)
}
