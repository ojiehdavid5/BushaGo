package handlers

import (
	"bytes"
	"encoding/json"
	 "fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Struct for incoming payload
type CreateRecipientPayload struct {
	Currency      string `json:"currency"`
	CountryCode   string `json:"country_code"`
	Type          string `json:"type"`
	BankName      string `json:"bank_name"`
	BankCode      string `json:"bank_code"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
}

// Struct for Busha recipients response
type Recipient struct {
	ID            string `json:"id"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	BankName      string `json:"bank_name"`
}

type RecipientsResponse struct {
	Data []Recipient `json:"data"`
}

func CreateRecipientHandler(c *fiber.Ctx) error {
	var payload CreateRecipientPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Busha API credentials — ideally load from .env
	apiKey := "STV0ckQwSTZaeDpPQTJJODh5RkMwbDhPaVFTV1VMVzBjRWozTWljdkVqdDhkdUVicW10andQZzcwUGE="
	profileID := "BUS_jlKUYwF9z1ynQZ98bWbaP"
	apiVersion := "2025-07-11"

	client := &http.Client{}

	// Step 1 — Get all existing recipients
	req, err := http.NewRequest("GET", "https://api.sandbox.busha.so/v1/recipients", nil)
	if err != nil {
		return err
	}

	// fmt.Println(req)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("X-BU-PROFILE-ID", profileID)
	req.Header.Set("X-BU-VERSION", apiVersion)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var recipientsList RecipientsResponse
	if err := json.Unmarshal(body, &recipientsList); err != nil {
		return err
	}
	fmt.Println(recipientsList)

	// Step 2 — Check if recipient exists
	for _, r := range recipientsList.Data {
		if r.AccountNumber == payload.AccountNumber && r.BankName == payload.BankName {
			return c.JSON(fiber.Map{
				"message": "Recipient already exists",
				"data":    r,
			})
		}
	}

	// Step 3 — Create recipient if not found
	createBody, _ := json.Marshal(payload)

	createReq, err := http.NewRequest("POST", "https://api.sandbox.busha.so/v1/recipients", bytes.NewBuffer(createBody))
	if err != nil {
		return err
	}
	createReq.Header.Set("Authorization", "Bearer "+apiKey)
	createReq.Header.Set("X-BU-PROFILE-ID", profileID)
	createReq.Header.Set("X-BU-VERSION", apiVersion)
	createReq.Header.Set("Content-Type", "application/json")

	createResp, err := client.Do(createReq)
	if err != nil {
		return err
	}
	defer createResp.Body.Close()

	createRespBody, _ := io.ReadAll(createResp.Body)

	if createResp.StatusCode != http.StatusCreated {
		return c.Status(createResp.StatusCode).SendString(string(createRespBody))
	}

	var createdRecipient map[string]interface{}
	json.Unmarshal(createRespBody, &createdRecipient)

	return c.Status(http.StatusCreated).JSON(createdRecipient)
}
