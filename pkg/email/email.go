package email

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
)

// Represents the data to populate in the email template
type EmailData struct {
	ClientFirstName string
	UserFirstName   string
	Requirements    string
	PaymentLink     string
	DashboardLink   string
}

// Generates the email body and subject based on the provided data
func GenerateEmailBody(data EmailData) (string, string, error) {
	// Define email subject and body
	subject := "[tronch.io] Your escrow has initiated!"
	body := `
Congrats {{.ClientFirstName}},

{{.UserFirstName}} has initiated an escrow for your project! 🎉

We’re tronch, the escrow service that will be handling payment!
You can rest assured knowing your funds are completely secure and will be released only once the project milestones have been completed to your standards.

The requirements that are set to be programmed into your smart contract are as follows:
{{.Requirements}}

Below, you’ll find two links:

- A secure link to make payment: 

{{.PaymentLink}}

- A link to your collaborative Dashboard, where you will confirm requirements, and check on project updates from a birds-eye view. We’ll email you as the project status changes:

{{.DashboardLink}}

Next steps would be to head over to the Dashboard, confirm the requirements, make payment, and leave the rest to us!
`

	// Create a new buffer for the email body
	var bodyBytes bytes.Buffer

	// Create a new template and parse the body
	tmpl, err := template.New("email").Parse(body)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse email template: %v", err)
	}

	// Populate the template with data
	err = tmpl.Execute(&bodyBytes, data)
	if err != nil {
		return "", "", fmt.Errorf("failed to execute email template: %v", err)
	}

	// Convert the buffer to string
	bodyStr := bodyBytes.String()

	return subject, bodyStr, nil
}

// Saves the given subject and content to a .txt file
func SaveToTxt(filename, subject, content string) error {
	// Combine subject and content with two spaces in between
	fullContent := fmt.Sprintf("%s\n\n%s", subject, content)

	err := os.WriteFile(filename+".txt", []byte(fullContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to save content to .txt file: %v", err)
	}
	return nil
}
