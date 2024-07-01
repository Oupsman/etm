package controllers

import (
	"ETM/models"
	"ETM/types"
	"github.com/SherClockHolmes/webpush-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
)

type Notifs struct {
	Pubkey  string
	PrivKey string
}

func (notif *Notifs) SendTelegramMessage(telegramConfig types.TelegramConfig, text string) error {
	client := &http.Client{}

	bot, err := tgbotapi.NewBotAPIWithClient(telegramConfig.Token, tgbotapi.APIEndpoint, client)
	if err != nil {
		return err
	}
	_, err = bot.Send(tgbotapi.NewMessage(telegramConfig.ChatID, text))
	if err != nil {
		return err
	}
	return nil
}

func (notif *Notifs) GenerateVapidKeys() (*models.Keys, error) {
	// Generate VAPID Keys
	privateVapidKey, publicVapidKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		log.Fatal("Failed to generated VAPID keys")
		return nil, err
	}
	return &models.Keys{
		Pubkey:  publicVapidKey,
		Privkey: privateVapidKey,
	}, nil
}

func (notif *Notifs) BrowserSend() {

}
