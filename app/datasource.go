package app

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type UserDB struct {
	id         int
	firstName  string
	lastName   string
	TelegramID sql.NullInt64
}

type DataSource struct {
	path string
}

func NewDataSource(path string) (*DataSource, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &DataSource{path: path}, nil
}

func (ds DataSource) Select(query string) *sql.Row {
	db, err := sql.Open("sqlite3", ds.path)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer db.Close()

	return db.QueryRow(query)
}

func (ds DataSource) Update(query string) *sql.Result {
	db, err := sql.Open("sqlite3", ds.path)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer db.Close()

	result, err := db.Exec(query)
	if err != nil {
		return nil
	}
	return &result
}

func (ds DataSource) FindUserByFirstNameAndLatName(firstName string, lastName string) *UserDB {
	query := fmt.Sprintf("SELECT * FROM user WHERE first_name = '%s' AND last_name = '%s' AND telegram_id = null", firstName, lastName)
	row := ds.Select(query)
	if row == nil {
		return nil
	}

	user := new(UserDB)
	err := row.Scan(&user.id, &user.firstName, &user.lastName, &user.TelegramID)
	if err != nil {
		log.Print(err)
		return nil
	}
	return user
}

func (ds DataSource) SetTelegramIdToUser(user *UserDB, telegramId int) bool {
	query := fmt.Sprintf("UPDATE user SET telegram_id = '%d' WHERE id = '%d'", telegramId, user.id)
	result := ds.Update(query)
	if result == nil {
		return false
	}
	return true
}

func (ds DataSource) FindUserByTelegramId(TelegramID int) *UserDB {
	query := fmt.Sprintf("SELECT * FROM user WHERE telegram_id = '%d'", TelegramID)
	row := ds.Select(query)
	if row == nil {
		return nil
	}

	user := new(UserDB)
	err := row.Scan(&user.id, &user.firstName, &user.lastName, &user.TelegramID)
	if err != nil {
		log.Print(err)
		return nil
	}
	return user
}
