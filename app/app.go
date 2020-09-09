package app

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"os"
)

var (
	//ReaderFile define for test
	ReaderFile = ioutil.ReadFile
)

type App struct {
	token  string
	config *Config
}

type Config struct {
	Welcome string `json:"welcome"`
}

func NewApp() *App {
	return &App{}
}

func (a *App) init() {
	a.token = os.Getenv("TOKEN")
	if a.token == "" {
		log.Panic("token is empty!")
	}
	a.config = &Config{}
	err := a.config.loadConfig()
	if err != nil {
		log.Panic(err.Error())
	}
}

func (c *Config) loadConfig() error {
	fileName := "../config/config.json"
	data, err := ReaderFile(fileName)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("Error to parse config %s", err)
	}
	return nil
}

func (a *App) Start() {
	a.init()
	bot, err := tgbotapi.NewBotAPI(a.token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
