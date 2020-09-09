package app

import (
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
