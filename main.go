package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
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

// Send SMS via TextBelt
func sendSMS(phoneNumber, message string) error {
	apiKey := os.Getenv("TEXTBELT_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("TEXTBELT_API_KEY environment variable not set")
	}

	values := url.Values{
		"phone":   {phoneNumber},
		"message": {message},
		"key":     {apiKey},
	}

	resp, err := http.PostForm("https://textbelt.com/text", values)
	if err != nil {
		return fmt.Errorf("error sending SMS: %v", err)
	}
	defer resp.Body.Close()

	// Ensure response status is OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response from TextBelt: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading TextBelt response: %v", err)
	}

	fmt.Println("TextBelt Response:", string(body))
	return nil
}

func main() {
	title, link, err := fetchDailyProblem()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	message := fmt.Sprintf("Today's LeetCode Problem: %s\n%s", title, link)

	phoneNumber := os.Getenv("PHONE_NUMBER")
	if phoneNumber == "" {
		phoneNumber = "+916363988392" // Default fallback
	}

	if err := sendSMS(phoneNumber, message); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("SMS sent successfully!")
}
