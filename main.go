package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
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

	log.Println("starting...")
	callMonitors()

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(tickerSeconds).Seconds().Do(func() {
		callMonitors()
	})
	scheduler.StartAsync()

	chanClose := make(chan os.Signal, 1)
	signal.Notify(chanClose, os.Interrupt, syscall.SIGTERM)
	<-chanClose
	log.Println("shutting down...")
	// scheduler.Stop()
}

func callMonitors() {
	log.Println("calling monitors...")
	for i, monitorFunc := range monitor.MonitorFuncs {
		result, err := monitorFunc()
		if err != nil {
			log.Printf("failed to monitor %d browser version: %s", i, err)
		} else {
			currVersion, err := version.NewVersion(result.Version)
			if err != nil {
				log.Printf("failed to parse current version %s %s: %s", result.Browser, result.Version, err)
				continue
			}

			var storedVersion *version.Version
			found, _ := global.DatabaseClient.Browser.FindFirst(
				db.Browser.BrowserName.Equals(result.Browser),
			).Exec(context.Background())

			var foundUpdate bool = false
			var majorUpdate bool = false

			if found != nil {
				storedVersion, err = version.NewVersion(found.Version)
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
				if err := global.TelegramBot.SendMessage(result.UpdateString(storedVersion)); err != nil {
					log.Printf("failed to send telegram message: %s", err)
				}
			}
		}
	}
}
