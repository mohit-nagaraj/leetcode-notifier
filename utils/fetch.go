package utils

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/mohit-nagaraj/leetcode-notifier/types"
)

// FetchDailyProblem fetches the daily LeetCode problem
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
	// link := "https://fast-leetcode.vercel.app" + leetCodeResponse.Data.ActiveDailyCodingChallengeQuestion.Link
	
	return title, link, nil
}

// FetchEasyProblemOfTheDay fetches the easy problem of the day from LeetCode
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

// FetchRandomCodeChefProblem fetches a random CodeChef problem from the CSV file
func FetchRandomCodeChefProblem(csvPath string) (types.CodeChefProblem, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return types.CodeChefProblem{}, fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return types.CodeChefProblem{}, fmt.Errorf("error reading CSV file: %v", err)
	}

	if len(records) <= 1 {
		return types.CodeChefProblem{}, fmt.Errorf("CSV file is empty or has no data rows")
	}

	// Skip header row and select a random problem
	dataRows := records[1:]
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(dataRows))
	selectedRow := dataRows[randomIndex]

	if len(selectedRow) < 4 {
		return types.CodeChefProblem{}, fmt.Errorf("invalid CSV row format")
	}

	problem := types.CodeChefProblem{
		Name:       selectedRow[0],
		Link:       selectedRow[1],
		Difficulty: selectedRow[2],
		Category:   selectedRow[3],
	}

	return problem, nil
}
