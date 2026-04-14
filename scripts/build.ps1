$ErrorActionPreference = "Stop"

$effectiveClientId = if ($env:DISCORD_PRESENCE_CLIENT_ID) { $env:DISCORD_PRESENCE_CLIENT_ID } else { $env:DISCORD_CLIENT_ID }

if (-not $effectiveClientId) {
  throw "DISCORD_PRESENCE_CLIENT_ID or DISCORD_CLIENT_ID is required."
}

New-Item -ItemType Directory -Force -Path bin, dist | Out-Null

go build -ldflags "-H=windowsgui -X main.discordClientID=$effectiveClientId" `
  -o bin/marathon-zelda-presence.exe `
  ./cmd/marathon-zelda-presence

Compress-Archive `
  -Path bin/marathon-zelda-presence.exe `
  -DestinationPath dist/marathon-zelda-presence-windows.zip `
  -Force
