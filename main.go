package main

import (
	"log"

	handlers "BushaGo/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {

		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	// Routes

	// app.Post("/recipient", handlers.CreateRecipientHandler)
	app.Post("/transaction", handlers.CreateTransactionHandler)

	log.Fatal(app.Listen(":3001"))
}
