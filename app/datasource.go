package app

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type UserDB struct {
	id         int
	firstName  string
	lastName   string
	telegramId int
}

func FindUserByFirstNameAndLatName(firstName string, lastName string) *UserDB {
	db, err := sql.Open("sqlite3", "./db/dev.db")
	checkError(err)
	defer db.Close()

	var user UserDB
	row := db.QueryRow("SELECT * FROM user WHERE first_name = $1 AND last_name = $2", firstName, lastName)
	err = row.Scan(&user.id, &user.firstName, &user.lastName, &user.telegramId)
	checkError(err)

	return &user
}

func SetTelegramIdToUser(user UserDB, telegramId int) {

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
