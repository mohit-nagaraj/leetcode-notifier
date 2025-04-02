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

// Fetch problems based on specific criteria and return one problem based on the current date
func FetchEasyProblemOfTheDay() (string, string, error) {
	query := `{
        problemsetQuestionList: questionList(
            categorySlug: "all-code-essentials"
            limit: 50
            skip: 0
            filters: {
                difficulty: EASY
                tags: ["array", "string"]
				paidOnly: false
            }
        ) {
            total: totalNum
            questions: data {
                acRate
                difficulty
                title
                titleSlug
				paidOnly
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
		fmt.Println(resp)
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

	questions := leetCodeResponse.Data.ProblemsetQuestionList.Questions
	if len(questions) == 0 {
		return "", "", fmt.Errorf("no problems found")
	}

	startDate := time.Date(2025, 3, 29, 0, 0, 0, 0, time.UTC)
	today := time.Now().UTC()

	daysSinceStart := int(today.Sub(startDate).Hours() / 24)

	index := daysSinceStart % len(questions)
	selectedProblem := questions[index]

	title := selectedProblem.Title
	link := "https://leetcode.com/problems/" + selectedProblem.TitleSlug

	return title, link, nil
}
