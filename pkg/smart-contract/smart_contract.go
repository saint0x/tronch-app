package smart_contract

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// Represents the current stage/status of the contract
type ContractStatus string

const (
	AwaitingConfirmation ContractStatus = "awaiting_confirmation"
	ContractConfirmed    ContractStatus = "contract_confirmed"
	PaymentMade          ContractStatus = "payment_made"
	ReqsCompleted        ContractStatus = "reqs_completed"
	ContractExecuted     ContractStatus = "contract_executed"
	PaymentReleased      ContractStatus = "payment_released"
)

// Initiates the contract and calculates payment amount minus our fee
func InitiateContract(paymentAmount float64) (float64, error) {
	// Calculate payment amount minus our fee (e.g., 5% fee)
	ourFee := 0.05 * paymentAmount
	netPaymentAmount := paymentAmount - ourFee

	return netPaymentAmount, nil
}

// Retrieves the current stage/status of the contract
func GetCurrentContractStatus() ContractStatus {
	// TODO: Implement logic to retrieve the current status from the contract
	return AwaitingConfirmation
}

// UpdateContractStatus updates the contract's current stage/status
func UpdateContractStatus(newStatus ContractStatus) error {
	// TODO: Implement logic to update the contract's status
	return nil
}

// Uses GPT-3.5 to extract requirements from user-provided parameters
func ExtractRequirements(ctx context.Context, clientName, clientEmail string, paymentAmount float64, requirements, description string) (string, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	if client == nil {
		return "", fmt.Errorf("failed to create OpenAI client")
	}

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("Extract the requirements for a Solidity smart contract for an escrow service based on the following details:\n\nClient Name: %s\nClient Email: %s\nPayment Amount: %.2f ETH\nUser Requirements: %s\nDescription: %s", clientName, clientEmail, paymentAmount, requirements, description),
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	log.Println("Successfully extracted requirements.")
	return resp.Choices[0].Message.Content, nil
}

// Generates a Solidity smart contract based on extracted requirements
func GenerateSmartContract(userInput map[string]string) (string, error) {
	// Read the Solidity contract template from solidity_template.json
	templateJSON, err := os.ReadFile("solidity_template.json")
	if err != nil {
		return "", fmt.Errorf("failed to read solidity_template.json: %v", err)
	}

	var templateData map[string]interface{}
	if err := json.Unmarshal(templateJSON, &templateData); err != nil {
		return "", fmt.Errorf("failed to unmarshal solidity_template.json: %v", err)
	}

	// Extract the contract template from the JSON data
	contractTemplate, ok := templateData["contractTemplate"].(string)
	if !ok {
		return "", fmt.Errorf("contractTemplate not found in solidity_template.json")
	}

	// Populate the contract template with user's input
	populatedContract := strings.ReplaceAll(contractTemplate, "{{clientName}}", userInput["clientName"])
	populatedContract = strings.ReplaceAll(populatedContract, "{{clientEmail}}", userInput["clientEmail"])
	populatedContract = strings.ReplaceAll(populatedContract, "{{paymentAmount}}", userInput["paymentAmount"])
	populatedContract = strings.ReplaceAll(populatedContract, "{{requirements}}", userInput["requirements"])
	populatedContract = strings.ReplaceAll(populatedContract, "{{description}}", userInput["description"])

	log.Println("Successfully generated smart contract.")
	return populatedContract, nil
}

// Saves the generated contract as a .sol file
func SaveContractAsSolidity(contractCode string, filePath string) error {
	solidityContent := fmt.Sprintf("/* SPDX-License-Identifier: MIT */\npragma solidity ^0.8.0;\n\n%s", contractCode)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(solidityContent)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully saved contract to %s\n", filePath)
	return nil
}
