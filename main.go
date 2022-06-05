package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jxxsharks/petitionbackend/config"
	"github.com/jxxsharks/petitionbackend/routes"
)

func main() {
	db := config.SetupDatabaseConnection()
	_ = db

	app := fiber.New()

	app.Use(logger.New(logger.Config{
		TimeFormat: "02-Jan-2006",
		TimeZone:   "Asia/Bangkok",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	routes.Routeinit(db, app)

	app.Listen(":8000")

}
