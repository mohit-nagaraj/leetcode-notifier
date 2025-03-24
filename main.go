package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// LeetCodeResponse represents the response structure from LeetCode's GraphQL API.
type LeetCodeResponse struct {
	Data struct {
		ActiveDailyCodingChallengeQuestion struct {
			Link     string `json:"link"`
			Question struct {
				Title string `json:"title"`
			} `json:"question"`
		} `json:"activeDailyCodingChallengeQuestion"`
	} `json:"data"`
}

// Fetch the daily LeetCode problem
func fetchDailyProblem() (string, string, error) {
	query := `{
		activeDailyCodingChallengeQuestion {
			link
			question {
				title
			}
		}
	}`

	requestBody, err := json.Marshal(map[string]string{
		"query": query,
	})
	if err != nil {
		return "", "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	resp, err := http.Post("https://leetcode.com/graphql", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", "", fmt.Errorf("error making request to LeetCode: %v", err)
	}
	defer resp.Body.Close()

	// Ensure response status is OK
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected response from LeetCode: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response body: %v", err)
	}

	var leetCodeResponse LeetCodeResponse
	if err := json.Unmarshal(body, &leetCodeResponse); err != nil {
		return "", "", fmt.Errorf("error parsing LeetCode response: %v", err)
	}

	title := leetCodeResponse.Data.ActiveDailyCodingChallengeQuestion.Question.Title
	link := "https://leetcode.com" + leetCodeResponse.Data.ActiveDailyCodingChallengeQuestion.Link

	return title, link, nil
}

// Send message via Telegram Bot
func sendTelegramMessage(botToken, chatID, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	payload := map[string]string{
		"chat_id": chatID,
		"text":    message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error sending Telegram message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram API error: %s", string(body))
	}

	return nil
}

func main() {
	title, link, err := fetchDailyProblem()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	message := fmt.Sprintf("Today's LeetCode Problem: %s\n%s\nDear students, please find the daily challenge posted for today ☝️", title, link)

	botToken := "8120973833:AAF0IgNyy3AhLUrwJcZzVABBeERj_NwMpJU"
	chatID := "-1002526221482"

	if err := sendTelegramMessage(botToken, chatID, message); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("SMS sent successfully!")
}
