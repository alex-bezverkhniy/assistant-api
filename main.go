package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

// Data Structures

// Option for user selection
type Option struct {
	gorm.Model
	ID        uint   `json:"id"`
	Body      string `json:"body"`
	MessageID int    `json:"nextMessageId"`
}

// List of Options
type Options []Option

// Message/Question
type Message struct {
	gorm.Model
	// ID      uint    `json:"id"`
	Body    string  `json:"body"`
	Options Options `json:"options"`
	FlowID  uint
}

// Flow - dialog flow
type Flow struct {
	gorm.Model
	Messages []Message
	Title    string
}

type FlowStorage struct {
	db *gorm.DB
}

func NewFlowStorage(db *gorm.DB) *FlowStorage {
	return &FlowStorage{
		db: db,
	}
}

func main() {
	db, err := SetupDB("assistant.db")
	if err != nil {
		panic(err.Error())
	}
	store := NewFlowStorage(db)
	app := Setup(store)

	log.Fatal(app.Listen(":3000"))
}

// Setup DB connection and does all migration and default data
func SetupDB(dbFile string) (*gorm.DB, error) {

	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect database")
		return nil, err
	}

	db.AutoMigrate(&Option{})
	db.AutoMigrate(&Message{})
	db.AutoMigrate(&Flow{})

	flow, err := LoadFlowFromJson("data.json")
	if err != nil {
		log.Fatalln("cannot load default data")
		return nil, err
	}

	db.Create(flow)

	return db, nil
}

func Setup(store *FlowStorage) *fiber.App {
	app := fiber.New()

	v1 := app.Group("/api/v1")

	assist := v1.Group("/assist")

	assist.Put("/:id?", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		if idStr == "" {
			idStr = "1"
		}

		id64, err := strconv.ParseUint(idStr, 10, 32)
		id := uint(id64)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(NewError(err.Error()))
		}

		m := store.GetByID(id)
		if m == nil || m.ID == 0 {
			return c.Status(fiber.StatusNotFound).JSON(NewError("no messages found"))
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

		id64, err := strconv.ParseUint(idStr, 10, 32)
		id := uint(id64)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(NewError(err.Error()))
		}

		m := store.GetByID(id)
		if m == nil || m.ID == 0 {
			return c.Status(fiber.StatusNotFound).JSON(NewError("message not found"))
		}
		return c.Status(fiber.StatusOK).JSON(m)
	})

	return app
}

// Create new dialog Flow
func LoadFlowFromJson(fn string) (*Flow, error) {
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
	return &Flow{
		Messages: msgs,
		Title:    "default",
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
func (s *FlowStorage) GetAll() []Message {
	var messages []Message
	s.db.Find(&messages)
	return messages
}

// Get message by ID
func (s *FlowStorage) GetByID(id uint) *Message {
	var msg *Message
	s.db.First(&msg, id)
	return msg
}
