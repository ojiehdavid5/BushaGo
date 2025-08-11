package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"BushaGo/handler"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	// Routes
	app.Post("/recipient", handlers.CreateRecipientHandler)

	log.Fatal(app.Listen(":3001"))
}
