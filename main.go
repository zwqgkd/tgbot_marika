package main

import (
	"context"
	"go_tgbot/command"
	"go_tgbot/config"
	"go_tgbot/database"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var bot *tgbotapi.BotAPI
var logger *log.Logger

func replyHandler(message *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	log.Println(text)
	msg.ReplyToMessageID = message.MessageID
	msg.ParseMode = "markdownV2"
	if _, err := bot.Send(msg); err != nil {
		msg.ParseMode = ""
		bot.Send(msg)
	}
}

func commandHandler(update tgbotapi.Update) {
	replyText := ""

	// Extract the command from the Message.
	switch update.Message.Command() {
	case "help":
		replyText = "I understand /ask"
	case "ask":
		if text, err := command.GeminiAsk(update.Message.Text); err != nil {
			logger.Println("Error making gemini request:", err)
			replyText = "I'm tired..."
		} else {
			replyText = text
		}
	default:
		replyText = "I don't know that command"
	}
	replyHandler(update.Message, replyText)
}

func uploadDocument(db *mongo.Database, message *tgbotapi.Message) {
	db.Collection(message.Chat.Title).InsertOne(context.TODO(), bson.D{
		{Key: "data", Value: time.Unix(int64(message.Date), 0)},
		{Key: "content", Value: message.Text},
		{Key: "userName", Value: message.From.UserName},
	}, nil)
}

func init() {
	go config.InitConfig()
	//connect to mongodb
	go database.ConnectMongoDB()
}
func main() {
	var err error
	//ser logger
	file, err := os.Create("log.txt")
	if err != nil {
		log.Println("Error creating log file:", err)
	}
	defer file.Close()
	logger = log.New(file, "[TgBot]", log.LstdFlags)

	//create botapi
	proxyURL, err := url.Parse(config.SetConfig.ProxyUrl)
	if err != nil {
		bot, err = tgbotapi.NewBotAPIWithClient(config.SetConfig.BotToken,
			"https://api.telegram.org/bot%s/%s", &http.Client{})
	} else {
		bot, err = tgbotapi.NewBotAPIWithClient(config.SetConfig.BotToken,
			"https://api.telegram.org/bot%s/%s", &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
				},
			})
	}

	tgbotapi.SetLogger(logger)
	if err != nil {
		logger.Println(err)
	}

	bot.Debug = true

	logger.Printf("Authorized on account %s", bot.Self.UserName)

	//handler updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		go func(update tgbotapi.Update) {
			if update.Message != nil && update.Message.Chat.Type != "supergroup" {
				replyHandler(update.Message, "I'm only for supergroup")
				return
			}
			if update.Message != nil {
				go uploadDocument(database.MongoDB, update.Message)
				if update.Message.IsCommand() { // ignore any non-command Messages
					commandHandler(update)
				}
			}
		}(update)
	}
}
