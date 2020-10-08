package app

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

type MessageTemplate struct {
	Welcome          string `json:"welcome"`
	AuthMsg          string `json:"auth_msg"`
	MainCommand      string `json:"main_command"`
	UndefinedCommand string `json:"undefined_command"`
}

type UI struct {
	template   *MessageTemplate
	auth       *Authorization
	dataSource *DataSource
}

func NewUI(dataSource *DataSource) *UI {
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

	return &UI{auth: NewAuthorization(dataSource), dataSource: dataSource, template: template}
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

func (u UI) getKeyboardButtons(titles []string) tgbotapi.ReplyKeyboardMarkup {
	var keyboard []tgbotapi.KeyboardButton
	for _, title := range titles {
		button := tgbotapi.NewKeyboardButton(title)
		keyboard = append(keyboard, button)
	}
	return tgbotapi.NewReplyKeyboard(keyboard)
}

func (u UI) chooseMessage(cmdText string, msg *tgbotapi.MessageConfig) *tgbotapi.MessageConfig {
	command := u.dataSource.FindCommand(cmdText)

	if command != nil {
		msg.Text = command.Text
		var buttons []string
		if command.Buttons != nil {
			buttons = command.Buttons
		} else {
			buttons = []string{u.template.MainCommand}
		}
		msg.ReplyMarkup = u.getKeyboardButtons(buttons)
		return msg
	} else {
		msg.Text = u.template.UndefinedCommand
		msg.ReplyMarkup = u.getKeyboardButtons([]string{u.template.MainCommand})
		return msg
	}
}
