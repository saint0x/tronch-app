package main

import (
  "log"
  "smart_contract/pkg/email"
)

func main() {
  // Mock input
  clientFirstName := "John"
  userFirstName := "Jane"
  requirements := "- Create a login system\n- Implement payment gateway"
  paymentLink := "https://example.com/payment"
  dashboardLink := "https://example.com/dashboard"

  // Create email data
  data := email.EmailData{
    ClientFirstName: clientFirstName,
    UserFirstName:   userFirstName,
    Requirements:    requirements,
    PaymentLink:     paymentLink,
    DashboardLink:   dashboardLink,
  }

  // Generate email body and subject
  subject, emailBody, err := email.GenerateEmailBody(data)
  if err != nil {
    log.Fatalf("Error generating email body: %v", err)
  }

  // Save email body and subject to txt file
  err = email.SaveToTxt("email", subject, emailBody)
  if err != nil {
    log.Fatalf("Error saving email to file: %v", err)
  }

  log.Printf("Email subject: %s", subject)
  log.Println("Email body saved to email.txt")
  log.Println("Test was successful!")
}
