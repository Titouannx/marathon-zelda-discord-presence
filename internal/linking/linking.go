// Package linking gere la liaison du compte au premier lancement, sans jamais
// embarquer de jeton dans le telechargement.
//
// Flux (boucle localhost, anti-CSRF) :
//  1. L'app ecoute sur 127.0.0.1:<port ephemere> (jamais 0.0.0.0).
//  2. Elle genere un `state` aleatoire et ouvre le navigateur sur
//     <webBase>/zelda/presence/connect?port=<port>&state=<state>.
//  3. L'utilisateur (deja connecte au site) clique « Connecter » ; le site
//     renvoie le navigateur vers http://127.0.0.1:<port>/ avec le meme `state`,
//     le `token` de presence et l'`statusUrl`.
//  4. L'app verifie le `state` (comparaison a temps constant), ecrit config.json
//     a cote de l'executable, puis arrete le serveur local.
//
// Le jeton n'autorise QUE la lecture du statut de presence (jeu en cours).
package linking

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Titouannx/marathon-zelda-discord-presence/internal/config"
	"github.com/Titouannx/marathon-zelda-discord-presence/internal/platform"
)

const linkTimeout = 5 * time.Minute

const successHTML = `<!doctype html><html lang="fr"><head><meta charset="utf-8">` +
	`<meta name="viewport" content="width=device-width,initial-scale=1">` +
	`<title>Marathon Zelda</title></head>` +
	`<body style="margin:0;min-height:100vh;display:flex;align-items:center;justify-content:center;` +
	`background:#0a0e06;color:#f4f2b4;font-family:system-ui,sans-serif;text-align:center;padding:24px">` +
	`<div><h1 style="font-size:22px;margin:0 0 8px">Compte connecté !</h1>` +
	`<p style="opacity:.85;margin:0">Tu peux fermer cet onglet et revenir à Discord.</p></div></body></html>`

const failureHTML = `<!doctype html><html lang="fr"><head><meta charset="utf-8">` +
	`<meta name="viewport" content="width=device-width,initial-scale=1">` +
	`<title>Marathon Zelda</title></head>` +
	`<body style="margin:0;min-height:100vh;display:flex;align-items:center;justify-content:center;` +
	`background:#0a0e06;color:#f4f2b4;font-family:system-ui,sans-serif;text-align:center;padding:24px">` +
	`<div><h1 style="font-size:22px;margin:0 0 8px">Échec de la connexion</h1>` +
	`<p style="opacity:.85;margin:0">Relance l'application pour réessayer.</p></div></body></html>`

type linkResult struct {
	token     string
	statusURL string
}

// EnsureLinked ne fait rien si une configuration valide existe deja. Sinon, il
// lance le flux de liaison (ouvre le navigateur) et ecrit config.json.
func EnsureLinked(webBase string) error {
	if _, err := config.Load(); err == nil {
		return nil
	}
	return runLinkFlow(webBase)
}

func runLinkFlow(webBase string) error {
	state, err := randomState()
	if err != nil {
		return fmt.Errorf("generation du state: %w", err)
	}

	// Ecoute UNIQUEMENT sur la boucle locale, port ephemere.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("ecoute locale: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	resultCh := make(chan linkResult, 1)
	mux := http.NewServeMux()
	mux.HandleFunc("/", makeCallbackHandler(state, resultCh))
	server := &http.Server{Handler: mux}

	go func() { _ = server.Serve(listener) }()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()

	connectURL := fmt.Sprintf(
		"%s/zelda/presence/connect?port=%d&state=%s",
		strings.TrimRight(webBase, "/"),
		port,
		url.QueryEscape(state),
	)
	if err := platform.OpenURL(connectURL); err != nil {
		return fmt.Errorf("ouverture du navigateur: %w", err)
	}

	select {
	case res := <-resultCh:
		return config.Save(config.Config{
			PresenceToken:       res.token,
			StatusURL:           res.statusURL,
			PollIntervalSeconds: 30,
		})
	case <-time.After(linkTimeout):
		return errors.New("delai de liaison depasse")
	}
}

func makeCallbackHandler(expectedState string, resultCh chan<- linkResult) http.HandlerFunc {
	var once sync.Once
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		state := query.Get("state")
		token := query.Get("token")
		statusURL := query.Get("statusUrl")

		valid := token != "" && statusURL != "" &&
			subtle.ConstantTimeCompare([]byte(state), []byte(expectedState)) == 1

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if !valid {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(failureHTML))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(successHTML))
		once.Do(func() {
			resultCh <- linkResult{token: token, statusURL: statusURL}
		})
	}
}

func randomState() (string, error) {
	buffer := make([]byte, 24)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}
