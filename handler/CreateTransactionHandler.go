package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

// RecipientRequest represents the incoming request body for a transaction.
type RecipientRequest struct {
	Amount         string         `json:"amount"`
	Currency       string         `json:"currency"`
	AccountDetails AccountDetails `json:"account_details"`
}



// AccountDetails holds bank account details for payouts.
type AccountDetails struct {
	AccountNo   string `json:"account_no"`
	AccountName string `json:"account_name"`
	Bank        string `json:"bank"`
}

// PayoutPayload represents the structure for creating a payout.
type PayoutPayload struct {
	SourceCurrency string     `json:"source_currency"`
	TargetCurrency string     `json:"target_currency"`
	SourceAmount   float64    `json:"source_amount"`
	PayIn          PayMethod  `json:"pay_in"`
	PayOut         PayOutData `json:"pay_out"`
}

// PayMethod describes the method of payment.
type PayMethod struct {
	Type string `json:"type"`
}

// PayOutData describes payout recipient details.
type PayOutData struct {
	Type        string `json:"type"`
	RecipientID string `json:"recipient_id"`
}

// QuoteResponse is the response from Busha's quote API.
type QuoteResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID string `json:"id"`
	} `json:"data"`
}

// CreateTransactionHandler handles creation of recipient, quote, and returns the quote ID.
func CreateTransactionHandler(c *fiber.Ctx) error {
	var payload RecipientRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Step 1: Create recipient
	recipientID, err := CreateRecipient(CreateRecipientPayload{
		Currency:      "NGN",
		CountryCode:   "NG",
		Type:          "ngn_bank",
		BankName:      payload.AccountDetails.Bank,
		BankCode:      "100004", // hardcoded, should ideally be dynamic
		AccountNumber: payload.AccountDetails.AccountNo,
		AccountName:   payload.AccountDetails.AccountName,
	})
	if err != nil {
		fmt.Println("Recipient creation error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create recipient",
		})
	}
	fmt.Println("Recipient created with ID:", recipientID)

	// Step 2: Create quote
	quotePayload := PayoutPayload{
		SourceCurrency: "USDT",
		TargetCurrency: "NGN",
		SourceAmount:   10,
		PayIn:          PayMethod{Type: "balance"},
		PayOut:         PayOutData{Type: "bank_transfer", RecipientID: recipientID},
	}

	quoteBody, err := json.Marshal(quotePayload)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to marshal quote payload",
		})
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.sandbox.busha.so/v1/quotes", bytes.NewBuffer(quoteBody))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create quote request",
		})
	}

	apiKey := os.Getenv("BUSHA_API_KEY")
	profileID := os.Getenv("BUSHA_PROFILE_ID")
	apiVersion := os.Getenv("BUSHA_API_VERSION")

	if apiKey == "" || profileID == "" || apiVersion == "" {
		fmt.Println("Missing Busha API credentials in environment variables")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send quote request",
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return c.Status(resp.StatusCode).SendString(string(body))
	}

	var quoteResp QuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quoteResp); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse quote response",
		})
	}

	// Step 3: Return quote ID
	return c.JSON(fiber.Map{
		"status":   "success",
		"quote_id": quoteResp.Data.ID,
	})
}
