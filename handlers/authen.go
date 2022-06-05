package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jxxsharks/petitionbackend/services"
)

type authenHandler struct {
	authSrv services.AuthenService
	jwtSrv  services.JWTServices
	psnSrv  services.PersonnelService
}

type AuthenHandler interface {
	Login(*fiber.Ctx) error
	Logout(*fiber.Ctx) error
}

func NewAuthenHandler(authenSrv services.AuthenService, jwtSrv services.JWTServices, psnSrv services.PersonnelService) AuthenHandler {
	return authenHandler{
		authSrv: authenSrv,
		jwtSrv:  jwtSrv,
		psnSrv:  psnSrv,
	}
}

func (h authenHandler) Login(c *fiber.Ctx) error {
	data := services.AuthenRequest{}
	c.BodyParser(&data)
	fmt.Println(data)
	user, err := h.authSrv.Login(data)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": err.Error(),
		})
	}
	role := "personnel"
	if user.IsStudent {
		role = "student"
	}
	token, err := h.jwtSrv.GenerateJWT(user.ID, role)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Crashed server",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	if user.PersonnelID != 0 {
		pname, _ := h.psnSrv.GetName(user.PersonnelID)
		user.PersonnelName = *pname
	}

	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"data":   user,
		"role":   user.IsStudent,
	})

}

func (h authenHandler) Logout(c *fiber.Ctx) error {
	h.jwtSrv.DeleteCookie(c)
	return c.JSON(fiber.Map{
		"message": "logout!",
	})
}
