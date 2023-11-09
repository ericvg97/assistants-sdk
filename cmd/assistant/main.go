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
		Tools: []Tool{
			{
				Type: "function",
				Function: Function{
					Name:        "scheduleNewPickup",
					Description: "Schedules a new pickup for the customer for the next business day in the specified hour, with a 2 hour window",
					Parameters: Parameters{
						Type: "object",
						Properties: map[string]Parameter{
							"time": {
								Type:        "string",
								Description: "Time of the pickup. Format: HH:MM",
							},
						},
						Required: []string{"time"},
					},
				},
			},
		},
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
will contact them as soon as possible. You can schedulePickups and retrieve labels but nothing else.
`

type CreateAssistantRequest struct {
	Instructions string `json:"instructions"`
	Name         string `json:"name"`
	Model        string `json:"model"`
	Tools        []Tool `json:"tools"`
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Parameters  Parameters `json:"parameters"`
}

type Parameters struct {
	Type       string               `json:"type"`
	Properties map[string]Parameter `json:"properties"`
	Required   []string             `json:"required"`
}

type Parameter struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type CreateAssistantResponse struct {
	ID string `json:"id"`
}
