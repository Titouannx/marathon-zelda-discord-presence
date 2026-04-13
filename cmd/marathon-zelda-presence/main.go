package main

import (
	"log"

	"github.com/Titouannx/marathon-zelda-discord-presence/internal/app"
	"github.com/getlantern/systray"
)

var discordClientID = ""

func main() {
	instance, err := app.New(discordClientID)
	if err != nil {
		log.Fatal(err)
	}

	systray.Run(instance.OnReady, instance.OnExit)
}
