package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ericvg97/assistants-sdk/utils"
)

const openaiApiUrl = "https://api.openai.com/v1"

const createThreadSuffix = "/threads"
const createMessageSuffix = "/threads/%v/messages"
const runThreadSuffix = "/threads/%v/runs"
const retrieveRunSuffix = "/threads/%v/runs/%v"
const listMessagesSuffix = "/threads/%v/messages"
const submitToolCallSuffix = "/threads/%v/runs/%v/submit_tool_outputs"

const assistantID = "asst_IcCi26qqxaTFYR5X4cKArjC4"

type Thread struct {
	ID string `json:"id"`
}

const returnProcess = `
this is the ReturnProcessState
	RefundTiming:         "ON_SENT",
	RefundedAt:     "",
	LogisticStatus: "SENT",
	RefundPaymentMethod: ORIGINAL_PAYMENT_METHOD
just answer this message with OK 
	`

func main() {
	threadID := CreateThread()

	scanner := bufio.NewScanner(os.Stdin)
	CreateMessage(returnProcess, threadID)

	for {
		fmt.Printf("%s", ColorReset)

		scanner.Scan()
		message := scanner.Text()

		fmt.Printf("%s", ColorGrey)

		CreateMessage(message, threadID)
		response := RunThread(threadID)

		fmt.Printf("%sClaudia Möller: %s%s\n", ColorBlue, response, ColorGrey)
	}
}

func CreateThread() string {
	requestURL := openaiApiUrl + createThreadSuffix

	resp := utils.DoRequest(requestURL, nil, "POST")

	var thread Thread
	err := json.NewDecoder(resp.Body).Decode(&thread)
	if err != nil {
		panic(err)
	}

	return thread.ID
}

type CreateMessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func CreateMessage(message string, threadID string) {
	requestURL := fmt.Sprintf(openaiApiUrl+createMessageSuffix, threadID)
	requestBody := CreateMessageRequest{
		Role:    "user",
		Content: message,
	}

	utils.DoRequest(requestURL, requestBody, "POST")

	// fmt.Printf("Message '%v' sent to thread '%v'\n", message, threadID)
}

type RunThreadRequest struct {
	AssistantID string `json:"assistant_id"`
}

func RunThread(threadID string) string {
	requestURL := fmt.Sprintf(openaiApiUrl+runThreadSuffix, threadID)
	requestBody := RunThreadRequest{
		AssistantID: assistantID,
	}

	resp := utils.DoRequest(requestURL, requestBody, "POST")
	var run Run
	err := json.NewDecoder(resp.Body).Decode(&run)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("Thread '%v' is running on %v with status %v", threadID, run.ID, run.Status)
	return PollRun(threadID, run.ID)
}

func PollRun(threadID, runID string) string {
	run := Run{}

	reqURL := fmt.Sprintf(openaiApiUrl+retrieveRunSuffix, threadID, runID)

	status := "queued"
	for status == "queued" || status == "in_progress" {
		fmt.Println("Sleeping as status is", status)
		time.Sleep(2 * time.Second)

		resp := utils.DoRequest(reqURL, nil, "GET")

		err := json.NewDecoder(resp.Body).Decode(&run)
		if err != nil {
			panic(err)
		}

		status = run.Status
	}

	fmt.Printf("Thread finished with status %v \n", status)

	if run.Status == "completed" {
		return GetLastMessage(threadID)
	} else if run.Status == "requires_action" {
		return HandleAction(run, threadID)
	}

	panic("Thread finished with status " + run.Status)
}

func HandleAction(run Run, threadID string) string {
	action := run.RequiredAction
	if action.Type != "submit_tool_outputs" {
		panic("Action type not supported: " + action.Type)
	}

	if len(action.SubmitToolOutputs.ToolCalls) > 1 {
		panic("Multiple tool calls not supported")
	}

	toolCall := action.SubmitToolOutputs.ToolCalls[0]
	if toolCall.Type != "function" {
		panic("Tool call type not supported: " + toolCall.Type)
	}

	function := toolCall.Function
	if function.Name != "scheduleNewPickup" {
		panic("Function name not supported: " + function.Name)
	}

	var arguments map[string]string

	err := json.Unmarshal([]byte(function.Arguments), &arguments)
	if err != nil {
		panic("Could not unmarshal arguments")
	}

	response := scheduleNewPickup(arguments["time"])
	responseJSON, err := json.Marshal(response)
	if err != nil {
		panic("Could not marshal response")
	}

	SubmitToolCall(threadID, run.ID, toolCall.ID, string(responseJSON))
	return PollRun(threadID, run.ID)
}

type SubmitToolCallRequest struct {
	ToolOutputs []ToolOutput `json:"tool_outputs"`
}

type ToolOutput struct {
	ToolCallID string `json:"tool_call_id"`
	Output     string `json:"output"`
}

func SubmitToolCall(threadID, runID, toolCallID, outputJSON string) {
	requestURL := fmt.Sprintf(openaiApiUrl+submitToolCallSuffix, threadID, runID)
	requestBody := SubmitToolCallRequest{
		ToolOutputs: []ToolOutput{
			{
				ToolCallID: toolCallID,
				Output:     outputJSON,
			},
		},
	}

	utils.DoRequest(requestURL, requestBody, "POST")
}

func scheduleNewPickup(time string) PickupResponse {
	fmt.Printf("%sAuthenticated call to schedule pickup: %s%s\n", ColorRed, time, ColorGrey)
	return PickupResponse{
		HttpCode: 200,
		Message:  "Pickup scheduled",
	}
}

type PickupResponse struct {
	HttpCode int    `json:"http_code"`
	Message  string `json:"message"`
}

type Messages struct {
	Messages []Message `json:"data"`
}

type Message struct {
	Content []Content `json:"content"`
}

type Content struct {
	Text Text `json:"text"`
}

type Text struct {
	Value string `json:"value"`
}

func GetLastMessage(threadID string) string {
	requestURL := fmt.Sprintf(openaiApiUrl+listMessagesSuffix, threadID)

	resp := utils.DoRequest(requestURL, nil, "GET")

	var messages Messages
	err := json.NewDecoder(resp.Body).Decode(&messages)
	if err != nil {
		panic(err)
	}

	return messages.Messages[0].Content[0].Text.Value
}

type Run struct {
	ID             string         `json:"id"`
	Status         string         `json:"status"`
	RequiredAction RequiredAction `json:"required_action"`
}

type RequiredAction struct {
	Type              string            `json:"type"`
	SubmitToolOutputs SubmitToolOutputs `json:"submit_tool_outputs"`
}

type SubmitToolOutputs struct {
	ToolCalls []ToolCall `json:"tool_calls"`
}

type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorReset  = "\033[0m"
	ColorGrey   = "\033[90m"
)
