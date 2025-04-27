package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wmpn/lms-go-react/db"
	"github.com/wmpn/lms-go-react/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create a new book
func CreateBook(c *fiber.Ctx) error {
	book := models.Book{}

	if err := c.BodyParser(&book); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body!", "details": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.Collection.InsertOne(ctx, book)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert book"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Book created", "result": result.InsertedID})
}

// Update a book
func UpdateBook(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Params("id")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	book := models.Book{}
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body", "details": err.Error()})
	}

	update := bson.M{
		"$set": bson.M{
			"title":  book.Title,
			"author": book.Author,
		},
	}

	_, err = db.Collection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update book"})
	}

	return c.JSON(fiber.Map{"message": "Book updated"})
}

// Get a books
func GetBook(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	book := models.Book{}
	err = db.Collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&book)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
	}

	return c.JSON(fiber.Map{"result": book})
}

func GetBooks(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection.Find(ctx, bson.M{}) // Cursor is like a pointer that lets us loop over the results
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch books"})
	}
	defer cursor.Close(ctx) // Ensures the cursor is closed when done reading

	var books []models.Book
	for cursor.Next(ctx) { // cursor.Next() moves to next result in the cursor
		var book models.Book
		cursor.Decode(&book) // MongoDB returns BSON, we decode it into our Book struct
		// Decode() takes a pointer to the struct where the result will be stored
		books = append(books, book)
	}

	return c.JSON(books) // Response status 200 by default
}

func DeleteBook(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	_, err = db.Collection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete book"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Book deleted"})
}
