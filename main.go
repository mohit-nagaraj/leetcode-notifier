package main

import (
	"fmt"
	"os"

	"github.com/mohit-nagaraj/leetcode-notifier/types"
	"github.com/mohit-nagaraj/leetcode-notifier/utils"
)

func sendTelegramMessage(botToken, chatID, message, message2, message3 string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	greet1payload := types.TelegramPayload{
		ChatID: chatID,
		Text:   "4th year problem",
	}
	if err := utils.SendMessageWorker(url, greet1payload); err != nil {
		return err
	}

	payload := types.TelegramPayload{
		ChatID: chatID,
		Text:   message,
	}
	if err := utils.SendMessageWorker(url, payload); err != nil {
		return err
	}

	greet2payload := types.TelegramPayload{
		ChatID: chatID,
		Text:   "3rd year problem",
	}
	if err := utils.SendMessageWorker(url, greet2payload); err != nil {
		return err
	}

	payload2 := types.TelegramPayload{
		ChatID: chatID,
		Text:   message2,
	}
	if err := utils.SendMessageWorker(url, payload2); err != nil {
		return err
	}

	greet3payload := types.TelegramPayload{
		ChatID: chatID,
		Text:   "2nd year problem",
	}
	if err := utils.SendMessageWorker(url, greet3payload); err != nil {
		return err
	}

	payload3 := types.TelegramPayload{
		ChatID: chatID,
		Text:   message3,
	}
	if err := utils.SendMessageWorker(url, payload3); err != nil {
		return err
	}

	return nil
}

func main() {
	title, link, err := utils.FetchDailyProblem()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	title2, link2, err2 := utils.FetchEasyProblemOfTheDay()
	if err2 != nil {
		fmt.Println("Error2:", err)
		return
	}

	codechefProblem, err3 := utils.FetchRandomCodeChefProblem("codechef_problems.csv")
	if err3 != nil {
		fmt.Println("Error3:", err3)
		return
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" || chatID == "" {
		fmt.Println("Error: Telegram credentials not set in environment variables")
		fmt.Println("Please set TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID")
		return
	}

	message := fmt.Sprintf("Today's LeetCode Problem: %s\n%s\nDear students, please find the daily challenge posted for today ☝️", title, link)
	message2 := fmt.Sprintf("Today's LeetCode Problem: %s\n%s\nDear students, please find the daily challenge posted for today ☝️", title2, link2)
	message3 := fmt.Sprintf("Today's CodeChef Problem: %s\n%s (%s, %s)\nDear students, please find the daily challenge posted for today ☝️", codechefProblem.Name, codechefProblem.Link, codechefProblem.Category, codechefProblem.Difficulty)

	if err := sendTelegramMessage(botToken, chatID, message, message2, message3); err != nil {
		fmt.Println("Error sending msg:", err)
		return
	}

	fmt.Println("Telegram notification sent successfully!")
}
