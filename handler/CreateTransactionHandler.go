package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type RecipientRequest struct {
    Amount         string        `json:"amount"`
    Currency       string         `json:"currency"`
    AccountDetails AccountDetails `json:"account_details"`
}

type AccountDetails struct {
    AccountNo   string  `json:"account_no"`
    AccountName string `json:"account_name"`
    Bank        string `json:"bank"`
}

type PayoutPayload struct{
SourceCurrency string     `json:"source_currency"`
    TargetCurrency string     `json:"target_currency"`
    SourceAmount   float64    `json:"source_amount"`
    PayIn          PayMethod  `json:"pay_in"`
    PayOut         PayOutData `json:"pay_out"`

}

type PayMethod struct {
    Type string `json:"type"`
}

type PayOutData struct {
    Type        string `json:"type"`
    RecipientID string `json:"recipient_id"`
} 

func CreateTransactionHandler(c *fiber.Ctx) error {

	var payload RecipientRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	RecipientID, err := CreateRecipient(CreateRecipientPayload{
		Currency:      "NGN",
		CountryCode:   "NG",
		Type:          "bank",
		BankName:      payload.AccountDetails.Bank,
		BankCode:      "100004",
		AccountNumber:  payload.AccountDetails.AccountNo,
		AccountName:   payload.AccountDetails.AccountName,
	})

	fmt.Println(err)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create recipient",
		})
	}

	fmt.Println("Recipient created with ID:", RecipientID)

	fmt.Println("CreateTransactionHandler called")

	return c.SendString("Transaction created")
}
