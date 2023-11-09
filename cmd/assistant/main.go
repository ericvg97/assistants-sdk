package main

import (
	"encoding/json"

	"github.com/ericvg97/assistants-sdk/utils"
)

const openaiApiUrl = "https://api.openai.com/v1"
const createAssistantSuffix = "/assistants"

func main() {
	requestBody := CreateAssistantRequest{
		Instructions: prompt,
		Name:         "Customer Support Rever",
		Model:        "gpt-4",
	}

	requestURL := openaiApiUrl + createAssistantSuffix

	resp := utils.DoRequest(requestURL, requestBody, "POST")

	var createAssistantResponse CreateAssistantResponse
	err := json.NewDecoder(resp.Body).Decode(&createAssistantResponse)
	if err != nil {
		panic(err)
	}

	println(createAssistantResponse.ID)
}

var prompt = `
You are a customer success/support assistant for REVER, a startup from Barcelona. REVER manages the returns of ecommerces and provides customer support for the reverse logistic
process. Customers are going to ask you questions about different topics. If you don't know the answer, ask them for their contact details (phone or email) and tell them someone
will contact them as soon as possible. Be always nice but not overwhelmingly nice. Try to give short responses and add questions for clarifications. Do not overwhelm the customer.
If someone asks you about a label tell them to:
1. Check SPAM
2. Check other tabs like promotions
3. Check they are looking at the correct email account
4. Wait for at least 2 days
5. Contact ops@itsrever.com
`

type CreateAssistantRequest struct {
	Instructions string `json:"instructions"`
	Name         string `json:"name"`
	Model        string `json:"model"`
}

type CreateAssistantResponse struct {
	ID string `json:"id"`
}
