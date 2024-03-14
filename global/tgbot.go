package global

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type telegramBot struct {
	bot       *tgbotapi.BotAPI
	channelId int64
}

var TelegramBot telegramBot

func (t *telegramBot) GetBot() *tgbotapi.BotAPI {
	return t.bot
}

func (t *telegramBot) Init(token string, channelId string) error {
	var err error
	t.bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	log.Printf("Authorized on telegram account %s\n", t.bot.Self.UserName)

	channelIdInt, err := strconv.ParseInt(channelId, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse channelId: %w", err)
	}

	t.channelId = channelIdInt

	return nil
}

func (t *telegramBot) Close() {
	t.bot.StopReceivingUpdates()
}

func (t *telegramBot) SendMessage(text string) error {
	if t.bot == nil {
		return fmt.Errorf("bot is not initialized")
	}

	if t.channelId == 0 {
		return fmt.Errorf("channelId is not set")
	}

	message := tgbotapi.NewMessage(t.channelId, text)
	message.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := t.bot.Send(message)
	return err
}
