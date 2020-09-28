package app

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type UserDB struct {
	id         int
	firstName  string
	lastName   string
	telegramId int
}

type DataSource struct {
	path string
}

func NewDataSource(path string) *DataSource {
	return &DataSource{path: path}
}

func (ds DataSource) getConnect() *sql.DB {
	db, err := sql.Open("sqlite3", ds.path)
	if err != nil {
		log.Print(err)
		return nil
	}
	return db
}

func (ds DataSource) FindUserByFirstNameAndLatName(firstName string, lastName string) *UserDB {
	connect := ds.getConnect()
	if connect == nil {
		return nil
	}
	defer connect.Close()

	user := new(UserDB)
	row := connect.QueryRow("SELECT * FROM user WHERE first_name = $1 AND last_name = $2", firstName, lastName)
	err := row.Scan(&user.id, &user.firstName, &user.lastName, &user.telegramId)
	if err != nil {
		log.Print(err)
		return nil
	}

	return user
}

func (ds DataSource) SetTelegramIdToUser(user *UserDB, telegramId int) bool {
	connect := ds.getConnect()
	if connect == nil {
		return false
	}
	defer connect.Close()

	result, err := connect.Exec("UPDATE user SET telegram_id = $1 WHERE id = $2", telegramId, user.id)
	if err != nil {
		return false
	}

	log.Print(result.RowsAffected())
	return true
}

func (ds DataSource) FindUserByTelegramId(telegramId int) *UserDB {
	connect := ds.getConnect()
	if connect == nil {
		return nil
	}
	defer connect.Close()

	var user UserDB
	row := connect.QueryRow("SELECT * FROM user WHERE telegram_id = $1", telegramId)
	err := row.Scan(&user.id, &user.firstName, &user.lastName, &user.telegramId)
	if err != nil {
		log.Print(err)
		return nil
	}

	return &user
}
