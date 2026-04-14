package main

import (
	"log"
	"os"

	"github.com/Titouannx/marathon-zelda-discord-presence/internal/app"
	"github.com/Titouannx/marathon-zelda-discord-presence/internal/platform"
	"github.com/getlantern/systray"
)

var discordClientID = ""

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--uninstall" {
		if err := app.UninstallCurrentInstallation(); err != nil {
			_ = platform.ShowError("Marathon Zelda", "Impossible de desinstaller l'activite Discord Marathon Zelda.")
			log.Print(err)
			return
		}

		_ = platform.ShowInfo("Marathon Zelda", "Activite Discord Marathon Zelda desinstallee.")
		return
	}

	instance, err := app.New(discordClientID)
	if err != nil {
		_ = platform.ShowError("Marathon Zelda", "Impossible de lancer l'activite Discord Marathon Zelda.")
		log.Print(err)
		return
	}

	systray.Run(instance.OnReady, instance.OnExit)
}
