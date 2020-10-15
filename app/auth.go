package app

import (
	"strings"
)

type Authorization struct {
	dataSource *DataSource
}

func NewAuthorization(source *DataSource) *Authorization {
	return &Authorization{dataSource: source}
}

func (a *Authorization) CheckAuthorization(userMessage *Message) bool {
	user := a.dataSource.FindUserByTelegramId(userMessage.TelegramID)
	return user != nil
}

func (a *Authorization) TryAuthorize(msg *Message) bool {
	data := strings.Split(msg.Text, " ")
	if len(data) < 2 {
		return false
	}

	user := a.dataSource.FindUserByFirstNameAndLatName(data[0], data[1])
	if user != nil {
		result := a.dataSource.SetTelegramIdToUser(user, msg.TelegramID)
		return result
	}
	return false
}
