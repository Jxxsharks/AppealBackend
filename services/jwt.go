package services

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type JwtClaims struct {
	jwt.StandardClaims
	Role string
}

type jwtService struct {
	privateKey *rsa.PrivateKey
	publickey  *rsa.PublicKey
	issuer     string
}

type JWTServices interface {
	GenerateJWT(int, string) (string, error)
	DeleteCookie(*fiber.Ctx)
	GetUserID(*fiber.Ctx) (*int, error)
}

func NewJWTService() JWTServices {
	return jwtService{
		privateKey: getPrivateKey(),
		publickey:  GetPublicKey(),
		issuer:     "PetitionAdmin",
	}
}

func getPrivateKey() *rsa.PrivateKey {
	private, err := ioutil.ReadFile("./certs/private.pem")
	if err != nil {
		fmt.Println(err)
	}

	rsaprivate, err := jwt.ParseRSAPrivateKeyFromPEM(private)
	if err != nil {
		fmt.Println(err)
	}

	return rsaprivate
}

func GetPublicKey() *rsa.PublicKey {
	public, err := ioutil.ReadFile("./certs/public.pem")
	if err != nil {
		fmt.Println(err)
	}

	rsapublic, err := jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		fmt.Println(err)
	}
	return rsapublic
}

func (s jwtService) GenerateJWT(id int, role string) (string, error) {

	claim := JwtClaims{}
	claim.Issuer = s.issuer
	claim.Subject = strconv.Itoa(id)
	claim.ExpiresAt = time.Now().Add(time.Hour).Unix()
	claim.Role = role
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	return token.SignedString(s.privateKey)
}

func (s jwtService) GetUserID(c *fiber.Ctx) (*int, error) {
	cookie := c.Cookies("jwt")

	key := GetPublicKey()

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	payload := token.Claims.(*jwt.StandardClaims)

	id, _ := strconv.Atoi(payload.Subject)

	return &id, nil
}

func (s jwtService) DeleteCookie(c *fiber.Ctx) {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
}
