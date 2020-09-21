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
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var user UserDB
	row := db.QueryRow("SELECT * FROM user WHERE first_name = $1 AND last_name = $2", firstName, lastName)
	err = row.Scan(&user.id, &user.firstName, &user.lastName, &user.telegramId)
	if err != nil {
		fmt.Println(err)
	}

	return &user
}

func SetTelegramIdToUser(user *UserDB, telegramId int) bool {
	db, err := sql.Open("sqlite3", "./db/dev.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	result, err := db.Exec("UPDATE user SET telegram_id = $1 WHERE id = $2", telegramId, user.id)
	if err != nil {
		return false
	}
	fmt.Println(result.RowsAffected())
	return true
}

func FindUserByTelegramId(telegramId int) *UserDB {
	db, err := sql.Open("sqlite3", "./db/dev.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var user UserDB
	row := db.QueryRow("SELECT * FROM user WHERE telegram_id = $1", telegramId)
	err = row.Scan(&user.id, &user.firstName, &user.lastName, &user.telegramId)
	if err != nil {
		fmt.Println(err)
	}

	return &user
}
