package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"go_tgbot/config"
	"log"
	"net/http"
	"net/url"
)

func GeminiAsk(contentRequest string) (string, error) {
	geminiEndpoint := "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent"

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "USER",
				"parts": []map[string]interface{}{
					{
						"text": contentRequest,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		return "", err
	}

	req, err := http.NewRequest("POST", geminiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		return "", err
	}

	req.Header.Set("x-goog-api-key", config.SetConfig.GeminiApiKey)
	req.Header.Set("Content-Type", "application/json")

	proxyURL, _ := url.Parse("http://localhost:7890")
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	// Make the HTTP request
	response, err := client.Do(req)
	if err != nil {
		log.Println("Error making HTTP request:", err)
		return "", err
	}
	defer response.Body.Close()

	// Read response body
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		log.Println("Error decoding JSON response:", err)
		return "", err
	}

	candidatesSlice, ok := result["candidates"].([]interface{})
	if !ok || len(candidatesSlice) == 0 {
		return "", errors.New("invalid candidates in JSON response")
	}

	content, ok := candidatesSlice[0].(map[string]interface{})["content"]
	if !ok {
		return "", errors.New("content not found in JSON response")
	}

	parts, ok := content.(map[string]interface{})["parts"]
	if !ok {
		return "", errors.New("parts not found in JSON response")
	}

	partsSlice, ok := parts.([]interface{})
	if !ok || len(partsSlice) == 0 {
		return "", errors.New("invalid parts in JSON response")
	}

	text, ok := partsSlice[0].(map[string]interface{})["text"]
	if !ok {
		return "", errors.New("text not found in JSON response")
	}

	textStr, ok := text.(string)
	if !ok {
		return "", errors.New("invalid text in JSON response")
	}
	return textStr, nil
}
