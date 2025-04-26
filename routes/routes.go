package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wmpn/lms-go-react/controllers"
)

func RegisterBookRoutes(app *fiber.App) {
	api := app.Group("/api/books")

	api.Post("/", controllers.CreateBook)
	api.Get("/", controllers.GetBooks)
	api.Get("/:id", controllers.GetBook)
	api.Put("/:id", controllers.UpdateBook)
	api.Delete("/:id", controllers.DeleteBook)
}
