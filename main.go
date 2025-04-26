package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}

var books []Book

func main() {
	fmt.Println("Hello, World!")
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "Hello..."})
	})
	
	// Create a new book
	app.Post("/api/book", func(c *fiber.Ctx) error {
    book := Book{}

    if err := c.BodyParser(&book); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request!", "details": err.Error()})
    }

    book.ID = len(books) + 1
    books = append(books, book)

    return c.Status(201).JSON(fiber.Map{"message": "Book created", "book": book})
	})

	// Update a book
	app.Put("/api/updatebook/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, book := range books {
			if fmt.Sprintf("%d", book.ID) == id {
				updatedBook := Book{}
				
				if err := c.BodyParser(&updatedBook); err != nil {
					return c.Status(400).JSON(fiber.Map{"error": "Invalid request", "details": err.Error()})
				}

				updatedBook.ID = book.ID // keep original ID
				books[i] = updatedBook

				return c.Status(200).JSON(fiber.Map{"message": "Book updated", "book": updatedBook})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Book not found"})
	})

	app.Get("/api/books", func(c *fiber.Ctx) error {
		if len(books) == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "No books found"})
		}
		return c.Status(200).JSON(fiber.Map{"books": books})
	})

	app.Delete("/api/book/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		if len(books) == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "No books found"})
		}

		for i, book := range books {
			if fmt.Sprintf("%d", book.ID) == id {
				books = append(books[:i], books[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"message": "Book deleted"})
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": "Book not found"})
	})

	log.Fatal(app.Listen(":" + PORT))
}
