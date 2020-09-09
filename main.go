package main

import "github.com/Help-in-forest/bot/app"

func main() {
	srv := app.NewApp()
	srv.Start()
}
