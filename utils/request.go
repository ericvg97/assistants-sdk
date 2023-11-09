package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func DoRequest(requestURL string, requestBody interface{}, method string) *http.Response {
	// fmt.Println(requestURL)

	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(reqBody))

	if method == "GET" {
		reqBody = nil
	}

	req, err := http.NewRequest(method, requestURL, bytes.NewBuffer(reqBody))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openAISecretKey)
	req.Header.Set("OpenAI-Beta", "assistants=v1")

	// command, _ := http2curl.GetCurlCommand(req)
	// fmt.Println(command)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		newStr := buf.String()
		fmt.Println(newStr)
		panic("Request failed")
	}

	return resp
}
