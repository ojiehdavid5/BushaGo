package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	// "log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type RecipientRequest struct {
	Amount         string         `json:"amount"`
	Currency       string         `json:"currency"`
	AccountDetails AccountDetails `json:"account_details"`
}

type AccountDetails struct {
	AccountNo   string `json:"account_no"`
	AccountName string `json:"account_name"`
	Bank        string `json:"bank"`
}

type PayoutPayload struct {
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

type QuoteResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID string `json:"id"`
	} `json:"data"`
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
		Type:          "ngn_bank",
		BankName:      payload.AccountDetails.Bank,
		BankCode:      "100004",
		AccountNumber: payload.AccountDetails.AccountNo,
		AccountName:   payload.AccountDetails.AccountName,
	})

	fmt.Println(err)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create recipient",
		})
	}

	fmt.Println("Recipient created with ID:", RecipientID)

	// Step 2: Create quote
	quotePayload := PayoutPayload{
		SourceCurrency: "USDT",
		TargetCurrency: "NGN",
		SourceAmount:   10,
		PayIn:          PayMethod{Type: "balance"},
		PayOut:         PayOutData{Type: "bank_transfer", RecipientID: RecipientID},
	}

	quoteBody, _ := json.Marshal(quotePayload)
	fmt.Println(bytes.NewBuffer(quoteBody))

	// Replace with your env variable
	client := &http.Client{}

		req, err := http.NewRequest("POST", "https://api.sandbox.busha.so/v1/quotes", bytes.NewBuffer(quoteBody))

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create quote request"})
	} 

		apiKey := os.Getenv("BUSHA_API_KEY")
	profileID := os.Getenv("BUSHA_PROFILE_ID")
	apiVersion := os.Getenv("BUSHA_API_VERSION")

		if apiKey == "" || profileID == "" || apiVersion == "" {
		 fmt.Println("missing Busha API credentials in environment variables")

	}


	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	// req.Header.Set("X-BU-PROFILE-ID", profileID)
	// req.Header.Set("X-BU-VERSION", apiVersion)

	resp, err := client.Do(req)
	if err != nil {
		// return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send quote request"})
				return  err

	}
	fmt.Println(resp)
fmt.Println(err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return c.Status(resp.StatusCode).SendString(string(body))
	}

	var quoteResp QuoteResponse
	fmt.Println(quoteResp)
	if err := json.NewDecoder(resp.Body).Decode(&quoteResp); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse quote response"})
	}

	// Step 3: Return quote ID
	return c.JSON(fiber.Map{
		"status":   "success",
		"quote_id": quoteResp.Data.ID,
	})
}
