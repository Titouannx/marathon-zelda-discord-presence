package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Titouannx/marathon-zelda-discord-presence/internal/assets"
	"github.com/Titouannx/marathon-zelda-discord-presence/internal/config"
	"github.com/Titouannx/marathon-zelda-discord-presence/internal/discord"
	"github.com/Titouannx/marathon-zelda-discord-presence/internal/model"
	"github.com/Titouannx/marathon-zelda-discord-presence/internal/platform"
	"github.com/getlantern/systray"
)

type App struct {
	cfg         config.Config
	rpc         *discord.Client
	httpClient  *http.Client
	cancel      context.CancelFunc
	lastProfile string
	autoStart   *systray.MenuItem
	openProfile *systray.MenuItem
	quit        *systray.MenuItem
	about       *systray.MenuItem
}

func New(discordClientID string) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	rpc, err := discord.New(discordClientID)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg: cfg,
		rpc: rpc,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (a *App) OnReady() {
	systray.SetIcon(assets.TrayIcon)
	systray.SetTitle("Marathon Zelda")
	systray.SetTooltip("Marathon Zelda - attente de Discord")

	a.openProfile = systray.AddMenuItem("Ouvrir mon profil", "Ouvre le profil Marathon Zelda actuel")
	a.autoStart = systray.AddMenuItem("Lancer au demarrage", "Active ou desactive le demarrage auto")
	a.about = systray.AddMenuItem("A propos", "Ouvre la page du marathon")
	systray.AddSeparator()
	a.quit = systray.AddMenuItem("Quitter", "Ferme le programme")

	if platform.IsAutoStartEnabled() {
		a.autoStart.Check()
	}

	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	go a.menuLoop()
	go a.pollLoop(ctx)
}

func (a *App) OnExit() {
	if a.cancel != nil {
		a.cancel()
	}
	a.rpc.Close()
}

func (a *App) menuLoop() {
	for {
		select {
		case <-a.openProfile.ClickedCh:
			if a.lastProfile != "" {
				_ = platform.OpenURL(a.lastProfile)
			}
		case <-a.about.ClickedCh:
			_ = platform.OpenURL("https://loon.bzh/zelda")
		case <-a.autoStart.ClickedCh:
			if platform.IsAutoStartEnabled() {
				if err := platform.RemoveAutoStart(); err == nil {
					a.autoStart.Uncheck()
				}
			} else {
				if err := platform.InstallAutoStart(); err == nil {
					a.autoStart.Check()
				}
			}
		case <-a.quit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func (a *App) pollLoop(ctx context.Context) {
	a.refresh()

	ticker := time.NewTicker(time.Duration(a.cfg.PollIntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.refresh()
		}
	}
}

func (a *App) refresh() {
	requestURL, err := url.Parse(a.cfg.StatusURL)
	if err != nil {
		systray.SetTooltip("Marathon Zelda - URL invalide")
		return
	}

	query := requestURL.Query()
	query.Set("token", a.cfg.PresenceToken)
	requestURL.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		systray.SetTooltip("Marathon Zelda - requete invalide")
		return
	}

	res, err := a.httpClient.Do(req)
	if err != nil {
		systray.SetTooltip("Marathon Zelda - reseau indisponible")
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		_ = a.rpc.Clear()
		systray.SetTooltip("Marathon Zelda - token invalide")
		return
	}
	if res.StatusCode >= 500 {
		systray.SetTooltip("Marathon Zelda - serveur temporairement indisponible")
		return
	}
	if res.StatusCode != http.StatusOK {
		systray.SetTooltip(fmt.Sprintf("Marathon Zelda - erreur HTTP %d", res.StatusCode))
		return
	}

	var payload model.PresenceStatus
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		systray.SetTooltip("Marathon Zelda - reponse invalide")
		return
	}

	if !payload.Active {
		a.lastProfile = ""
		_ = a.rpc.Clear()
		systray.SetTooltip("Marathon Zelda - aucune session active")
		return
	}

	if err := a.rpc.Set(payload); err != nil {
		systray.SetTooltip("Marathon Zelda - Discord Desktop indisponible")
		return
	}

	a.lastProfile = payload.ProfileURL
	systray.SetTooltip("Marathon Zelda - " + payload.GameName)
}
