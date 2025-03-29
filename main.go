package main

import (
	"fmt"
	"os"

	"github.com/mohit-nagaraj/leetcode-notifier/types"
	"github.com/mohit-nagaraj/leetcode-notifier/utils"
)

func sendTelegramMessage(botToken, chatID, message, message2 string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	greet1payload := types.TelegramPayload{
		ChatID: chatID,
		Text:   "6th sem problem",
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
		Text:   "4th sem problem",
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

	return nil
}

func main() {
	title, link, err := utils.FetchDailyProblem()
	if err != nil {
		fmt.Println("Error:", err)
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
	message2 := fmt.Sprintf("Today's LeetCode Problem: %s\n%s\nDear students, please find the daily challenge posted for today ☝️", title, link)

	if err := sendTelegramMessage(botToken, chatID, message, message2); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Telegram notification sent successfully!")
}
