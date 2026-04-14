package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/Titouannx/marathon-zelda-discord-presence/internal/assets"
	"github.com/Titouannx/marathon-zelda-discord-presence/internal/config"
	"github.com/Titouannx/marathon-zelda-discord-presence/internal/discord"
	"github.com/Titouannx/marathon-zelda-discord-presence/internal/model"
	"github.com/Titouannx/marathon-zelda-discord-presence/internal/platform"
	"github.com/getlantern/systray"
)

type App struct {
	cfg           config.Config
	rpc           *discord.Client
	httpClient    *http.Client
	cancel        context.CancelFunc
	lastPresence  *model.PresenceStatus
	currentStatus *model.PresenceStatus
	openMarathon  *systray.MenuItem
	uninstall     *systray.MenuItem
	quit          *systray.MenuItem
}

const (
	appTitle                  = "Marathon Zelda"
	installedAppName          = "Marathon Zelda Activite Discord"
	installationMarkerName    = ".marathon-zelda-presence-installed"
	installationSuccessText   = "Activite Discord Marathon Zelda installee.\nCelle-ci s'affichera lors du prochain demarrage de session."
	uninstallationSuccessText = "Activite Discord Marathon Zelda desinstallee.\nLe demarrage automatique a ete retire. Tu peux maintenant supprimer ce dossier."
	zeldaSiteURL              = "https://www.loon.bzh/zelda"
)

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
	systray.SetTitle(appTitle)
	systray.SetTooltip(appTitle + " - attente de Discord")

	a.openMarathon = systray.AddMenuItem("Ouvrir la page du marathon", "Ouvre la page Zelda sur loon.bzh")
	a.uninstall = systray.AddMenuItem("Desinstaller", "Retire le demarrage auto et ferme le programme")
	systray.AddSeparator()
	a.quit = systray.AddMenuItem("Quitter", "Ferme le programme")
	a.syncTrayActions(nil)
	_, _ = a.ensureInstalled()

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
		case <-a.openMarathon.ClickedCh:
			if err := platform.OpenURL(zeldaSiteURL); err != nil {
				_ = platform.ShowError(appTitle, "Impossible d'ouvrir la page du marathon.")
			}
		case <-a.uninstall.ClickedCh:
			_ = a.performUninstall()
			_ = platform.ShowInfo(appTitle, uninstallationSuccessText)
			systray.Quit()
			return
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
		a.lastPresence = nil
		a.currentStatus = nil
		a.syncTrayActions(nil)
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
		a.lastPresence = nil
		a.currentStatus = nil
		a.syncTrayActions(nil)
		_ = a.rpc.Clear()
		systray.SetTooltip("Marathon Zelda - aucune session active")
		return
	}

	current := payload
	a.currentStatus = &current
	a.syncTrayActions(a.currentStatus)

	if shouldUpdatePresence(a.lastPresence, payload) {
		if err := a.rpc.Set(payload); err != nil {
			systray.SetTooltip("Marathon Zelda - Discord Desktop indisponible")
			return
		}
		sent := payload
		a.lastPresence = &sent
	}

	systray.SetTooltip("Marathon Zelda - " + payload.GameName)
}

func (a *App) syncTrayActions(current *model.PresenceStatus) {
	if a.openMarathon == nil {
		return
	}

	if current != nil && current.Active {
		a.openMarathon.Hide()
		return
	}

	a.openMarathon.Show()
}

func (a *App) ensureInstalled() (bool, error) {
	if platform.IsAutoStartEnabled() {
		_ = a.writeInstallationMarker()
		_ = platform.EnsureAppRegistration(installedAppName)
		return true, nil
	}

	hasMarker, err := a.hasInstallationMarker()
	if err != nil {
		return false, err
	}
	if hasMarker {
		return false, nil
	}

	if err := platform.InstallAutoStart(); err != nil {
		_ = platform.ShowError(appTitle, "Impossible d'activer le demarrage automatique pour Marathon Zelda.")
		return false, err
	}
	if err := platform.EnsureAppRegistration(installedAppName); err != nil {
		_ = platform.ShowError(appTitle, "Impossible d'enregistrer la desinstallation Windows pour Marathon Zelda.")
	}

	if err := a.writeInstallationMarker(); err != nil {
		return true, err
	}

	_ = platform.ShowInfo(appTitle, installationSuccessText)
	return true, nil
}

func (a *App) hasInstallationMarker() (bool, error) {
	path, err := installationMarkerPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (a *App) writeInstallationMarker() error {
	path, err := installationMarkerPath()
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte("installed\n"), 0o644)
}

func (a *App) removeInstallationMarker() error {
	return removeInstallationMarkerFile()
}

func (a *App) performUninstall() error {
	_ = a.rpc.Clear()
	_ = platform.RemoveAutoStart()
	_ = platform.RemoveAppRegistration()
	return removeInstallationMarkerFile()
}

func UninstallCurrentInstallation() error {
	_ = platform.RemoveAutoStart()
	_ = platform.RemoveAppRegistration()
	return removeInstallationMarkerFile()
}

func removeInstallationMarkerFile() error {
	path, err := installationMarkerPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func installationMarkerPath() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(executable), installationMarkerName), nil
}

func resolveProfileTarget(lastProfile string) string {
	if lastProfile != "" {
		return lastProfile
	}

	return zeldaSiteURL
}

func shouldUpdatePresence(current *model.PresenceStatus, next model.PresenceStatus) bool {
	if current == nil {
		return true
	}

	return *current != next
}

func visibleTrayMenuLabels(current *model.PresenceStatus) []string {
	if current != nil && current.Active {
		return []string{"Desinstaller", "Quitter"}
	}

	return []string{"Ouvrir la page du marathon", "Desinstaller", "Quitter"}
}
