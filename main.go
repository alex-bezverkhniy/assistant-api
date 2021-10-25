package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Data Structures

// Option for user selection
type Option struct {
	ID            int    `json:"id"`
	Body          string `json:"body"`
	NextMessageID int    `json:"nextMessageId"`
}

// List of Options
type Options []Option

// Message/Question
type Message struct {
	ID      int     `json:"id"`
	Body    string  `json:"body"`
	Options Options `json:"options"`
}

// Store
type Store struct {
	messages []Message
	fileName string
}

func main() {
	app := Setup()

	log.Fatal(app.Listen(":3000"))
}

func Setup() *fiber.App {
	app := fiber.New()

	store, err := NewStore("data.json")
	if err != nil {
		log.Fatal("Cannot read DB file")
	}

	v1 := app.Group("/api/v1")

	assist := v1.Group("/assist")

	assist.Put("/:id?", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		if idStr == "" {
			idStr = "0"
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(NewError(err.Error()))
		}

		m, err := store.GitByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(NewError(err.Error()))
		}
		return c.Status(fiber.StatusOK).JSON(m)
	})

	assistantDB := v1.Group("/assistant/db")
	assistantDB.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(store.GetAll())
	})

	assistantDB.Get("/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		if idStr == "" {
			idStr = "0"
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(NewError(err.Error()))
		}

		m, err := store.GitByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(NewError(err.Error()))
		}
		return c.Status(fiber.StatusOK).JSON(m)
	})

	return app
}

// Create new Store
func NewStore(fn string) (*Store, error) {
	// Read file
	jsonBytes, err := os.ReadFile(fn)
	if err != nil || len(jsonBytes) == 0 {
		fmt.Println("Cannot read file: ", fn)
		fmt.Println("Error:", err)
		return nil, err
	}

	var msgs []Message
	// Unmarshal it to DB
	if err = json.Unmarshal(jsonBytes, &msgs); err != nil {
		fmt.Println("Cannot unmarshal json: ", fn)
		fmt.Println("Error:", err)
		return nil, err
	}
	return &Store{
		messages: msgs,
		fileName: fn,
	}, nil
}

// Create Error message
func NewError(msg string) fiber.Map {
	return fiber.Map{
		"status":  "error",
		"message": msg,
	}
}

// Get all messages
func (s *Store) GetAll() []Message {
	return s.messages
}

// Get message by ID
func (s *Store) GitByID(id int) (*Message, error) {
	for _, msg := range s.messages {
		if msg.ID == id {
			return &msg, nil
		}
	}
	return nil, errors.New("no messages found")
}
