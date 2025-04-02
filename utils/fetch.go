package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mohit-nagaraj/leetcode-notifier/types"
)

// Fetch the daily LeetCode problem
func FetchDailyProblem() (string, string, error) {
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

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected response from LeetCode: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response body: %v", err)
	}

	var leetCodeResponse types.LeetCodeResponse
	if err := json.Unmarshal(body, &leetCodeResponse); err != nil {
		return "", "", fmt.Errorf("error parsing LeetCode response: %v", err)
	}

	title := leetCodeResponse.Data.ActiveDailyCodingChallengeQuestion.Question.Title
	link := "https://leetcode.com" + leetCodeResponse.Data.ActiveDailyCodingChallengeQuestion.Link

	return title, link, nil
}

// Fetch the easy problem of the day from LeetCode
func FetchEasyProblemOfTheDay() (string, string, error) {
	query := `{
        problemsetQuestionList: questionList(
            categorySlug: "all-code-essentials"
            limit: 100
            skip: 0
            filters: {
                difficulty: EASY
                tags: ["array", "string"]
                premiumOnly: false
            }
        ) {
            total: totalNum
            questions: data {
                acRate
                difficulty
                title
                titleSlug
                isPaidOnly
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

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected response from LeetCode: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response body: %v", err)
	}

	var leetCodeResponse types.EasyProblemsResponse
	if err := json.Unmarshal(body, &leetCodeResponse); err != nil {
		return "", "", fmt.Errorf("error parsing LeetCode response: %v", err)
	}

	// Filter out any premium problems that might have slipped through
	var nonPremiumQuestions []types.Problem
	for _, question := range leetCodeResponse.Data.ProblemsetQuestionList.Questions {
		if !question.IsPaidOnly {
			nonPremiumQuestions = append(nonPremiumQuestions, question)
		}
	}

	if len(nonPremiumQuestions) == 0 {
		return "", "", fmt.Errorf("no non-premium problems found")
	}

	// Select problem based on current date
	startDate := time.Date(2025, 3, 29, 0, 0, 0, 0, time.UTC)
	today := time.Now().UTC()
	daysSinceStart := int(today.Sub(startDate).Hours() / 24)
	index := daysSinceStart % len(nonPremiumQuestions)
	selectedProblem := nonPremiumQuestions[index]

	title := selectedProblem.Title
	link := "https://leetcode.com/problems/" + selectedProblem.TitleSlug
	fmt.Printf("Selected problem: %s\n", link)
	return title, link, nil
}
