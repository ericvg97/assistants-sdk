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
const assistantID = "asst_tz3woggpuACf21pADXtnpLfh"

type CreateThreadRespones struct {
	ID string `json:"id"`
}

func main() {
	threadID := CreateThread()
	// fmt.Println("ThreadID:", threadID)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		message := scanner.Text()
		done := make(chan bool)
		go func() {
			for {
				for _, r := range `-\|/` {
					select {
					case <-done:
						return
					default:
						fmt.Printf("\rClaudia Möller: %c", r)
						time.Sleep(100 * time.Millisecond)
					}
				}
			}
		}()

		CreateMessage(message, threadID)
		response := RunThread(threadID)
		done <- true

		fmt.Printf("\rClaudia Möller: %s\n", response)
	}

}

func CreateThread() string {
	requestURL := openaiApiUrl + createThreadSuffix

	resp := utils.DoRequest(requestURL, nil, "POST")

	var createThreadResponse CreateThreadRespones
	err := json.NewDecoder(resp.Body).Decode(&createThreadResponse)
	if err != nil {
		panic(err)
	}

	return createThreadResponse.ID
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

type ThreadResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func RunThread(threadID string) string {
	requestURL := fmt.Sprintf(openaiApiUrl+runThreadSuffix, threadID)
	requestBody := RunThreadRequest{
		AssistantID: assistantID,
	}

	resp := utils.DoRequest(requestURL, requestBody, "POST")
	var runThreadResponse ThreadResponse
	err := json.NewDecoder(resp.Body).Decode(&runThreadResponse)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("Thread '%v' is running on %v with status %v", threadID, runThreadResponse.ID, runThreadResponse.Status)
	PollRun(threadID, runThreadResponse.ID)
	return GetLastMessage(threadID)
}

func PollRun(threadID, runID string) {
	reqURL := fmt.Sprintf(openaiApiUrl+retrieveRunSuffix, threadID, runID)

	status := "queued"
	for status == "queued" || status == "in_progress" {
		// fmt.Println("Sleeping as status is", status)
		time.Sleep(1 * time.Second)

		resp := utils.DoRequest(reqURL, nil, "GET")

		var threadResponse ThreadResponse
		err := json.NewDecoder(resp.Body).Decode(&threadResponse)
		if err != nil {
			panic(err)
		}

		status = threadResponse.Status
	}

	// fmt.Printf("Thread finished with status %v \n", status)
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
