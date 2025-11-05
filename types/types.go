package types

// LeetCodeResponse represents the response from LeetCode's GraphQL API for daily problems
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

// TelegramPayload represents the payload structure for sending messages via Telegram API
type TelegramPayload struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

// Problem represents a LeetCode problem with its metadata
type Problem struct {
	Title      string  `json:"title"`
	TitleSlug  string  `json:"titleSlug"`
	Difficulty string  `json:"difficulty"`
	AcRate     float64 `json:"acRate"`
	IsPaidOnly bool    `json:"isPaidOnly"`
}

// EasyProblemsResponse represents the response from LeetCode's GraphQL API for easy problems
type EasyProblemsResponse struct {
	Data struct {
		ProblemsetQuestionList struct {
			Questions []Problem `json:"questions"`
		} `json:"problemsetQuestionList"`
	} `json:"data"`
}

// CodeChefProblem represents a CodeChef problem with its metadata
type CodeChefProblem struct {
	Name       string
	Link       string
	Difficulty string
	Category   string
}
