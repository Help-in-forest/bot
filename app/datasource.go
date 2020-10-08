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
	teamID     int
}

type Team struct {
	id   int
	name string
}

type Command struct {
	ID      int
	Command string
	Text    string
	Buttons []string
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

func (ds DataSource) querySelectMany(query string) *sql.Rows {
	db, err := sql.Open("sqlite3", ds.path)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return nil
	}
	return rows
}

func (ds DataSource) querySelectOne(query string) *sql.Row {
	db, err := sql.Open("sqlite3", ds.path)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer db.Close()

	return db.QueryRow(query)
}

func (ds DataSource) queryUpdate(query string) *sql.Result {
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
	query := fmt.Sprintf("SELECT * FROM user WHERE first_name = '%s' AND last_name = '%s' AND telegram_id is null", firstName, lastName)
	row := ds.querySelectOne(query)
	if row == nil {
		return nil
	}

	user := new(UserDB)
	err := row.Scan(&user.id, &user.firstName, &user.lastName, &user.TelegramID, &user.teamID)
	if err != nil {
		log.Print(err)
		return nil
	}
	return user
}

func (ds DataSource) SetTelegramIdToUser(user *UserDB, telegramId int) bool {
	query := fmt.Sprintf("UPDATE user SET telegram_id = '%d' WHERE id = '%d'", telegramId, user.id)
	result := ds.queryUpdate(query)
	if result == nil {
		return false
	}
	return true
}

func (ds DataSource) FindUserByTelegramId(TelegramID int) *UserDB {
	query := fmt.Sprintf("SELECT * FROM user WHERE telegram_id = '%d'", TelegramID)
	row := ds.querySelectOne(query)
	if row == nil {
		return nil
	}

	user := new(UserDB)
	err := row.Scan(&user.id, &user.firstName, &user.lastName, &user.TelegramID, &user.teamID)
	if err != nil {
		log.Print(err)
		return nil
	}
	return user
}

func (ds DataSource) FindCommand(cmdText string) *Command {
	query := fmt.Sprintf("SELECT id, command, text FROM command WHERE command = '%s'", cmdText)
	row := ds.querySelectOne(query)
	if row == nil {
		return nil
	}

	cmd := new(Command)
	err := row.Scan(&cmd.ID, &cmd.Command, &cmd.Text)
	if err != nil {
		log.Print(err)
		return nil
	}

	query = fmt.Sprintf("SELECT cmd.command FROM command cmd INNER JOIN command_keyboard ck ON ck.child_id = cmd.id WHERE parent_id = '%d'", cmd.ID)
	rows := ds.querySelectMany(query)
	if rows == nil {
		return cmd
	}

	var buttons []string
	for rows.Next() {
		var title string
		err := rows.Scan(&title)
		if err != nil {
			fmt.Println(err)
			continue
		}
		buttons = append(buttons, title)
	}
	cmd.Buttons = buttons
	return cmd
}
