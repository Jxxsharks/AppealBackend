package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jxxsharks/petitionbackend/services"
)

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	key := services.GetPublicKey()

	token, err := jwt.ParseWithClaims(cookie, &services.JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated",
			"status":  fiber.StatusUnauthorized,
		})
	}
	payload := token.Claims.(*services.JwtClaims)
	isStudent := strings.Contains(c.Path(), "/student")
	fmt.Println(payload)
	if payload.Role == "personnel" && isStudent || payload.Role == "student" && !isStudent {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	return c.Next()
}
