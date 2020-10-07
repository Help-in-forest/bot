package app

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

type MessageTemplate struct {
	Welcome    string `json:"welcome"`
	AuthMsg    string `json:"auth_msg"`
	Authorized string `json:"authorized"`
	TeamsTitle string `json:"teams_button_title"`
}

type UI struct {
	template *MessageTemplate
	auth     *Authorization
}

func NewUI(authorization *Authorization) *UI {
	templatePath := os.Getenv("MESSAGE_TEMPLATE_PATH")
	if templatePath == "" {
		log.Panic("MESSAGE_TEMPLATE_PATH is empty!")
	}
	data, err := ReaderFile(templatePath)
	if err != nil {
		log.Panic("Invalid MESSAGE_TEMPLATE_PATH. File does not exist!")
	}

	template := new(MessageTemplate)
	if err := json.Unmarshal(data, template); err != nil {
		log.Panic(err.Error())
	}

	return &UI{auth: authorization, template: template}
}

func (u UI) HandleMessage(userMessage *Message) *tgbotapi.MessageConfig {
	msg := new(tgbotapi.MessageConfig)
	msg.ChatID = userMessage.CharID

	if userMessage.Text == "/start" {
		msg.Text = u.template.Welcome
		return msg
	}
	if !u.auth.CheckAuthorization(userMessage) {
		msg.Text = u.template.AuthMsg
		return msg
	}

	return u.chooseMessage(userMessage.Text, msg)
}

func (u UI) getMainMenu() tgbotapi.ReplyKeyboardMarkup {
	buttons := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(u.template.TeamsTitle),
		tgbotapi.NewKeyboardButton("Спринт"),
	}
	return tgbotapi.NewReplyKeyboard(buttons)
}

func (u UI) getTeamsMenu() tgbotapi.ReplyKeyboardMarkup {
	buttons := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("Команда 1"),
		tgbotapi.NewKeyboardButton("Палата 6"),
		tgbotapi.NewKeyboardButton("Risk&Dream"),
	}
	return tgbotapi.NewReplyKeyboard(buttons)
}

func (u UI) chooseMessage(command string, msg *tgbotapi.MessageConfig) *tgbotapi.MessageConfig {
	switch command {
	case u.template.TeamsTitle:
		msg.Text = u.template.TeamsTitle
		msg.ReplyMarkup = u.getTeamsMenu()
		return msg
	default:
		msg.Text = command
		msg.ReplyMarkup = u.getMainMenu()
		return msg
	}
}
