package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// CreateRecipientPayload represents the payload for creating a recipient.
type CreateRecipientPayload struct {
	Currency      string `json:"currency"`
	CountryCode   string `json:"country_code"`
	Type          string `json:"type"`
	BankName      string `json:"bank_name"`
	BankCode      string `json:"bank_code"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
}

// Recipient represents a Busha recipient.
type Recipient struct {
	ID            string `json:"id"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	BankName      string `json:"bank_name"`
}

// RecipientsResponse is the response structure for listing recipients.
type RecipientsResponse struct {
	Data []Recipient `json:"data"`
}

// CreateRecipientResponse is the response structure for creating a recipient.
type CreateRecipientResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Data    Recipient `json:"data"`
}

// CreateRecipient checks if a recipient exists; if not, creates one and returns its ID.
func CreateRecipient(payload CreateRecipientPayload) (string, error) {
	apiKey := os.Getenv("BUSHA_API_KEY")
	profileID := os.Getenv("BUSHA_PROFILE_ID")
	apiVersion := os.Getenv("BUSHA_API_VERSION")

	if apiKey == "" || profileID == "" || apiVersion == "" {
		return "", errors.New("missing Busha API credentials in environment variables")
	}

	client := &http.Client{}

	// Step 1 — Get all existing recipients
	req, err := http.NewRequest("GET", "https://api.sandbox.busha.so/v1/recipients", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("X-BU-PROFILE-ID", profileID)
	req.Header.Set("X-BU-VERSION", apiVersion)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var recipientsList RecipientsResponse
	if err := json.Unmarshal(body, &recipientsList); err != nil {
		return "", err
	}

	// Step 2 — Check if recipient exists
	for _, r := range recipientsList.Data {
		if r.AccountNumber == payload.AccountNumber && r.BankName == payload.BankName {
			fmt.Println("Recipient already exists:", r.ID)
			return r.ID, nil
		}
	}

	// Step 3 — Create recipient if not found
	createBody, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	createReq, err := http.NewRequest("POST", "https://api.sandbox.busha.so/v1/recipients", bytes.NewBuffer(createBody))
	if err != nil {
		return "", err
	}
	createReq.Header.Set("Authorization", "Bearer "+apiKey)
	createReq.Header.Set("X-BU-PROFILE-ID", profileID)
	createReq.Header.Set("X-BU-VERSION", apiVersion)
	createReq.Header.Set("Content-Type", "application/json")

	createResp, err := client.Do(createReq)
	if err != nil {
		return "", err
	}
	defer createResp.Body.Close()

	createRespBody, err := io.ReadAll(createResp.Body)
	if err != nil {
		return "", err
	}

	if createResp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create recipient: %s", string(createRespBody))
	}

	// Unmarshal create recipient response
	var createdRecipientResp CreateRecipientResponse
	if err := json.Unmarshal(createRespBody, &createdRecipientResp); err != nil {
		return "", err
	}

	if createdRecipientResp.Data.ID == "" {
		return "", errors.New("recipient created but ID is empty in response")
	}

	fmt.Println("Recipient created successfully:", createdRecipientResp.Data.ID)
	return createdRecipientResp.Data.ID, nil
}
