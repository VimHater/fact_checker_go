package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func FactCheck(fact string) (string, error) {
	url := "https://api.perplexity.ai/chat/completions"

	reqBody := ChatRequest{
		Model: "sonar",
		Messages: []Message{
			{Role: "system", Content: "You are a fact-checking assistant. Verify claims and cite evidence if possible."},
			{Role: "user", Content: fact},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer " + os.Getenv("PPLX_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	respBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBytes, &chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) > 0 {
		return chatResp.Choices[0].Message.Content, nil
	}

	return "No response from model", nil
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Usage: need arguments")
		return	
	}

	result, err := FactCheck(args[0])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Fact-check result:\n", result)
}

