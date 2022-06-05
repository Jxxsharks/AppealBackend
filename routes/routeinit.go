package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jxxsharks/petitionbackend/handlers"
	"github.com/jxxsharks/petitionbackend/middleware"
	"github.com/jxxsharks/petitionbackend/repository"
	"github.com/jxxsharks/petitionbackend/services"
	"gorm.io/gorm"
)

func Routeinit(db *gorm.DB, c *fiber.App) {
	personnelRepo := repository.NewPersonnelRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	appealRepo := repository.NewAppealRepository(db)
	authenSrv := services.NewAuthenService(personnelRepo, studentRepo)
	personnelSrv := services.NewPersonnelService(personnelRepo)
	studentSrv := services.NewStudentService(studentRepo)
	appealSrv := services.NewAppealService(appealRepo)
	jwtSrv := services.NewJWTService()
	authenHandler := handlers.NewAuthenHandler(authenSrv, jwtSrv, personnelSrv)
	personnelHandler := handlers.NewPersonnelHandler(personnelSrv, jwtSrv)
	studentHandler := handlers.NewStudentHandler(studentSrv, jwtSrv)
	appealHandler := handlers.NewAppealHandler(appealSrv, personnelSrv, studentSrv, jwtSrv)

	c.Post("/login", authenHandler.Login)

	personnel := c.Group("/personnel")
	personnelAuthen := personnel.Use(middleware.IsAuthenticated)
	personnelAuthen.Get("petition/type::type/position::position", appealHandler.GetPetitionOfPersonnel)
	personnelAuthen.Get("petition/petitionid::petitionid", appealHandler.GetScorePetition)
	personnelAuthen.Put("/password", personnelHandler.UpdatePassword)
	personnelAuthen.Put("/info", personnelHandler.UpdatePersonnel)
	personnelAuthen.Put("/image", personnelHandler.UpdateImage)
	personnelAuthen.Put("/petition/personnel", appealHandler.UpdatePersonnelPetition)

	student := c.Group("/student")
	studentAuthen := student.Use(middleware.IsAuthenticated)
	studentAuthen.Get("/personnel", studentHandler.GetPersonnel)
	studentAuthen.Get("/subjects", studentHandler.GetSubjects)
	studentAuthen.Get("/subjects/:subject_id", studentHandler.GetSubject)

	studentAuthen.Post("/petition", appealHandler.CreateAppeal)
	studentAuthen.Get("/petition/petitionid::petitionid", appealHandler.GetScorePetition)
	studentAuthen.Put("/petition/personnel", appealHandler.UpdatePersonnelPetition)
	studentAuthen.Get("/petition/type::type", appealHandler.GetPetitionOfStudent)
	studentAuthen.Get("/petition/type::type/:id", personnelHandler.GetName)
	studentAuthen.Put("/password", studentHandler.UpdatePassword)
	studentAuthen.Put("/info", studentHandler.UpdateStudent)
	studentAuthen.Put("/image", studentHandler.UpdateImage)
	studentAuthen.Post("/logout", authenHandler.Logout)

}
