# Marathon Zelda Discord Presence

Client local open-source pour afficher automatiquement la progression du Marathon Zelda sur le profil Discord d'un utilisateur.

## Ce que fait le programme

- lit `config.json` a cote de l'executable
- interroge `https://loon.bzh/api/zelda/presence/status?token=...` toutes les 30 secondes
- met a jour la Rich Presence Discord quand une session Zelda est active
- efface la presence quand aucune session n'est active
- s'enregistre au demarrage de l'OS au premier lancement
- expose une icone discrete dans la zone de notification ou la barre de menu

## Ce que le programme ne fait pas

- ne lit pas vos messages Discord
- n'accede pas a votre compte LOONDASHBOARD hors du token de presence
- ne peut ni demarrer ni arreter une session Zelda
- ne collecte pas de telemetrie cote client

## Configuration attendue

Le site genere automatiquement un `config.json` au telechargement.

```json
{
  "presenceToken": "opaque-token",
  "statusUrl": "https://loon.bzh/api/zelda/presence/status"
}
```

## Build

Le `Discord Application ID` est injecte au build via `-ldflags`.

### Windows

```powershell
$env:DISCORD_CLIENT_ID="123456789012345678"
go build -ldflags "-H=windowsgui -X main.discordClientID=$env:DISCORD_CLIENT_ID" -o bin/marathon-zelda-presence.exe ./cmd/marathon-zelda-presence
Compress-Archive -Path bin/marathon-zelda-presence.exe -DestinationPath dist/marathon-zelda-presence-windows.zip -Force
```

### macOS Intel

```bash
DISCORD_CLIENT_ID=123456789012345678 GOOS=darwin GOARCH=amd64 \
  go build -ldflags "-X main.discordClientID=${DISCORD_CLIENT_ID}" \
  -o bin/marathon-zelda-presence-darwin-amd64 ./cmd/marathon-zelda-presence
zip -j dist/marathon-zelda-presence-macos-intel.zip bin/marathon-zelda-presence-darwin-amd64
```

### macOS Apple Silicon

```bash
DISCORD_CLIENT_ID=123456789012345678 GOOS=darwin GOARCH=arm64 \
  go build -ldflags "-X main.discordClientID=${DISCORD_CLIENT_ID}" \
  -o bin/marathon-zelda-presence-darwin-arm64 ./cmd/marathon-zelda-presence
zip -j dist/marathon-zelda-presence-macos-arm64.zip bin/marathon-zelda-presence-darwin-arm64
```

## Release assets attendus par LOONDASHBOARD

- `marathon-zelda-presence-windows.zip`
- `marathon-zelda-presence-macos-intel.zip`
- `marathon-zelda-presence-macos-arm64.zip`

## Remarques macOS

La premiere version n'est pas signee ni notarized. Le premier lancement demandera donc une validation manuelle dans Gatekeeper.
