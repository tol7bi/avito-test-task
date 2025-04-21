package app_setup

import (
	"pvz-backend/internal/http"
	"pvz-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(r fiber.Router) {

	r.Post("/register", http.Register)
	r.Post("/login", http.Login)
	r.Post("/dummyLogin", http.DummyLogin)


	// защищённые 
	auth := r.Group("/", middleware.JWTMiddleware())

	// общие 
	auth.Get("/pvz", http.GetPVZListHandler)

	// только moderator
	auth.Post("/pvz", http.CreatePVZHandler, middleware.RequireRole("moderator"))


	// только employee
	employee := auth.Group("/", middleware.RequireRole("employee"))
	employee.Post("/pvz/:pvzId/delete_last_product", http.DeleteLastProductHandler)
	employee.Post("/pvz/:pvzId/close_last_reception", http.CloseLastReceptionHandler)
	employee.Post("/receptions", http.CreateReceptionHandler)
	employee.Post("/products", http.AddProductHandler)
}
