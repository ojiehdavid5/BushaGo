package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type RecipientRequest struct {
    Amount         float64        `json:"amount"`
    Currency       string         `json:"currency"`
    AccountDetails AccountDetails `json:"account_details"`
}

type AccountDetails struct {
    AccountNo   int64  `json:"account_no"`
    AccountName string `json:"account_name"`
    Bank        string `json:"bank"`
}

type PayoutPayload struct{

	
}
func CreateTransactionHandler(c *fiber.Ctx) error {

	var payload RecipientRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	fmt.Println(payload)

	RecipientID, err := CreateRecipient(CreateRecipientPayload{
		Currency:      payload.Currency,
		CountryCode:   "NG",
		Type:          "bank",
		BankName:      payload.AccountDetails.Bank,
		BankCode:      "123",
		AccountNumber: fmt.Sprintf("%d", payload.AccountDetails.AccountNo),
		AccountName:   payload.AccountDetails.AccountName,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create recipient",
		})
	}

	fmt.Println("Recipient created with ID:", RecipientID)

	fmt.Println("CreateTransactionHandler called")

	return c.SendString("Transaction created")
}
