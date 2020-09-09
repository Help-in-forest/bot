package app

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenApp_Init_LoadTokenFromEnv(t *testing.T) {
	app := NewApp()

	os.Setenv("TOKEN", "test")
	app.init()

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
