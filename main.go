package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

const interval = 4 * time.Hour

func buildMessage(results []PriceResult) string {
	var sb strings.Builder
	sb.WriteString("*Crypto & Metals Update*\n")
	sb.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	sb.WriteString("\n\n")

	for _, r := range results {
		sb.WriteString(fmt.Sprintf("• %s: $%.2f\n", r.Name, r.Price))
	}

	sb.WriteString("\n_Next update in 4h_")
	return sb.String()
}

func runOnce() {
	results := GetAllAssets()

	msg := buildMessage(results)

	if err := SendTelegramMessage(msg); err != nil {
		log.Println("send failed:", err)
		return
	}
	log.Println("prices sent successfully")
}

func main() {
	log.Println("Bot started — interval:", interval)

	// Run immediately on start
	runOnce()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		runOnce()
	}
}