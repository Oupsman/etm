package controllers

import (
	"ETM/pkg/models"
	"ETM/pkg/types"
	"encoding/json"
	"errors"
	"github.com/SherClockHolmes/webpush-go"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
)

type Notifs struct {
	Pubkey  string
	PrivKey string
}

func SendTelegramMessage(telegramConfig types.TelegramConfig, text string) error {
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

func GenerateVapidKeys() error {
	// Generate VAPID Keys
	privateVapidKey, publicVapidKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		log.Fatal("Failed to generated VAPID keys")
		return err
	}
	keys := &models.Keys{
		Pubkey:  publicVapidKey,
		Privkey: privateVapidKey,
	}
	err = models.SaveKeys(keys)
	if err != nil {
		log.Fatal("Failed to save VAPID keys to database")
		return err
	}
	return nil
}

func GetVAPIDKey(c *gin.Context) {

	keys, err := models.GetKeys()

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"public_key": keys.Pubkey,
	})
}

func BrowserSend(message string, browserConfig string) error {

	// Get VAPID keys
	keys, err := models.GetKeys()
	if err != nil {
		return err
	}

	s := &webpush.Subscription{}

	err = json.Unmarshal([]byte(browserConfig), s)
	if err != nil {
		return err
	}

	// Send notification
	response, err := webpush.SendNotification([]byte(message), s, &webpush.Options{
		Subscriber:      "oupsman@oupsman.fr", // Do not include "mailto:"
		VAPIDPublicKey:  keys.Pubkey,
		VAPIDPrivateKey: keys.Privkey,
		//		Topic:           "Game changed price",
		TTL: 120,
	})
	if err != nil {
		return err
	}
	if response == nil {
		return errors.New("no response from web push server")
	}
	if response.StatusCode != 201 {
		return errors.New("failed to send notification")
	}
	defer response.Body.Close()

	return nil
}
