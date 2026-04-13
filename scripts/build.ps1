$ErrorActionPreference = "Stop"

if (-not $env:DISCORD_CLIENT_ID) {
  throw "DISCORD_CLIENT_ID is required."
}

New-Item -ItemType Directory -Force -Path bin, dist | Out-Null

go build -ldflags "-H=windowsgui -X main.discordClientID=$env:DISCORD_CLIENT_ID" `
  -o bin/marathon-zelda-presence.exe `
  ./cmd/marathon-zelda-presence

Compress-Archive `
  -Path bin/marathon-zelda-presence.exe `
  -DestinationPath dist/marathon-zelda-presence-windows.zip `
  -Force
