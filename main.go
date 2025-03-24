package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

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

func fetchDailyProblem() (string, string, error) {
	query := `{
		activeDailyCodingChallengeQuestion {
			link
			question {
				title
			}
		}
	}`

	requestBody, _ := json.Marshal(map[string]string{
		"query": query,
	})

	resp, err := http.Post("https://leetcode.com/graphql", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", "", fmt.Errorf("error making request to LeetCode: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var leetCodeResponse LeetCodeResponse
	if err := json.Unmarshal(body, &leetCodeResponse); err != nil {
		return "", "", fmt.Errorf("error parsing LeetCode response: %v", err)
	}

	title := leetCodeResponse.Data.ActiveDailyCodingChallengeQuestion.Question.Title
	link := "https://leetcode.com" + leetCodeResponse.Data.ActiveDailyCodingChallengeQuestion.Link

	return title, link, nil
}

func sendSMS(phoneNumber, message string) error {
	values := url.Values{
		"phone":   {phoneNumber},
		"message": {message},
		"key":     {"912980295454a1f77c627ad3631f4fdf7f679bc4TKapiAVSXvpHHrLcMMT2cihp9_test"},
	}

	resp, err := http.PostForm("https://textbelt.com/text", values)
	if err != nil {
		return fmt.Errorf("error sending SMS: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("TextBelt Response:", string(body))
	fmt.Printf("message: %v\n", message)
	return nil
}

func main() {
	title, link, err := fetchDailyProblem()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	message := fmt.Sprintf("Today's LeetCode Problem: %s\n%s", title, link)
	phoneNumber := "+916363988392"

	if err := sendSMS(phoneNumber, message); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("SMS sent successfully!")
}
