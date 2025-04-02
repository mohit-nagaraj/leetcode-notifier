package types

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

type TelegramPayload struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type Problem struct {
	Title      string  `json:"title"`
	TitleSlug  string  `json:"titleSlug"`
	Difficulty string  `json:"difficulty"`
	AcRate     float64 `json:"acRate"`
	PaidOnly   bool    `json:"paidOnly"`
}

type EasyProblemsResponse struct {
	Data struct {
		ProblemsetQuestionList struct {
			Questions []Problem `json:"questions"`
		} `json:"problemsetQuestionList"`
	} `json:"data"`
}
