$ErrorActionPreference = 'Stop'

$gopath = go env GOPATH
$rsrc = Join-Path $gopath 'bin\rsrc.exe'
if (-not (Test-Path $rsrc)) {
  throw 'rsrc.exe not found. Run: go install github.com/akavel/rsrc@latest'
}

& $rsrc -arch amd64 -ico internal/assets/tray_icon_windows.ico -o cmd/marathon-zelda-presence/resource_windows_amd64.syso
