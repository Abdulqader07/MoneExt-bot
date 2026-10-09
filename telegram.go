package main

import (
	"fmt"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	botOnce sync.Once
	bot     *tgbotapi.BotAPI
	botErr  error
)

// getBot initializes the bot exactly once, even if called concurrently.
func getBot() (*tgbotapi.BotAPI, error) {
	botOnce.Do(func() {
		if err := godotenv.Load(); err != nil {
			botErr = fmt.Errorf("error loading .env: %w", err)
			return
		}
		token := os.Getenv("BOT_TOKEN")
		if token == "" {
			botErr = fmt.Errorf("BOT_TOKEN is not set")
			return
		}
		bot, botErr = tgbotapi.NewBotAPI(token)
	})
	return bot, botErr
}

// SendTelegramMessage sends one message to the CHAT_ID in .env.
func SendTelegramMessage(message string) error {
	b, err := getBot()
	if err != nil {
		return err
	}

	chatID := os.Getenv("CHAT_ID")
	if chatID == "" {
		return fmt.Errorf("CHAT_ID is not set")
	}

	var chatIDInt int64
	if _, err := fmt.Sscanf(chatID, "%d", &chatIDInt); err != nil {
		return fmt.Errorf("invalid CHAT_ID: %w", err)
	}

	msg := tgbotapi.NewMessage(chatIDInt, message)
	_, err = b.Send(msg)
	return err
}

// BroadcastAlerts sends multiple messages concurrently.
// Each alert goes to every chat ID provided. Errors are collected.
func BroadcastAlerts(chatIDs []int64, alerts []string) []error {
	b, err := getBot()
	if err != nil {
		return []error{err}
	}

	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		errs   []error
	)

	for _, chatID := range chatIDs {
		for _, alert := range alerts {
			wg.Add(1)
			go func(cid int64, text string) {
				defer wg.Done()
				msg := tgbotapi.NewMessage(cid, text)
				if _, err := b.Send(msg); err != nil {
					mu.Lock()
					errs = append(errs, err)
					mu.Unlock()
				}
			}(chatID, alert)
		}
	}

	wg.Wait()
	return errs
}
