package main

import (
	"log"
	"os"
	"strings"

	"github.com/Titouannx/marathon-zelda-discord-presence/internal/app"
	"github.com/Titouannx/marathon-zelda-discord-presence/internal/linking"
	"github.com/Titouannx/marathon-zelda-discord-presence/internal/platform"
	"github.com/getlantern/systray"
)

var discordClientID = ""

const defaultWebBase = "https://www.loon.bzh"

func webBase() string {
	if override := strings.TrimSpace(os.Getenv("MARATHON_WEB_BASE")); override != "" {
		return override
	}
	return defaultWebBase
}

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

	// Au premier lancement (aucun config.json valide a cote de l'executable), on
	// lie le compte via la boucle localhost : le navigateur s'ouvre, l'utilisateur
	// clique « Connecter » et le jeton est ecrit localement. Ensuite, ce bloc ne
	// fait rien (demarrage silencieux, y compris au demarrage de session Windows).
	if err := linking.EnsureLinked(webBase()); err != nil {
		_ = platform.ShowError("Marathon Zelda", "Connexion du compte impossible.\nRelance le programme pour reessayer.")
		log.Print(err)
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
