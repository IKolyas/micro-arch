package infrastructure

import (
	"github.com/IKolyas/otus-highload/internal/infrastructure/controller"
	"github.com/IKolyas/otus-highload/internal/infrastructure/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func Router() *fiber.App {

	router := fiber.New()
	router.Use(logger.New())
	//---------------------------------------------------------------------------------------//
	router.Post("/login", controller.Login)
	router.Post("/register", controller.Register)
	router.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))
	//---------------------------------------------------------------------------------------//
	router.Get("/faker/create/:count", middleware.JWTProtected, controller.FakerUser)
	//---------------------------------------------------------------------------------------//
	v1 := router.Group("/api/v1")
	v1.Use(middleware.JWTProtected)
	v1.Get("/users/find/firstName::firstName/secondName::secondName", controller.SearchUser)
	v1.Get("/users/:id", controller.GetUser)

	return router
}
