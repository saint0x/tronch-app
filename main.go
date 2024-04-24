package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "text/template"

  "smart_contract/pkg/smart-contract"
)

type ContractData struct {
  ClientFirstName string `json:"client_first_name"`
  UserFirstName   string `json:"user_first_name"`
  Requirements    string `json:"requirements"`
  Description     string `json:"description"`
}

func main() {
  // Serve static files
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

  // Handle API endpoints
  http.HandleFunc("/generate_contract", GenerateContract)

  // Serve index.html as the frontend
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
      http.Error(w, "Internal Server Error", http.StatusInternalServerError)
      return
    }
    tmpl.Execute(w, nil)
  })

  fmt.Println("Server is running on port 8080...")
  http.ListenAndServe(":8080", nil)
}

func GenerateContract(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Printf("Error reading request body: %v", err)
    http.Error(w, "Failed to read request body", http.StatusInternalServerError)
    return
  }

  // Check if the request body is empty
  if len(body) == 0 {
    log.Println("Request body is empty")
    http.Error(w, "Request body is empty", http.StatusBadRequest)
    return
  }

  var data ContractData
  err = json.Unmarshal(body, &data)
  if err != nil {
    log.Printf("Error decoding request body: %v", err)
    http.Error(w, "Failed to decode request body", http.StatusBadRequest)
    return
  }

  log.Printf("Received contract data: %+v", data)

  // Extract requirements
  ctx := r.Context() // You can pass context if needed
  requirements, err := smart_contract.ExtractRequirements(ctx, data.ClientFirstName, "", 0, data.Requirements, data.Description)
  if err != nil {
    log.Printf("Error extracting requirements: %v", err)
    http.Error(w, "Failed to extract requirements", http.StatusInternalServerError)
    return
  }

  log.Printf("Extracted requirements: %s", requirements)

  // Generate the smart contract
  userInput := map[string]string{
    "requirements": requirements,
    "description":  data.Description,
  }
  contractCode, err := smart_contract.GenerateSmartContract(userInput)
  if err != nil {
    log.Printf("Error generating smart contract: %v", err)
    http.Error(w, "Failed to generate smart contract", http.StatusInternalServerError)
    return
  }

  // Save the generated contract as a .sol file
  err = smart_contract.SaveContractAsSolidity(contractCode, "contract.sol")
  if err != nil {
    log.Printf("Error saving smart contract: %v", err)
    http.Error(w, "Failed to save smart contract", http.StatusInternalServerError)
    return
  }

  w.Write([]byte("Contract generated and saved successfully"))
}
