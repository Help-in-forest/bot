package app

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

type App struct {
	token      string
	dataSource *DataSource
	ui         *UI
}

type Message struct {
	CharID     int64
	TelegramID int
	UserName   string
	Text       string
}

func NewApp() *App {
	return &App{}
}

func (a *App) init() {
	a.token = os.Getenv("TOKEN")
	if a.token == "" {
		log.Panic("token is empty!")
	}

	dataSourcePath := os.Getenv("DB_PATH")
	if dataSourcePath == "" {
		log.Panic("DB_PATH is empty!")
	}
	ds, err := NewDataSource(dataSourcePath)
	if err != nil {
		log.Panic("Invalid DB_PATH. DB does not exist!")
	}
	a.dataSource = ds
	a.ui = NewUI(ds)
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
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.CallbackQuery != nil {
			fmt.Print(update)
		}

		userMsg := &Message{
			CharID:     update.Message.Chat.ID,
			TelegramID: update.Message.From.ID,
			UserName:   update.Message.From.UserName,
			Text:       update.Message.Text,
		}

		msg := a.ui.HandleMessage(userMsg)
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
