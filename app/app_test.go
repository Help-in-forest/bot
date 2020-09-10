package app

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SetUp() *App {
	app := NewApp()
	os.Setenv("TOKEN", "test")
	return app
}

func CleanUp() {
	os.Setenv("TOKEN", "")
}

func TestWhenApp_Init_LoadTokenFromEnv(t *testing.T) {
	app := NewApp()

	os.Setenv("TOKEN", "test")
	app.init()
	// Clean up
	os.Setenv("TOKEN", "")

	assert.Equal(t, "test", app.token)
}

func TestWhenApp_Init_PanicsWithoutToken(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Error must be thrown when token is empty!")
		}
		if r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	app := NewApp()

	app.init()
}

func TestWhenApp_Init_LoadMessages(t *testing.T) {
	app := NewApp()

	os.Setenv("TOKEN", "test")
	app.init()
	// Clean up
	os.Setenv("TOKEN", "")

	assert.Equal(t, "Welcome", app.config.Welcome)
}

func TestWhenApp_GetStartMessage_ShowWelcomeMessage(t *testing.T) {
	app := NewApp()
	app.config.Welcome = "Welcome"

	msg := app.chooseMsg("/start")

	assert.Equal(t, "Welcome", msg)
}

func TestWhenApp_Launch_LoadAuthData(t *testing.T) {
	defer CleanUp()
	app := SetUp()

	app.init()

	assert.Equal(t, 1, len(app.users))
}
