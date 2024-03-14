package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/joho/godotenv"

	"github.com/snekROmonoro/Browser-Version-Monitor/db"
	"github.com/snekROmonoro/Browser-Version-Monitor/global"
	"github.com/snekROmonoro/Browser-Version-Monitor/monitor"
)

func main() {
	// load .env
	godotenv.Load()

	_tickerSeconds := os.Getenv("TICKER_SECONDS")
	if _tickerSeconds == "" {
		log.Panicln("TICKER_SECONDS is not set")
	}

	tickerSeconds, err := strconv.Atoi(_tickerSeconds)
	if err != nil {
		log.Panicf("failed to parse TICKER_SECONDS: %s", err)
	}

	log.Printf("TICKER_SECONDS: %d", tickerSeconds)

	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramBotToken != "" {
		telegramChannelId := os.Getenv("TELEGRAM_CHANNEL_ID")
		if telegramChannelId == "" {
			log.Panicln("TELEGRAM_CHANNEL_ID is not set")
		}

		if err := global.TelegramBot.Init(telegramBotToken, telegramChannelId); err != nil {
			log.Panicf("failed to initialize telegram bot: %s", err)
		}

		defer global.TelegramBot.Close()
	}

	if err := global.InitDatabase(); err != nil {
		log.Panicf("failed to initialize database: %s", err)
	}

	defer global.CloseDatabase()

	chanClose := make(chan os.Signal, 1)
	signal.Notify(chanClose, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(time.Duration(tickerSeconds) * time.Second)
	defer ticker.Stop()

	log.Println("starting...")
	callMonitors()

	for {
		select {
		case <-chanClose:
			log.Println("shutting down...")
			return
		case <-ticker.C:
			callMonitors()
		}
	}
}

func callMonitors() {
	for _, monitorFunc := range monitor.MonitorFuncs {
		result, err := monitorFunc()
		if err != nil {
			log.Printf("failed to monitor browser version: %s", err)
		} else {
			currVersion, err := version.NewVersion(result.Version)
			if err != nil {
				log.Printf("failed to parse current version %s %s: %s", result.Browser, result.Version, err)
				continue
			}

			found, _ := global.DatabaseClient.Browser.FindFirst(
				db.Browser.BrowserName.Equals(result.Browser),
			).Exec(context.Background())

			var foundUpdate bool = false
			var majorUpdate bool = false

			if found != nil {
				storedVersion, err := version.NewVersion(found.Version)
				if err != nil {
					log.Printf("failed to parse stored version %s %s: %s", result.Browser, found.Version, err)
					continue
				}

				if currVersion.GreaterThan(storedVersion) {
					global.DatabaseClient.Browser.FindUnique(
						db.Browser.ID.Equals(found.ID),
					).Update(
						db.Browser.BrowserName.Set(result.Browser),
						db.Browser.Version.Set(result.Version),
					).Exec(context.Background())

					foundUpdate = true
					majorUpdate = currVersion.Segments()[0] > storedVersion.Segments()[0]
				}
			} else {
				foundUpdate = true
				majorUpdate = true

				global.DatabaseClient.Browser.CreateOne(
					db.Browser.BrowserName.Set(result.Browser),
					db.Browser.Version.Set(result.Version),
				).Exec(context.Background())
			}

			if foundUpdate {
				log.Printf("new version found (major %v): %s %s", majorUpdate, result.Browser, result.Version)
				global.TelegramBot.SendMessage(fmt.Sprintf("Browser `%s` got an update\nCurrent version: `%s`\nMajor: `%v`", result.Browser, result.Version, majorUpdate))
			}
		}
	}
}
