package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
	sysPromptPath := "./sysPrompt.txt"

	data, err := os.ReadFile(sysPromptPath)
	if err != nil {
		fmt.Println("no system promt file found", err)
	}

	sysPrompt := string(data)

	url := "https://api.perplexity.ai/chat/completions"

	reqBody := ChatRequest{
		Model: "sonar",
		Messages: []Message{
			{Role: "system", Content: sysPrompt},
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

	req.Header.Set("Authorization", "Bearer "+os.Getenv("PPLX_API_KEY"))
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
	App := app.New()
	Window := App.NewWindow("FactCheck")

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter fact")

	output := widget.NewLabel("")
	output.Wrapping = fyne.TextWrapWord

	scroll := container.NewScroll(output)
	scroll.SetMinSize(fyne.NewSize(300, 200))

	button:= widget.NewButton("check", func() {
		if entry.Text == "" {
			output.SetText("no string found")
			return
		}

		output.SetText("Loading...")

		go func() {
			result, err := FactCheck(entry.Text)

			fyne.Do(func() {
				if err != nil {
					output.SetText("Bad request")
				} else {
					output.SetText(result)
				}
			})
		}()
	})

	content := container.NewVBox(
		entry,
		button,
		output,
	)
	Window.SetContent(content)
	Window.ShowAndRun()
}
