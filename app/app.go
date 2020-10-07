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
	token        string
	config       *Config
	users        map[string][]User
	authorized   map[string]struct{}
	homeKeyboard tgbotapi.InlineKeyboardMarkup
	dataSource   *DataSource
	auth         *Authorization
}

type Config struct {
	Welcome    string `json:"welcome"`
	AuthMsg    string `json:"auth_msg"`
	Authorized string `json:"authorized"`
	TeamsTitle string `json:"teams_button_title"`
}

type User struct {
	Name    string
	Surname string
	Data    string
}

type Message struct {
	UserName   string
	Text       string
	TelegramID int
}

func NewApp() *App {
	return &App{config: &Config{}, users: map[string][]User{}, authorized: map[string]struct{}{}}
}

func (a *App) init() {
	a.token = os.Getenv("TOKEN")
	if a.token == "" {
		log.Panic("token is empty!")
	}
	err := a.config.loadConfig()
	if err != nil {
		log.Panic(err.Error())
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

	a.auth = NewAuthorization(a.dataSource)

	a.homeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(a.config.TeamsTitle, a.config.TeamsTitle),
		),
	)
}

func (c *Config) loadConfig() error {
	fileName := "../config/custom.json"
	data, err := ReaderFile(fileName)
	if err != nil {
		data, err = ReaderFile("../config/config.json")
		if err != nil {
			data, err = ReaderFile("config/config.json")
			if err != nil {
				return err
			}
		}
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

		if update.CallbackQuery != nil {
			fmt.Print(update)
		}

		userMsg := &Message{UserName: update.Message.From.UserName, Text: update.Message.Text, TelegramID: update.Message.From.ID}
		text := a.handle(userMsg)
		keyboard := a.chooseKeyboard(text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)

		if len(keyboard.InlineKeyboard) > 0 {
			msg.ReplyMarkup = keyboard
		}

		bot.Send(msg)
	}
}

func (a *App) chooseMsg(command string) string {
	switch command {
	case "/start":
		return a.config.Welcome
	default:
		return command
	}
}

func (a *App) chooseKeyboard(text string) tgbotapi.InlineKeyboardMarkup {
	switch text {
	case a.config.Authorized:
		return a.homeKeyboard
	default:
		return tgbotapi.InlineKeyboardMarkup{}
	}
}

func (a *App) handle(msg *Message) string {
	var authorized bool
	if !a.auth.CheckAuthorization(msg.TelegramID) {
		authorized = a.auth.Authorize(msg)
		if !authorized {
			return a.config.AuthMsg
		}
		return a.config.Authorized
	}
	return a.chooseMsg(msg.Text)
}
